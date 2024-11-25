package handlers

import (
	"context"
	"io"
	"net/http"
	"time"
)

func isAlive() bool {
	// if the database doesn't respond within 5s something is wrong
	ctx, cancel := context.WithTimeout(*dbCtx, time.Duration(5*time.Second))
	defer cancel()
	_, err := dbConn.Exec(ctx, "-- ping")
	return err == nil
}

func HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	if !isAlive() {
		w.WriteHeader(http.StatusInternalServerError)
		_,_ = io.WriteString(w, "false")
		return
	} else {
		w.WriteHeader(http.StatusOK)
		_,_ =io.WriteString(w, "true")
		return
	}
}
