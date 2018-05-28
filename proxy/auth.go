package proxy

import (
	"crypto/sha1"
	"encoding/base64"
	"sync"
	"hash"
)

type Auther interface {
	Auth(login, pass string) bool
}

type Auth struct {
	AuthEnable bool
	Users      map[string]string
	Shapool    sync.Pool
}

func (p *Auth) Auth(login, pass string) bool {
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

func (p *Auth) hashSha(password string) string {
	s := p.Shapool.Get()
	if s == nil {
		s = sha1.New()
	}

	ss := s.(hash.Hash)

	ss.Write([]byte(password))
	passwordSum := []byte(ss.Sum(nil))
	ss.Reset()
	p.Shapool.Put(s)
	return base64.StdEncoding.EncodeToString(passwordSum)
}
