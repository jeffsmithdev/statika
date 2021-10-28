package tasks

import (
	"log"
	"net/http"
	"statika/models"
)

func Server(cfg *models.Config) {
	fs := http.FileServer(http.Dir(cfg.OutputDir))
	http.Handle("/", fs)

	var hostname, host, port string

	if host = cfg.DevServerHost; host == "" {
		host = "localhost"
	}

	if port = cfg.DevServerPort; port == "" {
		port = "8001"
	}

	hostname = host + ":" + port

	log.Println("Listening on http://" + hostname)
	err := http.ListenAndServe(hostname, nil)
	if err != nil {
		log.Fatal(err)
	}
}
