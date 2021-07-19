//+build linux freebsd dragonfly darwin

package gnet

import (
	"github.com/panjf2000/gnet/internal/socket"
	"github.com/panjf2000/gnet/logging"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"sync"
)

type client struct {
	Server
	fd      int
	mutex   sync.Mutex
	address net.Addr
	once    sync.Once
}

func newClient(server Server) *client {
	return &client{
		server,
		-1,
		sync.Mutex{},
		nil,
		sync.Once{},
	}
}

func (c *client) Dial(addr string) (Conn, error) {
	options := c.svr.ln.sockopts
	network, address := parseProtoAddr(addr)
	localNetwork, localAddress := c.svr.ln.network, c.svr.ln.addr
	sa, na, err := socket.GetUDPSocketAddr(network, address)
	if err != nil {
		return nil, err
	}
	if c.fd < 0 {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		if c.fd < 0 {
			fd, lna, err := socket.UDPSocket(localNetwork, localAddress, options...)
			if err != nil {
				return nil, err
			}
			c.fd = fd
			c.address = lna
		}
	}
	connection := &conn{
		fd:         c.fd,
		sa:         sa,
		localAddr:  c.address,
		remoteAddr: na,
	}
	return connection, nil
}

func (c *client) Close() {
	c.once.Do(func() {
		if c.fd > 0 {
			logging.LogErr(os.NewSyscallError("close", unix.Close(c.fd)))
		}
	})
}
