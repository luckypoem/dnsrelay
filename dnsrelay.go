package main

import (
	"./dnsrelay"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"github.com/labstack/gommon/log"
)

var (
	config_path = flag.String("config", "dnsrelay.toml", "Config file")
)

func main() {
	// this is your domain. All records will be scoped under it, e.g.,
	// 'test.docker' below.

	if *config_path == "" {
		panic("Arguments missing")
	}

	config, err := dnsrelay.NewConfig(*config_path)
	if err != nil {
		panic(err)
	}

	_, err = dnsrelay.NewDNSServer(config, nil)
	if err != nil {
		panic(err)
	}
	log.Printf("Server started!\n")

	// Wait for SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}