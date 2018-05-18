package proxy

import (
	"io"
	"net"
	"sync"
)

//Command struct for command
type Command struct {
	rqst    *Rqst
	conn    net.Conn
	connOut net.Conn
	wg      sync.WaitGroup
}

const (
	cmdConnect = 0x01
	cmdBind    = 0x02
	cmdUDP     = 0x03

	rspSuccess           = 0x00
	rspServerError       = 0x01
	rspCommandNotSupport = 0x07
)

//Fire method which check command and run it
func (c *Command) Fire() (rsp byte, err error) {

	switch c.rqst.cmd {
	case cmdConnect:
		return c.connect()
	case cmdBind:
		return rspCommandNotSupport, nil
	case cmdUDP:
		return rspCommandNotSupport, nil
	default:
		return rspCommandNotSupport, nil
	}
}

func (c *Command) connect() (rsp byte, err error) {
	to := c.rqst.FQDN()

	if c.connOut == nil {
		if c.connOut, err = net.Dial("tcp", to); err != nil {
			return rspServerError, err
		}
	}

	defer c.connOut.Close()

	r := &Rsps{
		socksVer: socks5Ver,
		rsp:      rspSuccess,
		reserved: reservedSymbol,
	}

	err = r.parseAddr(c.connOut.LocalAddr().String())
	if err != nil {
		return rspServerError, nil
	}

	buf, err := r.Bytes()

	if err != nil {
		return rspServerError, err
	}

	if _, err = c.conn.Write(buf); err != nil {
		return rspServerError, err
	}

	c.proxy()

	return
}

func (c *Command) proxy() {
	c.wg.Add(2)

	go func(to io.Writer, from io.Reader) {
		defer c.wg.Done()
		io.Copy(to, from)
		if conn, ok := to.(net.Conn); ok {
			conn.Close()
		}
	}(c.connOut, c.conn)

	go func(to io.Writer, from io.Reader) {
		defer c.wg.Done()
		io.Copy(to, from)
		if conn, ok := from.(net.Conn); ok {
			conn.Close()
		}
	}(c.conn, c.connOut)

	c.wg.Wait()
}
