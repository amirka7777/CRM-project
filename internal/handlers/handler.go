package handlers

import "net/http"

type HandlerSupport struct{}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{"message": "pong"}`))
}
