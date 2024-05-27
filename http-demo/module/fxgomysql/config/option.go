package config

import (
	"google.golang.org/grpc/test/bufconn"
)

// ConfigOptions are options for the [GoMySQLServerFactory] implementations.
type ConfigOptions struct {
	Protocol Protocol
	Listener *bufconn.Listener
	Socket   string
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

// DefaultMySQLServerConfigOptions are the default options used in the [DefaultGoMySQLServerFactory].
func DefaultMySQLServerConfigOptions() ConfigOptions {
	return ConfigOptions{
		Protocol: TCPProtocol,
		Listener: bufconn.Listen(1024 * 1024),
		Socket:   DefaultSocket,
		User:     DefaultUser,
		Password: DefaultPassword,
		Host:     DefaultHost,
		Port:     DefaultPort,
		Database: DefaultDatabase,
	}
}

// MySQLServerConfigOption are functional options for the [GoMySQLServerFactory] implementations.
type MySQLServerConfigOption func(o *ConfigOptions)

func WithProtocol(protocol Protocol) MySQLServerConfigOption {
	return func(o *ConfigOptions) {
		o.Protocol = protocol
	}
}

func WithSocket(socket string) MySQLServerConfigOption {
	return func(o *ConfigOptions) {
		if socket != "" {
			o.Socket = socket
		}
	}
}

func WithUser(user string) MySQLServerConfigOption {
	return func(o *ConfigOptions) {
		if user != "" {
			o.User = user
		}
	}
}

func WithPassword(password string) MySQLServerConfigOption {
	return func(o *ConfigOptions) {
		if password != "" {
			o.Password = password
		}
	}
}

func WithHost(host string) MySQLServerConfigOption {
	return func(o *ConfigOptions) {
		if host != "" {
			o.Host = host
		}
	}
}

func WithPort(port int) MySQLServerConfigOption {
	return func(o *ConfigOptions) {
		if port != 0 {
			o.Port = port
		}
	}
}

func WithDatabase(database string) MySQLServerConfigOption {
	return func(o *ConfigOptions) {
		if database != "" {
			o.Database = database
		}
	}
}
