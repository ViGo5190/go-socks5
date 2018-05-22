package proxy

import (
	"github.com/rs/zerolog/log"
)

//Proxy is main struct for all magic
type Proxy struct {
	Listen string
}

//Start method for start proxy
func (p *Proxy) Start() {
	server := &Server{}

	log.Info().Msgf("start %s", p.Listen)
	server.ListenAndServe("tcp", p.Listen)
}
