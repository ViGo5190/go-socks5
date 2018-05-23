package main

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"

	proxy2 "github.com/vigo5190/go-socks5/proxy"
)

func main() {
	log.Info().Msg("vigo5190/go-socks5 Started")
	defer log.Info().Msg("vigo5190/go-socks5 Stop")

	cfgFile := flag.String("c", "config.toml", "config file")

	flag.Parse()

	proxy := proxy2.Proxy{}

	if _, err := toml.DecodeFile(*cfgFile, &proxy); err != nil {
		log.Error().Err(err)
		panic(err)
	}

	proxy.Start()
}
