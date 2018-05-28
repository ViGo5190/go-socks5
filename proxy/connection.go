package proxy

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"unsafe"

	"github.com/rs/zerolog/log"
)

const (
	socks5Ver = 0x05

	authNoAuth = 0x00
	//authGSSAPI = 0x01
	authUP = 0x02

	authNoMethods = 0xFF

	authUPSuccess = 0x00
	authUPFailure = 0x01
)

var (
	errUserPasswordVersion = fmt.Errorf("user password wrong version")
	errUserPassword        = fmt.Errorf("user password wrong")
	errNoAuthMethod        = fmt.Errorf("no auth method")
)

//Connection container for data
type Connection struct {
	conn    net.Conn
	auth    Authorizer
	ver     byte
	methods []byte
}

//Reset use for reset data
func (c *Connection) Reset() {
	c.conn = nil
	c.auth = nil
	c.ver = 0x0
	c.methods = c.methods[:0]
}

//Serve serve connection: handshake + cmd
func (c *Connection) Serve() {
	defer c.conn.Close()

	if err := binary.Read(c.conn, binary.BigEndian, &c.ver); err != nil {
		log.Error().Msgf("Error on read data from connection : %v", err)
		return
	}

	if c.ver != socks5Ver {
		err := errors.New("unsupported protocol version")
		log.Error().Err(err)
		return
	}

	if err := c.handshake(); err != nil {
		log.Error().Err(err)
		return
	}

	err := c.cmd()
	if err != nil {
		log.Error().Msgf("Connection.cmd error: %v", err)
	}
}

func (c *Connection) handshake() (err error) {
	var acceptableAuthMethodSize byte
	if err = binary.Read(c.conn, binary.BigEndian, &acceptableAuthMethodSize); err != nil {
		return
	}

	c.methods = make([]byte, acceptableAuthMethodSize)
	if _, err = io.ReadFull(c.conn, c.methods); err != nil {
		return
	}

	if !c.shouldAuth() && bytes.IndexByte(c.methods, authNoAuth) != -1 {
		return c.handleNoAuth()
	}

	if c.shouldAuth() && bytes.IndexByte(c.methods, authUP) != -1 {
		if err = c.handleAuthUserPassword(); err != nil {
			c.conn.Write([]byte{socks5Ver, authUPFailure})
		}
		c.conn.Write([]byte{socks5Ver, authUPSuccess})
		return
	}

	c.conn.Write([]byte{socks5Ver, authNoMethods})
	return errNoAuthMethod
}

func (c *Connection) shouldAuth() (should bool) {
	return c.auth != nil && c.auth.ShouldAuth()
}

func (c *Connection) handleAuthUserPassword() (err error) {
	if _, err = c.conn.Write([]byte{socks5Ver, authUP}); err != nil {
		return
	}

	var ver byte
	if err = binary.Read(c.conn, binary.BigEndian, &ver); err != nil {
		return
	}

	if ver != 0x01 {
		return errUserPasswordVersion
	}

	var usernameLen byte
	if err = binary.Read(c.conn, binary.BigEndian, &usernameLen); err != nil {
		return
	}

	username := make([]byte, usernameLen)
	if _, err = io.ReadFull(c.conn, username); err != nil {
		return
	}

	var passwordLen byte
	if err = binary.Read(c.conn, binary.BigEndian, &passwordLen); err != nil {
		return
	}

	password := make([]byte, passwordLen)
	if _, err = io.ReadFull(c.conn, password); err != nil {
		return
	}

	if !c.auth.AuthLoginPassword(*(*string)(unsafe.Pointer(&username)), password) {
		return errUserPassword
	}
	return nil
}

func (c *Connection) handleNoAuth() (err error) {
	_, err = c.conn.Write([]byte{socks5Ver, authNoAuth})

	return err
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

	rsp, err := cmd.Fire()
	if err != nil {

		return err
	}

	if rsp != rspSuccess {
		return c.writeErrorResponce(r, rsp)
	}

	return nil
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
