package main

import (
	"flag"

	"github.com/rs/zerolog/log"
	proxy2 "github.com/vigo5190/go-socks5/proxy"
)

func main() {

	port := flag.String("port", "8008", "listen port")
	listenAddr := flag.String("addr", "0.0.0.0", "listen addr")
	flag.Parse()

	//customFormatter := new(log.TextFormatter)
	//customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	//customFormatter.FullTimestamp = true

	//lg := log.New()
	//lg.Formatter = customFormatter

	log.Info().Msg("vigo5190/go-socks5 Started")
	defer log.Info().Msg("vigo5190/go-socks5 Stop")

	proxy := proxy2.Proxy{
		Listen: *listenAddr + ":" + *port,
	}
	proxy.Start()
}
