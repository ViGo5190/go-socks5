package proxy

import (
	"crypto/sha1"
	"encoding/base64"
	"hash"
	"sync"
)

//Authorizer interface for auth
type Authorizer interface {
	AuthLoginPassword(login string, pass []byte) bool
	ShouldAuth() bool
}

//Auth container for auth info
type Auth struct {
	AuthEnable bool
	Users      map[string]string
	shapool    sync.Pool
}

//ShouldAuth return true if aith enable
func (p *Auth) ShouldAuth() bool {
	return p.AuthEnable
}

//AuthLoginPassword return true if auth enable and login/password corrected
func (p *Auth) AuthLoginPassword(login string, pass []byte) bool {
	if !p.AuthEnable {
		return true
	}

	spass, ok := p.Users[login]

	if !ok {
		return false
	}

	hashedPass := p.hashSha(pass)

	return hashedPass == spass
}

func (p *Auth) hashSha(password []byte) string {
	s := p.shapool.Get()
	if s == nil {
		s = sha1.New()
	}

	ss := s.(hash.Hash)
	defer ss.Reset()
	defer p.shapool.Put(ss)

	ss.Write(password)
	passwordSum := ss.Sum(nil)

	return base64.StdEncoding.EncodeToString(passwordSum)
}
