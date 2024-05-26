package config

import (
	"google.golang.org/grpc/test/bufconn"
)

// ConfigOptions are options for the [MySQLMemoryServerFactory] implementations.
type ConfigOptions struct {
	Protocol Protocol
	Listener *bufconn.Listener
	Socket   string
	User     string
	Password string
	Host     string
	Port     int
	Database string
	Options  string
}

// DefaultMySQLMemoryConfigOptions are the default options used in the [DefaultMySQLMemoryServerFactory].
func DefaultMySQLMemoryConfigOptions() ConfigOptions {
	return ConfigOptions{
		Protocol: TCPProtocol,
		Listener: bufconn.Listen(1024),
		Socket:   DefaultSocket,
		User:     DefaultUser,
		Password: DefaultPassword,
		Host:     DefaultHost,
		Port:     DefaultPort,
		Database: DefaultDatabase,
		Options:  DefaultOptions,
	}
}

// MySQLMemoryConfigOption are functional options for the [MySQLMemoryServerFactory] implementations.
type MySQLMemoryConfigOption func(o *ConfigOptions)

func WithProtocol(protocol Protocol) MySQLMemoryConfigOption {
	return func(o *ConfigOptions) {
		o.Protocol = protocol
	}
}

func WithSocket(socket string) MySQLMemoryConfigOption {
	return func(o *ConfigOptions) {
		if socket != "" {
			o.Socket = socket
		}
	}
}

func WithUser(user string) MySQLMemoryConfigOption {
	return func(o *ConfigOptions) {
		if user != "" {
			o.User = user
		}
	}
}

func WithPassword(password string) MySQLMemoryConfigOption {
	return func(o *ConfigOptions) {
		if password != "" {
			o.Password = password
		}
	}
}

func WithHost(host string) MySQLMemoryConfigOption {
	return func(o *ConfigOptions) {
		if host != "" {
			o.Host = host
		}
	}
}

func WithPort(port int) MySQLMemoryConfigOption {
	return func(o *ConfigOptions) {
		if port != 0 {
			o.Port = port
		}
	}
}

func WithDatabase(database string) MySQLMemoryConfigOption {
	return func(o *ConfigOptions) {
		if database != "" {
			o.Database = database
		}
	}
}

func WithOptions(options string) MySQLMemoryConfigOption {
	return func(o *ConfigOptions) {
		if options != "" {
			o.Options = options
		}
	}
}
