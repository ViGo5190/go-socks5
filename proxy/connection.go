package proxy

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

const (
	socks5Ver = 0x05

	authNoAuth       = 0x00
	authGSSAPI       = 0x01
	authUserPassword = 0x02

	authNoMethods = 0xFF
)

type Connection struct {
	conn   net.Conn
	ctx    *context.Context
	logger *log.Logger

	ver     byte
	methods []byte
}

func NewConnection(conn net.Conn, logger *log.Logger) *Connection {
	ctx := context.Background()
	return &Connection{
		conn:   conn,
		ctx:    &ctx,
		logger: logger,
	}
}

func (c *Connection) Serve() (err error) {
	defer c.conn.Close()

	if err = binary.Read(c.conn, binary.BigEndian, &c.ver); err != nil {
		c.logger.Error(err)
		return
	}

	if c.ver != socks5Ver {
		err = errors.New("unsupported protocol version")
		c.logger.Error(err)
		return
	}

	if err = c.handshake(); err != nil {
		return err
	}

	err = c.cmd()
	if err != nil {
		c.logger.Error(err)
	}
	return
}

func (c *Connection) handshake() (err error) {
	var n int64
	var acceptableAuthMethodSize byte
	if err = binary.Read(c.conn, binary.BigEndian, &acceptableAuthMethodSize); err != nil {
		return
	}
	n++

	c.methods = make([]byte, acceptableAuthMethodSize)
	if _, err = io.ReadFull(c.conn, c.methods); err != nil {
		return
	}
	n += int64(acceptableAuthMethodSize)

	if bytes.IndexByte(c.methods, authNoAuth) != -1 {
		return c.handleNoAuth()
	}

	//TODO: add another methods

	c.conn.Write([]byte{socks5Ver, authNoMethods})
	return
}

func (c *Connection) handleNoAuth() (err error) {
	c.logger.Debug("handleNoAuth")
	c.logger.Debug(c.conn.RemoteAddr().String())

	remAddrstring := c.conn.RemoteAddr().String()

	_, _, err = net.SplitHostPort(remAddrstring)
	if err != nil {
		return
	}

	_, err = c.conn.Write([]byte{socks5Ver, authNoAuth})

	if err != nil {
		c.logger.Error(err)
		return
	}
	return
}

func (c *Connection) cmd() (err error) {

	r := &Rqst{}
	if err = r.fromReader(c.conn); err != nil {
		return
	}

	cmd := &Command{
		rqst: r,
		conn: c.conn,
	}

	err, rsp := cmd.Fire()
	if err != nil {

		return err
	}

	if rsp != rspSuccess {
		return c.writeErrorResponce(r, rsp)
	}

	c.logger.Info("done")
	return
}

func (c *Connection) writeErrorResponce(r *Rqst, errCode byte) (err error) {
	_, err = c.conn.Write([]byte{
		r.socksVer,
		errCode,
		r.reserved,
		r.addressType,
		0, 0, 0, 0,
		0, 0,
	})
	return err
}
