package fxmysqlmemory

import (
	"context"
	basesql "database/sql"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"github.com/ankorstore/yokai-showroom/http-demo/module/fxmysqlmemory/config"
	"github.com/ankorstore/yokai-showroom/http-demo/module/fxmysqlmemory/server"
	yokaiconfig "github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	gomysqlserver "github.com/dolthub/go-mysql-server/server"
	"github.com/go-sql-driver/mysql"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

var dialerOnce sync.Once

// ModuleName is the module name.
const ModuleName = "mysqlmemory"

// FxMySQLMemoryModule is the [Fx] log module.
//
// [Fx]: https://github.com/uber-go/fx
var FxMySQLMemoryModule = fx.Module(
	ModuleName,
	fx.Provide(
		fx.Annotate(
			server.NewDefaultMySQLMemoryServerFactory,
			fx.As(new(server.MySQLMemoryServerFactory)),
		),
		fx.Annotate(
			NewFxMySQLMemoryModuleInfo,
			fx.As(new(interface{})),
			fx.ResultTags(`group:"core-module-infos"`),
		),
		NewFxMySQLMemoryConfig,
		NewFxMySQLMemoryServer,
	),
)

// FxMySQLMemoryConfigParam allows injection of the required dependencies in [NewFxMySQLMemoryConfig].
type FxMySQLMemoryConfigParam struct {
	fx.In
	Config *yokaiconfig.Config
}

// NewFxMySQLMemoryConfig returns a new [FxMySQLMemoryDSNParam] instance
func NewFxMySQLMemoryConfig(p FxMySQLMemoryConfigParam) (*config.MySQLMemoryServerConfig, error) {
	// protocol
	protocolConfig := p.Config.GetString("modules.mysqlmemory.config.protocol")
	protocol := config.FetchProtocol(protocolConfig)
	if protocol == config.UnknownProtocol {
		return nil, fmt.Errorf("unknown protocol: %s", protocolConfig)
	}

	// config
	return config.NewMySQLMemoryServerConfig(
		config.WithProtocol(protocol),
		config.WithSocket(p.Config.GetString("modules.mysqlmemory.config.socket")),
		config.WithUser(p.Config.GetString("modules.mysqlmemory.config.user")),
		config.WithPassword(p.Config.GetString("modules.mysqlmemory.config.password")),
		config.WithHost(p.Config.GetString("modules.mysqlmemory.config.host")),
		config.WithPort(p.Config.GetInt("modules.mysqlmemory.config.port")),
		config.WithDatabase(p.Config.GetString("modules.mysqlmemory.config.database")),
		config.WithOptions(p.Config.GetString("modules.mysqlmemory.config.options")),
	), nil
}

// FxMySQLMemoryServerParam allows injection of the required dependencies in [NewFxMySQLMemoryServer].
type FxMySQLMemoryServerParam struct {
	fx.In
	LifeCycle      fx.Lifecycle
	Factory        server.MySQLMemoryServerFactory
	ServerConfig   *config.MySQLMemoryServerConfig
	Config         *yokaiconfig.Config
	Logger         *log.Logger
	LogBuffer      logtest.TestLogBuffer
	TracerProvider oteltrace.TracerProvider
}

// NewFxMySQLMemoryServer returns a new [server.Server] instance
func NewFxMySQLMemoryServer(p FxMySQLMemoryServerParam) (*gomysqlserver.Server, error) {
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

	// server start
	p.Logger.Info().Msgf("starting mysql memory server on port %d", p.ServerConfig.Port())
	go srv.Start()

	// server lifecycle
	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if p.ServerConfig.Protocol() != config.BufConnProtocol {
				p.Logger.Info().Msg("stopping mysql memory server")

				return srv.Close()
			}

			return nil
		},
	})

	return srv, nil
}

func ConnectFxTestMySQLMemoryServer() fx.Option {
	return fx.Decorate(
		func(db *basesql.DB, serverConfig *config.MySQLMemoryServerConfig) (*basesql.DB, error) {
			// activate dial context
			dialerOnce.Do(func() {
				mysql.RegisterDialContext(config.DefaultNetwork, func(ctx context.Context, addr string) (net.Conn, error) {
					return serverConfig.Listener().DialContext(ctx)
				})
			})

			// create connector
			connector, err := mysql.NewConnector(&mysql.Config{
				DBName:                  serverConfig.Database(),
				Addr:                    "bufconn",
				Net:                     config.DefaultNetwork,
				User:                    serverConfig.User(),
				Passwd:                  serverConfig.Password(),
				AllowNativePasswords:    true,
				AllowCleartextPasswords: true,
			})
			if err != nil {
				return nil, err
			}

			// reopen connection to mysql memory server
			newDb := basesql.OpenDB(connector)
			if err != nil {
				return nil, err
			}

			// close existing connection
			err = db.Close()
			if err != nil {
				return nil, err
			}

			return newDb, nil
		},
	)
}
