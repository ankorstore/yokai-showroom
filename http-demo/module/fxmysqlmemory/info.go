package fxmysqlmemory

import (
	"github.com/ankorstore/yokai-showroom/http-demo/module/fxmysqlmemory/config"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
)

// FxMySQLMemoryModuleInfo is a module info collector for fxmysqlmemory.
type FxMySQLMemoryModuleInfo struct {
	server *server.Server
	config *config.MySQLMemoryServerConfig
}

// NewFxMySQLMemoryModuleInfo returns a new [FxMySQLMemoryModuleInfo].
func NewFxMySQLMemoryModuleInfo(server *server.Server, config *config.MySQLMemoryServerConfig) *FxMySQLMemoryModuleInfo {
	return &FxMySQLMemoryModuleInfo{
		server: server,
		config: config,
	}
}

// Name return the name of the module info.
func (i *FxMySQLMemoryModuleInfo) Name() string {
	return ModuleName
}

// Data return the data of the module info.
func (i *FxMySQLMemoryModuleInfo) Data() map[string]interface{} {
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
			"options":  i.config.Options(),
		},
		"sessions": sessionVars,
	}
}
