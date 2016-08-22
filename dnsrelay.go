package main

import (
	"./dnsrelay"
	"flag"
	"os"
	"os/signal"
	"syscall"
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

	ds, err := dnsrelay.NewDNSServer(config)
	if err != nil {
		panic(err)
	}

	if err := ds.Listen(); err != nil {
		panic(err)
	}

	// Wait for SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}