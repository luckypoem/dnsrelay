package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"log"
	"github.com/FTwOoO/dnsrelay/dns"
)

var (
	config_path = flag.String("config", "dns.toml", "Config file")
)

func main() {
	// this is your domain. All records will be scoped under it, e.g.,
	// 'test.docker' below.

	if *config_path == "" {
		panic("Arguments missing")
	}

	config, err := dns.NewConfig(*config_path)
	if err != nil {
		panic(err)
	}

	_, err = dns.NewDNSServer(config, false)
	if err != nil {
		panic(err)
	}
	log.Print("Server started!\n")

	// Wait for SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}