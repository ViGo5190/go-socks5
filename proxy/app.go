package proxy

import (
	log "github.com/sirupsen/logrus"
)

//Proxy is main struct for all magic
type Proxy struct {
	Listen string
	Log    *log.Logger
}

//Start method for start proxy
func (p *Proxy) Start() {
	server := &Server{Logger: p.Log}

	log.Infof("start %s", p.Listen)
	server.ListenAndServe("tcp", p.Listen)
}
