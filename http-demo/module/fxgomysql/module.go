package fxgomysql

import (
	"context"
	basesql "database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/ankorstore/yokai-showroom/http-demo/module/fxgomysql/config"
	"github.com/ankorstore/yokai-showroom/http-demo/module/fxgomysql/server"
	yokaiconfig "github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	gomysqlserver "github.com/dolthub/go-mysql-server/server"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

var dialerOnce sync.Once

// ModuleName is the module name.
const ModuleName = "gomysqlserver"

// FxGoMySQLServerModule is the [Fx] go mysql module.
//
// [Fx]: https://github.com/uber-go/fx
var FxGoMySQLServerModule = fx.Options(
	fx.Module(
		ModuleName,
		fx.Provide(
			fx.Annotate(
				server.NewDefaultGoMySQLServerFactory,
				fx.As(new(server.GoMySQLServerFactory)),
			),
			fx.Annotate(
				NewFxGoMySQLServerModuleInfo,
				fx.As(new(interface{})),
				fx.ResultTags(`group:"core-module-infos"`),
			),
			NewFxGoMySQLServerConfig,
			NewFxGoMySQLServer,
		),
	),
	//ConnectFxGoMySQLServer(),
	AAA(),
)

// FxGoMySQLServerConfigParam allows injection of the required dependencies in [NewGoMySQLServerConfig].
type FxGoMySQLServerConfigParam struct {
	fx.In
	Config *yokaiconfig.Config
}

// NewFxGoMySQLServerConfig returns a new [config.GoMySQLServerConfig] instance
func NewFxGoMySQLServerConfig(p FxGoMySQLServerConfigParam) (*config.GoMySQLServerConfig, error) {
	// protocol
	protocolConfig := p.Config.GetString("modules.gomysqlserver.config.protocol")
	protocol := config.FetchProtocol(protocolConfig)
	if protocol == config.UnknownProtocol {
		return nil, fmt.Errorf("unknown protocol: %s", protocolConfig)
	}

	// config
	return config.NewGoMySQLServerConfig(
		config.WithProtocol(protocol),
		config.WithSocket(p.Config.GetString("modules.gomysqlserver.config.socket")),
		config.WithUser(p.Config.GetString("modules.gomysqlserver.config.user")),
		config.WithPassword(p.Config.GetString("modules.gomysqlserver.config.password")),
		config.WithHost(p.Config.GetString("modules.gomysqlserver.config.host")),
		config.WithPort(p.Config.GetInt("modules.gomysqlserver.config.port")),
		config.WithDatabase(p.Config.GetString("modules.gomysqlserver.config.database")),
	), nil
}

// FxGoMySQLServerParam allows injection of the required dependencies in [NewFxMySQLMemoryServer].
type FxGoMySQLServerParam struct {
	fx.In
	LifeCycle      fx.Lifecycle
	Factory        server.GoMySQLServerFactory
	ServerConfig   *config.GoMySQLServerConfig
	Config         *yokaiconfig.Config
	Logger         *log.Logger
	LogBuffer      logtest.TestLogBuffer
	TracerProvider oteltrace.TracerProvider
	DB             *basesql.DB
}

// NewFxGoMySQLServer returns a new [gomysqlserver.Server] instance
func NewFxGoMySQLServer(p FxGoMySQLServerParam) (*gomysqlserver.Server, error) {
	// server options
	options := []server.MySQLMemoryServerOption{
		server.WithConfig(p.ServerConfig),
	}

	logOutput := io.Discard
	if p.Config.GetBool("modules.mysqlmemory.log.enabled") {
		logOutput = os.Stdout
		options = append(options, server.WithLogOutput(logOutput))
	}

	if p.Config.GetBool("modules.mysqlmemory.trace.enabled") {
		options = append(options, server.WithTracer(p.TracerProvider.Tracer(ModuleName)))
	}

	// server creation
	srv, err := p.Factory.Create(options...)
	if err != nil {
		return nil, err
	}

	fmt.Printf("****** PROTOCOL: %v\n", p.ServerConfig.Protocol())

	// server start
	p.Logger.Info().Msgf("starting go mysql server on port %d", p.ServerConfig.Port())
	go srv.Start()

	// server lifecycle
	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if p.ServerConfig.Protocol() != config.BufConnProtocol {
				p.Logger.Info().Msg("stopping go mysql server")

				return srv.Close()
			}

			return nil
		},
	})

	return srv, nil
}

func AAA() fx.Option {
	return fx.Invoke(
		func(db *basesql.DB, server *gomysqlserver.Server) *basesql.DB {
			return db
		},
	)
}

func ConnectFxGoMySQLServer() fx.Option {
	return fx.Decorate(
		func(db *basesql.DB, serverConfig *config.GoMySQLServerConfig) (*basesql.DB, error) {
			if serverConfig.Protocol() == config.BufConnProtocol {
				fmt.Println("****** IN CONNECT")
				connector, err := db.Driver().(driver.DriverContext).OpenConnector(serverConfig.String())
				if err != nil {
					return nil, err
				}

				err = db.Close()
				if err != nil {
					return nil, err
				}

				return basesql.OpenDB(connector), nil
			}

			return db, nil
		},
	)
}
