package handlers

import "net/http"

var IsAlive func() bool

func HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	if IsAlive == nil || !IsAlive() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("false"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
		return
	}
}
