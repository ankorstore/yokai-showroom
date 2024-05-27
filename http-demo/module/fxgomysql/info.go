package fxgomysql

import (
	"github.com/ankorstore/yokai-showroom/http-demo/module/fxgomysql/config"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
)

// FxGoMySQLServerModuleInfo is a module info collector for fxgomysql.
type FxGoMySQLServerModuleInfo struct {
	server *server.Server
	config *config.GoMySQLServerConfig
}

// NewFxGoMySQLServerModuleInfo returns a new [FxMySQLMemoryModuleInfo].
func NewFxGoMySQLServerModuleInfo(server *server.Server, config *config.GoMySQLServerConfig) *FxGoMySQLServerModuleInfo {
	return &FxGoMySQLServerModuleInfo{
		server: server,
		config: config,
	}
}

// Name return the name of the module info.
func (i *FxGoMySQLServerModuleInfo) Name() string {
	return ModuleName
}

// Data return the data of the module info.
func (i *FxGoMySQLServerModuleInfo) Data() map[string]interface{} {
	sessionVars := make(map[uint32]interface{})

	i.server.SessionManager().Iter(func(session sql.Session) (stop bool, err error) {
		sessionVars[session.ID()] = session.GetAllSessionVariables()

		return false, nil
	})

	return map[string]interface{}{
		"config": map[string]interface{}{
			"protocol": i.config.Protocol(),
			"socket":   i.config.Socket(),
			"user":     i.config.User(),
			"password": i.config.Password(),
			"host":     i.config.Host(),
			"port":     i.config.Port(),
			"database": i.config.Database(),
		},
		"sessions": sessionVars,
	}
}
