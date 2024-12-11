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

	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%s", webPort),
		Handler: app.Routes(),
	}

	log.Printf("Starting broker service on port: %s\n", webPort)

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
