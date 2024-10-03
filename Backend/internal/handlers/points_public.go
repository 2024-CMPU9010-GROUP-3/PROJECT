//go:build public

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	db "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/db/public"
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/util"
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/wrappers"
)

const floatPrecision = 32

func (p *PointsHandler) HandleGetByRadius(w http.ResponseWriter, r *http.Request) {
	// parameters long, lat, radius are required
	params := r.URL.Query()
	long, err_long := strconv.ParseFloat(params.Get("long"), floatPrecision)
	lat, err_lat := strconv.ParseFloat(params.Get("lat"), floatPrecision)
	radius, err_radius := strconv.ParseFloat(params.Get("radius"), floatPrecision)

	// bad request if any parameters can't be parsed to float
	if err_long != nil || err_lat != nil || err_radius != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	x1 := long - radius
	y1 := lat - radius
	x2 := long + radius
	y2 := lat + radius

	points, err := dbQueries.GetPointsInEnvelope(*dbCtx, db.GetPointsInEnvelopeParams{X1: x1, Y1: y1, X2: x2, Y2: y2})
	if err != nil {
		log.Printf("Could not get points from database: %v\n", err)
	}
	w.Header().Set(contentType, applicationJson)
	err = json.NewEncoder(w).Encode(wrappers.FromRowList(points))
	util.CheckResponseError(err, w)
}

func (p *PointsHandler) HandleGetPointDetails(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("GET PointDetails")
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}
