package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/ankorstore/yokai-showroom/http-demo/module/fxmysqlmemory/config"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	gsql "github.com/dolthub/go-mysql-server/sql"
	vsql "github.com/dolthub/vitess/go/mysql"
	"github.com/sirupsen/logrus"
)

var loggerOnce sync.Once

type MySQLMemoryServerFactory interface {
	Create(options ...MySQLMemoryServerOption) (*server.Server, error)
}

type DefaultMySQLMemoryServerFactory struct{}

func NewDefaultMySQLMemoryServerFactory() *DefaultMySQLMemoryServerFactory {
	return &DefaultMySQLMemoryServerFactory{}
}

func (f *DefaultMySQLMemoryServerFactory) Create(options ...MySQLMemoryServerOption) (*server.Server, error) {
	// resolve options
	appliedOptions := DefaultMySQLMemoryServerOptions()
	for _, opt := range options {
		opt(&appliedOptions)
	}

	// configure logger
	loggerOnce.Do(func() {
		logrus.SetOutput(appliedOptions.LogOutput)
	})

	// create engine
	memDB := memory.NewDatabase(appliedOptions.Config.Database())
	memDB.BaseDatabase.EnablePrimaryKeyIndexes()
	memProvider := memory.NewDBProvider(memDB)
	memEngine := sqle.NewDefault(memProvider)

	// create default user
	catalog := memEngine.Analyzer.Catalog.MySQLDb
	memEditor := catalog.Editor()
	defer memEditor.Close()
	catalog.AddSuperUser(memEditor, appliedOptions.Config.User(), appliedOptions.Config.Host(), appliedOptions.Config.Password())

	// session builder
	memSessionBuilder := memory.NewSessionBuilder(memProvider)

	// create server config
	memServerConfig := server.Config{
		Tracer: appliedOptions.Tracer,
	}

	switch protocol := appliedOptions.Config.Protocol(); protocol {
	case config.TCPProtocol:
		memServerConfig.Protocol = "tcp"
		memServerConfig.Address = fmt.Sprintf("%s:%d", appliedOptions.Config.Host(), appliedOptions.Config.Port())
	case config.SocketProtocol:
		memServerConfig.Protocol = "unix"
		memServerConfig.Socket = appliedOptions.Config.Socket()
	case config.BufConnProtocol:
		memServerConfig.Protocol = "bufconn"
		memServerConfig.Listener = appliedOptions.Config.Listener()

		catalog.AddSuperUser(memEditor, appliedOptions.Config.User(), config.DefaultBufConnAddress, appliedOptions.Config.Password())

		memSessionBuilder = func(ctx context.Context, c *vsql.Conn, addr string) (gsql.Session, error) {
			host := ""
			user := ""
			mysqlConnectionUser, ok := c.UserData.(gsql.MysqlConnectionUser)
			if ok {
				host = mysqlConnectionUser.Host
				user = mysqlConnectionUser.User
			}

			client := gsql.Client{
				Address:      host,
				User:         user,
				Capabilities: c.Capabilities,
			}

			return memory.NewSession(
				gsql.NewBaseSessionWithClientServer(addr, client, c.ConnectionID),
				memProvider,
			), nil
		}
	case config.UnknownProtocol:
		return nil, fmt.Errorf("unknown protocol: %s", protocol)
	}

	return server.NewServer(memServerConfig, memEngine, memSessionBuilder, nil)
}
