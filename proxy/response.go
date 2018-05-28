package proxy

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

//Rsps container for data
type Rsps struct {
	socksVer    byte
	rsp         byte
	reserved    byte
	addressType byte

	addr []byte
	port uint16
}

var (
	//errRspServerMsg            = errors.New("general error")
	errRspCommandNotSupportMsg = errors.New("command not support")
	errRspAddressToShort       = errors.New("address to short")
	errRspAddressToLong        = errors.New("address to long")
	errRspEmptyIP              = errors.New("empty ip")
	errRspEmptyPort            = errors.New("empty port")
)

//Bytes bytes of response
func (r *Rsps) Bytes() (buf []byte, err error) {
	if r.socksVer != socks5Ver {
		err = errRspCommandNotSupportMsg
		return
	}

	buf = make([]byte, 0, net.IPv6len+8)

	buf = append(buf, r.socksVer, r.rsp, r.reserved, r.addressType)
	switch r.addressType {
	case addressTypeIPv4:
		if len(r.addr) < net.IPv4len {
			err = errRspAddressToShort
			return
		}
		buf = append(buf, r.addr[:net.IPv4len]...)
	case addressTypeDomain:
		if len(r.addr) > 255 {
			err = errRspAddressToLong
			return
		}
		buf = append(buf, byte(len(r.addr)))
		buf = append(buf, r.addr...)
	case addressTypeIPv6:
		if len(r.addr) < net.IPv6len {
			err = errRspAddressToShort
			return
		}
		buf = append(buf, r.addr[:net.IPv6len]...)
	}
	buf = append(buf, 0, 0)

	binary.BigEndian.PutUint16(buf[len(buf)-2:], r.port)

	return buf, nil
}

func (r *Rsps) parseAddr(addr string) (err error) {

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return
	}

	ip := net.ParseIP(host)

	if ip == nil {
		return errRspEmptyIP
	}

	if ipv4 := ip.To4(); ipv4 != nil {
		r.addressType = addressTypeIPv4
		r.addr = ipv4[:net.IPv4len]
	} else {
		r.addressType = addressTypeIPv6
		r.addr = ip[:net.IPv6len]
	}

	prt, err := strconv.Atoi(port)
	if err != nil {
		return errRspEmptyPort
	}
	r.port = uint16(prt)

	return nil
}
