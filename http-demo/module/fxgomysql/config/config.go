package config

import (
	"fmt"

	"google.golang.org/grpc/test/bufconn"
)

const (
	DefaultNetwork        = "mysqlmemory"
	DefaultSocket         = "/tmp/mysql.sock"
	DefaultUser           = "user"
	DefaultPassword       = "password"
	DefaultHost           = "localhost"
	DefaultPort           = 3306
	DefaultDatabase       = "db"
	DefaultBufConnAddress = "bufconn"
)

type GoMySQLServerConfig struct {
	protocol Protocol
	listener *bufconn.Listener
	socket   string
	user     string
	password string
	host     string
	port     int
	database string
}

func NewGoMySQLServerConfig(options ...MySQLServerConfigOption) *GoMySQLServerConfig {
	// resolve options
	appliedOptions := DefaultMySQLServerConfigOptions()
	for _, opt := range options {
		opt(&appliedOptions)
	}

	// create config
	return &GoMySQLServerConfig{
		protocol: appliedOptions.Protocol,
		listener: appliedOptions.Listener,
		socket:   appliedOptions.Socket,
		user:     appliedOptions.User,
		password: appliedOptions.Password,
		host:     appliedOptions.Host,
		port:     appliedOptions.Port,
		database: appliedOptions.Database,
	}
}

func (c *GoMySQLServerConfig) Protocol() Protocol {
	return c.protocol
}

func (c *GoMySQLServerConfig) Listener() *bufconn.Listener {
	return c.listener
}

func (c *GoMySQLServerConfig) Socket() string {
	return c.socket
}

func (c *GoMySQLServerConfig) User() string {
	return c.user
}

func (c *GoMySQLServerConfig) Password() string {
	return c.password
}

func (c *GoMySQLServerConfig) Host() string {
	return c.host
}

func (c *GoMySQLServerConfig) Port() int {
	return c.port
}

func (c *GoMySQLServerConfig) Database() string {
	return c.database
}

func (c *GoMySQLServerConfig) String() string {
	var str string

	switch c.Protocol() {
	case TCPProtocol:
		str = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			c.User(),
			c.Password(),
			c.Host(),
			c.Port(),
			c.Database(),
		)
	case SocketProtocol:
		str = fmt.Sprintf(
			"%s:%s@unix(%s)/%s",
			c.User(),
			c.Password(),
			c.Socket(),
			c.Database(),
		)
	case BufConnProtocol:
		str = fmt.Sprintf(
			"%s:%s@mysqlmemory(bufconn)/%s",
			c.User(),
			c.Password(),
			c.Database(),
		)
	}

	return str
}
