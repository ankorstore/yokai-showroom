package config

import "strings"

// Protocol is an enum for the supported database connection protocols.
type Protocol string

const (
	UnknownProtocol Protocol = "unknown"
	TCPProtocol     Protocol = "tcp"
	SocketProtocol  Protocol = "socket"
	BufConnProtocol Protocol = "bufconn"
)

// String returns a string representation of the [Protocol].
func (d Protocol) String() string {
	return string(d)
}

// FetchProtocol returns a [Protocol] for a given name.
func FetchProtocol(name string) Protocol {
	//nolint:exhaustive
	switch d := Protocol(strings.ToLower(name)); d {
	case TCPProtocol,
		SocketProtocol,
		BufConnProtocol:
		return d
	default:
		return UnknownProtocol
	}
}
