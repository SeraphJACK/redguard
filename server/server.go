package server

import (
	"flag"
	"log"
	"net/http"
)

var listen = flag.String("listen", "0.0.0.0:8081", "HTTP Server listen address")

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/completion", handleCompletion)

	log.Fatal(http.ListenAndServe(*listen, mux))
}
