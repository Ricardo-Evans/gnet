//+build linux freebsd dragonfly darwin

package gnet

import (
	"github.com/panjf2000/gnet/internal/socket"
)

type client struct {
	server *Server
}

func newClient(server *Server) *client {
	return &client{
		server,
	}
}

func (c *client) Dial(addr string) (Conn, error) {
	network, address := parseProtoAddr(addr)
	sa, na, err := socket.GetUDPSocketAddr(network, address)
	if err != nil {
		return nil, err
	}
	el := c.server.svr.lb.next(na)
	connection := newUDPConn(el.ln.fd, el, sa)
	return connection, nil
}
