package gnet

import (
	"net"
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
	remoteAddress, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return nil, err
	}
	el := c.server.svr.lb.next(remoteAddress)
	connection := newUDPConn(el, el.svr.ln.addr, remoteAddress)
	return connection, nil
}
