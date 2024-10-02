//go:build public

package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	db "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/db/public"
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/util"
)

func (p *PointsHandler) HandleGetByRadius(w http.ResponseWriter, r *http.Request) {
	points, err := dbQueries.GetPointsInEnvelope(*dbCtx, db.GetPointsInEnvelopeParams{X1: -1, Y1: -1, X2: 1, Y2: 1})
	if err != nil {
		log.Printf("Could not get points from database: %v\n", err)
	}
	log.Printf("%v\n", points)
	data := util.Placeholder("GET ByRadius")
	w.Header().Set(contentType, applicationJson)
	err = json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}

func (p *PointsHandler) HandleGetPointDetails(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("GET PointDetails")
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}
