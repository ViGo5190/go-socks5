package proxy

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
)

const (
	addressTypeIPv4   = 0x01
	addressTypeDomain = 0x03
	addressTypeIPv6   = 0x04

	reservedSymbol = 0x00
)

var (
	//ErrorAddressTypeNotSupported address type not supported
	ErrorAddressTypeNotSupported = errors.New("address type not supported")

	//ErrorWrongReservedSymbol wrong reserved symbol
	ErrorWrongReservedSymbol = errors.New("wrong reserved symbol")
)

//Rqst container for data
type Rqst struct {
	socksVer    byte
	cmd         byte
	reserved    byte
	addressType byte

	addr []byte
	port uint16
}

//FQDN return fqdn for given data
func (r *Rqst) FQDN() string {
	var host string
	switch r.addressType {
	case addressTypeIPv4:
		host = net.IPv4(r.addr[0], r.addr[1], r.addr[2], r.addr[3]).String()
	case addressTypeDomain:
		host = string(r.addr)
	case addressTypeIPv6:
		host = net.IP(r.addr).String()
	default:
		host = "<unsupported address type>"
	}
	return host + ":" + strconv.Itoa(int(r.port))
}

func (r *Rqst) fromReader(src io.Reader) (err error) {
	if err = binary.Read(src, binary.BigEndian, &r.socksVer); err != nil {
		return
	}

	if err = binary.Read(src, binary.BigEndian, &r.cmd); err != nil {
		return
	}

	if err = binary.Read(src, binary.BigEndian, &r.reserved); err != nil {
		return
	}

	if r.reserved != reservedSymbol {
		return ErrorWrongReservedSymbol
	}

	if err = binary.Read(src, binary.BigEndian, &r.addressType); err != nil {
		return
	}

	var ln byte
	switch r.addressType {
	case addressTypeIPv4:
		ln = net.IPv4len
	case addressTypeDomain:
		if err = binary.Read(src, binary.BigEndian, &ln); err != nil {
			return
		}
	case addressTypeIPv6:
		ln = net.IPv6len
	default:
		return ErrorAddressTypeNotSupported
	}

	r.addr = make([]byte, ln)

	if _, err = io.ReadFull(src, r.addr); err != nil {
		return
	}

	if err = binary.Read(src, binary.BigEndian, &r.port); err != nil {
		return
	}

	return nil
}
