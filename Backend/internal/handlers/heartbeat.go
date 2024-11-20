package handlers

import (
	"context"
	"net/http"
	"time"
)

func isAlive() bool {
	// if the database doesn't respond within 5s something is wrong
	ctx, cancel := context.WithTimeout(*dbCtx, time.Duration(5*time.Second))
	defer cancel()
	_, err := dbConn.Exec(ctx, ";")
	return err == nil
}

func HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	if !isAlive() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("false"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
		return
	}
}
