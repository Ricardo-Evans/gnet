//+build linux freebsd dragonfly darwin

package gnet

import (
	"github.com/panjf2000/gnet/internal/socket"
)

type client struct {
	Server
}

func (c *client) Dial(addr string, localAddr string, socketOpts ...socket.Option) (Conn, error) {
	network, address := parseProtoAddr(addr)
	sa, na, err := socket.GetUDPSocketAddr(network, address)
	if err != nil {
		return nil, err
	}
	localNetwork, localAddress := parseProtoAddr(localAddr)
	fd, lna, err := socket.UDPSocket(localNetwork, localAddress, socketOpts...)
	if err != nil {
		return nil, err
	}
	el := c.svr.lb.next(na)
	err = el.poller.AddRead(fd)
	connection := &conn{
		fd:         fd,
		sa:         sa,
		localAddr:  lna,
		remoteAddr: na,
	}
	return connection, nil
}
