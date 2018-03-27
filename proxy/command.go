package proxy

import (
	"io"
	"net"
	"sync"
)

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

func (c *Command) Fire() (err error, rsp byte) {

	switch c.rqst.cmd {
	case cmdConnect:
		return c.connect()
	case cmdBind:
		return nil, rspCommandNotSupport
	case cmdUDP:
		return nil, rspCommandNotSupport
	default:
		return nil, rspCommandNotSupport
	}
}

func (c *Command) connect() (err error, rsp byte) {
	to := c.rqst.FQDN()

	if c.connOut == nil {
		if c.connOut, err = net.Dial("tcp", to); err != nil {
			return err, rspServerError
		}
	}

	defer c.connOut.Close()

	r := &Rspc{
		socksVer: socks5Ver,
		rsp:      rspSuccess,
		reserved: reservedSymbol,
	}

	err = r.parseAddr(c.connOut.LocalAddr().String())
	if err != nil && err == rspServerErrorMsg {
		return nil, rspServerError
	}

	buf, err := r.Bytes()

	if err != nil {
		return err, rspServerError
	}

	if _, err = c.conn.Write(buf); err != nil {
		return err, rspServerError
	}

	errCh := make(chan error, 2)

	go c.proxy(c.conn, c.connOut, errCh)
	go c.proxy(c.connOut, c.conn, errCh)

	c.wg.Wait()

	if err := <-errCh; err != nil {
		return err, rspServerError
	}

	return
}

func (c *Command) proxy(to io.Writer, from io.Reader, errCh chan error) {
	c.wg.Add(1)
	defer c.wg.Done()

	_, err := io.Copy(to, from)
	if conn, ok := to.(net.Conn); ok {
		conn.Close()
	}
	errCh <- err

}
