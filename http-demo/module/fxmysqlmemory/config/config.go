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
	DefaultOptions        = "parseTime=true"
	DefaultBufConnAddress = "bufconn"
)

type MySQLMemoryServerConfig struct {
	protocol Protocol
	listener *bufconn.Listener
	socket   string
	user     string
	password string
	host     string
	port     int
	database string
	options  string
}

func NewMySQLMemoryServerConfig(options ...MySQLMemoryConfigOption) *MySQLMemoryServerConfig {
	// resolve options
	appliedOptions := DefaultMySQLMemoryConfigOptions()
	for _, opt := range options {
		opt(&appliedOptions)
	}

	// create config
	return &MySQLMemoryServerConfig{
		protocol: appliedOptions.Protocol,
		listener: appliedOptions.Listener,
		socket:   appliedOptions.Socket,
		user:     appliedOptions.User,
		password: appliedOptions.Password,
		host:     appliedOptions.Host,
		port:     appliedOptions.Port,
		database: appliedOptions.Database,
		options:  appliedOptions.Options,
	}
}

func (d *MySQLMemoryServerConfig) Protocol() Protocol {
	return d.protocol
}

func (d *MySQLMemoryServerConfig) Listener() *bufconn.Listener {
	return d.listener
}

func (d *MySQLMemoryServerConfig) Socket() string {
	return d.socket
}

func (d *MySQLMemoryServerConfig) User() string {
	return d.user
}

func (d *MySQLMemoryServerConfig) Password() string {
	return d.password
}

func (d *MySQLMemoryServerConfig) Host() string {
	return d.host
}

func (d *MySQLMemoryServerConfig) Port() int {
	return d.port
}

func (d *MySQLMemoryServerConfig) Database() string {
	return d.database
}

func (d *MySQLMemoryServerConfig) Options() string {
	return d.options
}

func (d *MySQLMemoryServerConfig) String() string {
	var str string

	switch d.protocol {
	case TCPProtocol:
		str = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			d.User(),
			d.Password(),
			d.Host(),
			d.Port(),
			d.Database(),
		)
	case SocketProtocol:
		str = fmt.Sprintf(
			"%s:%s@unix(%s)/%s",
			d.User(),
			d.Password(),
			d.Socket(),
			d.Database(),
		)
	case BufConnProtocol:
		str = "bufconn"
	}

	if len(d.options) > 0 {
		str = fmt.Sprintf("%s?%s", str, d.Options())
	}

	return str
}
