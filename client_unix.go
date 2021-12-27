//go:build linux || freebsd || dragonfly || darwin
// +build linux freebsd dragonfly darwin

package gnet

import (
	"github.com/panjf2000/gnet/internal/socket"
	"sync"
)

type client struct {
	server *Server
	mutex  *sync.Mutex // used to provide thread safety, for load balancer only currently
}

func newClient(server *Server) *client {
	return &client{
		server: server,
		mutex:  &sync.Mutex{},
	}
}

func (c *client) Dial(addr string) (Conn, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	network, address := parseProtoAddr(addr)
	sa, na, err := socket.GetUDPSocketAddr(network, address)
	if err != nil {
		return nil, err
	}
	el := c.server.svr.lb.next(na)
	connection := newUDPConn(el.ln.fd, el, sa)
	return connection, nil
}
