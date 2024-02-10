package server

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

var listen = flag.String("listen", "0.0.0.0:8081", "HTTP Server listen address")

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/completion", handleCompletion)
	mux.HandleFunc("/api/packet", handleGetRedPacket)

	log.Fatal(http.ListenAndServe(*listen, mux))
}

func handleGetRedPacket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	code := r.URL.Query().Get("code")
	if code != RedPacketCode {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(RedPacketResult)
}
