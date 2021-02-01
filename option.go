package mcserver

import "net"

type Option func(*MCServer)

func WithListener(lis net.Listener) Option {
	return func(srv *MCServer) {
		srv.listener = lis
	}
}
