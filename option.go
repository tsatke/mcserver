package mcserver

import (
	"net"

	"github.com/rs/zerolog"
)

// Option is an API function that can be passed into New to customize the created server
// with optional arguments.
type Option func(*MCServer)

// WithListener tells the server, which listener he should use. If this is given, the server
// will not allocate a separate listener.
func WithListener(lis net.Listener) Option {
	return func(srv *MCServer) {
		srv.listener = lis
	}
}

// WithLogger will make the server use the given logger to write logs.
func WithLogger(log zerolog.Logger) Option {
	return func(srv *MCServer) {
		srv.log = log
	}
}
