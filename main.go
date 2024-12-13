package main

import (
	"api/config"
	"api/server"
	"log"
	"net"
	"net/http"
)

func main() {
	log.Printf("Serveur démarré sur http://%s\n", net.JoinHostPort(config.App.Server.URL, config.App.Server.Port))
	log.Fatal(http.ListenAndServe(net.JoinHostPort(config.App.Server.URL, config.App.Server.Port), server.ServeMux()))
}
