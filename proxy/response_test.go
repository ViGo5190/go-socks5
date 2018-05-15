package proxy

import (
	"testing"
)

func TestParseAddrIpV4(t *testing.T) {
	r := &Rsps{
		socksVer: socks5Ver,
		rsp:      rspSuccess,
		reserved: reservedSymbol,
		addr:     []byte{},
	}
	if err := r.parseAddr("8.8.8.8:88"); err != nil {
		t.Error(err)
	}

	if string(r.addr) != string([]byte{8, 8, 8, 8}) {
		t.Errorf("Should r.addr=[8 8 8 8], got %v", r.addr)
	}

	if r.addressType != addressTypeIPv4 {
		t.Errorf("Should r.addressType=1, got %v", r.addressType)
	}

	if r.port != 88 {
		t.Errorf("Should r.port=88, got %v", r.port)
	}
}

func TestParseAddrIpV6(t *testing.T) {
	r := &Rsps{
		socksVer: socks5Ver,
		rsp:      rspSuccess,
		reserved: reservedSymbol,
		addr:     []byte{},
	}
	if err := r.parseAddr("[2001:db8:0:0:0:0:2:1]:8888"); err != nil {
		t.Error(err)
	}

	if string(r.addr) != string([]byte{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 1}) {
		t.Errorf("Should r.addr=[32 1 13 184 0 0 0 0 0 0 0 0 0 2 0 1], got %v", r.addr)
	}

	if r.addressType != addressTypeIPv6 {
		t.Errorf("Should r.addressType=1, got %v", r.addressType)
	}

	if r.port != 8888 {
		t.Errorf("Should r.port=8888, got %v", r.port)
	}
}

func TestSplitHostError(t *testing.T) {
	r := &Rsps{
		socksVer: socks5Ver,
		rsp:      rspSuccess,
		reserved: reservedSymbol,
		addr:     []byte{},
	}
	err := r.parseAddr("vig.gs")
	if err.Error() != "address vig.gs: missing port in address" {
		t.Errorf("Should 'address vig.gs: missing port in address', got '%v'", err.Error())
	}
}

func TestSplitHostError2(t *testing.T) {
	r := &Rsps{
		socksVer: socks5Ver,
		rsp:      rspSuccess,
		reserved: reservedSymbol,
		addr:     []byte{},
	}
	err := r.parseAddr("vig.gs:80:0")
	if err.Error() != "address vig.gs:80:0: too many colons in address" {
		t.Errorf("Should 'address vig.gs:80:0: too many colons in address', got '%v'", err.Error())
	}
}

func TestParseIPEmptyIp(t *testing.T) {
	r := &Rsps{
		socksVer: socks5Ver,
		rsp:      rspSuccess,
		reserved: reservedSymbol,
		addr:     []byte{},
	}
	err := r.parseAddr("ya.ru:80")
	if err != errRspEmptyIP {
		t.Errorf("Should '%v', got '%v'", errRspEmptyIP, err.Error())
	}
}

func TestParsePortErr(t *testing.T) {
	r := &Rsps{
		socksVer: socks5Ver,
		rsp:      rspSuccess,
		reserved: reservedSymbol,
		addr:     []byte{},
	}
	err := r.parseAddr("8.8.8.8:foo")
	if err != errRspEmptyPort {
		t.Errorf("Should '%v', got '%v'", errRspEmptyPort, err.Error())
	}
}
