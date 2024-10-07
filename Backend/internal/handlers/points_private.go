//go:build private

package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	db "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/db/private"
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/util"
	geos "github.com/twpayne/go-geos/geometry"
)

type PointWrapper struct {
	Longlat *geos.Geometry `json:"longlat"`
	Type    string         `json:"type"`
	Details any            `json:"details"` // potentially unsafe, but we need to accept any json object here
}

func (p *PointsHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	dbQueries := db.New(dbConn)

	var point PointWrapper
	err := json.NewDecoder(r.Body).Decode(&point)
	if err != nil {
		log.Printf("Could not decode request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encodedJson, err := json.Marshal(point.Details)
	if err != nil {
		log.Printf("Could not encode point details: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createPointParams := db.CreatePointParams{
		Longlat: point.Longlat.Geom,
		Type:    db.PointType(point.Type),
		Details: encodedJson,
	}

	log.Printf("Decoded point: %+v\n", point)
	log.Printf("Create point params: %+v\n", createPointParams)

	count, err := dbQueries.CreatePoint(*dbCtx, createPointParams)
	if err != nil {
		log.Printf("Could not save point to database: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Saved point to databse: %+v\n", count)

	w.WriteHeader(http.StatusCreated)
}

func (p *PointsHandler) HandlePut(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("PUT points")
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}

func (p *PointsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("DELETE points")
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}
