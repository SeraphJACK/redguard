package server

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"

	"git.s8k.top/SeraphJACK/redguard/slog"
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

	slog.Log(calcRealIP(r), "GetRedPacket", code, "")

	if code != RedPacketCode {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(RedPacketResult)
}

func calcRealIP(r *http.Request) string {
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	addr := r.RemoteAddr
	colonIdx := strings.LastIndex(addr, ":")
	if colonIdx > 0 {
		addr = addr[:colonIdx]
	}
	return addr
}
