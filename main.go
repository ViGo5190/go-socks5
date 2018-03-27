package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vigo5190/go-socks5/proxy"
)

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true

	lg := log.New()
	lg.Formatter = customFormatter

	lg.Info("eApi Service Start")
	defer lg.Info("eApi Service Stop")

	lg.SetLevel(log.DebugLevel)

	server := proxy.New(lg)

	server.ListenAndServe("tcp", "127.0.0.1:8008")

}
