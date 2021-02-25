package main

import (
	"flag"
	"gitlab.com/spacecowboy/paymail-server/config"
	"gitlab.com/spacecowboy/paymail-server/v1"
	"log"
	"net/http"
)

func main() {
	configPath := flag.String("config", "config.toml", "path to config file")

	flag.Parse()

	config, err := config.ReadConfig(*configPath)

	if err != nil {
		log.Fatalln(err)
	}

	wellknown, err := v1.GetWellKnownBsvAliasResponse(config)

	if err != nil {
		log.Fatalln(err)
	}

	handler := v1.GetHandler(config)

	log.Printf("%s/.well-known/bsvalias", config.Server.ListenAddress)
	log.Println(string(wellknown))
	log.Fatalln(http.ListenAndServe(config.Server.ListenAddress, handler))
}
