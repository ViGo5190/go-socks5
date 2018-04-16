package proxy

import (
	log "github.com/sirupsen/logrus"
)

type Proxy struct {
	Listen string
	Log    *log.Logger
}

func (p *Proxy) Start() {
	server := New(p.Log)

	log.Infof("start %s", p.Listen)
	server.ListenAndServe("tcp", p.Listen)
}
