package main

import (
	"broker/cmd/api"
	"fmt"
	"log"
	"net/http"
)

const webPort = "8010"

func main() {
	app := &api.App{
		Cfg: &api.Config{},
	}

	log.Printf("Starting broker service on port: %s\n", webPort)

	err := http.ListenAndServe(fmt.Sprintf(":%s", webPort), app.Routes())
	if err != nil {
		log.Panic(err)
	}
}
