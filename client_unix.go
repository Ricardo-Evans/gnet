//+build linux freebsd dragonfly darwin

package gnet

import (
	"github.com/panjf2000/gnet/internal/logging"
	"github.com/panjf2000/gnet/internal/socket"
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

func (c *client) Dial(addr string, localAddr string, socketOpts ...Option) (Conn, error) {
	option := loadOptions(socketOpts...)
	var options []socket.Option
	if option.ReusePort {
		options = append(options, socket.Option{SetSockopt: socket.SetReuseport, Opt: 1})
	}
	if option.SocketRecvBuffer > 0 {
		options = append(options, socket.Option{SetSockopt: socket.SetRecvBuffer, Opt: option.SocketRecvBuffer})
	}
	if option.SocketSendBuffer > 0 {
		options = append(options, socket.Option{SetSockopt: socket.SetSendBuffer, Opt: option.SocketSendBuffer})
	}
	network, address := parseProtoAddr(addr)
	localNetwork, localAddress := parseProtoAddr(addr)
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
