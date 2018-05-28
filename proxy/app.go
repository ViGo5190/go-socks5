package proxy

import (
	"github.com/rs/zerolog/log"
)

//Proxy is main struct for all magic
type Proxy struct {
	Listen     string `toml:"listen"`
	AuthEnable bool   `toml:"auth"`
	Users      []user `toml:"users"`
}

type user struct {
	Login string `toml:"login"`
	Pass  string `toml:"pass"`
}

//Start method for start proxy
func (p *Proxy) Start() {

	server := &Server{
		auth: p.authorizer(),
	}

	log.Info().Msgf("start %s", p.Listen)
	server.ListenAndServe("tcp", p.Listen)
}

func (p *Proxy) authorizer() Authorizer {
	upmap := map[string]string{}
	for _, u := range p.Users {
		upmap[u.Login] = u.Pass
	}
	return &Auth{
		Users:      upmap,
		AuthEnable: p.AuthEnable,
	}
}
