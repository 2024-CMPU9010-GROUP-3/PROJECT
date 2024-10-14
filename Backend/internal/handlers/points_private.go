//go:build private

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	db "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/db/private"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	// geos "github.com/twpayne/go-geom"
)

type PointDto struct {
	Longlat geojson.Geometry `json:"longlat"`
	Type    string           `json:"type"`
	Details any              `json:"details"` // potentially unsafe, but we need to accept any json object here
}

// HandlePost creates new point
// @Summary      Create new point
// @Tags         Points Private
// @Description  Creates a new point in the database
// @Accept       json
// @Produce      json
// @Success      201
// @Router       /v1/private/points/ [post]
func (p *PointsHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	dbQueries := db.New(dbConn)

	var point PointDto
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

	geometry, err := point.Longlat.Decode()
	if err != nil {
		log.Printf("Could not decode geojson from request: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pt, ok := geometry.(*geom.Point)
	if !ok {
		log.Printf("Could not convert geometry to point")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createPointParams := db.CreatePointParams{
		Longlat: pt,
		Type:    db.PointType(point.Type),
		Details: encodedJson,
	}

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

	// careful: if an ID was given that doesn't exist, no changes are made to the database
	// there is currently no mechanism that notifies the frontend if any actual changes were made

	pointIdPathParam := r.PathValue("id")
	pointId, err := strconv.ParseInt(pointIdPathParam, 10, 64)

	// bad request if id can't be parsed to int
	if err != nil {
		log.Printf("Invalid path parameter: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbQueries := db.New(dbConn)

	var point PointDto
	err = json.NewDecoder(r.Body).Decode(&point)
	if err != nil {
		log.Printf("Could not decode request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(point.Type) == 0 {
		log.Println("Field \"Type\" cannot be empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encodedJson, err := json.Marshal(point.Details)
	if err != nil {
		log.Printf("Could not encode point details: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	geometry, err := point.Longlat.Decode()
	if err != nil {
		log.Printf("Could not decode geojson from request: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pt, ok := geometry.(*geom.Point)
	if !ok {
		log.Printf("Could not convert geometry to point")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatePointParams := db.UpdatePointParams{
		ID:      pointId,
		Longlat: pt,
		Type:    db.PointType(point.Type),
		Details: encodedJson,
	}

	err = dbQueries.UpdatePoint(*dbCtx, updatePointParams)
	if err != nil {
		log.Printf("Could not update point in database: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Executed update query on databse for point id: %+v\n", pointId)

	w.WriteHeader(http.StatusOK)
}

func (p *PointsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {

	pointIdPathParam := r.PathValue("id")
	pointId, err := strconv.ParseInt(pointIdPathParam, 10, 64)

	// bad request if id can't be parsed to int
	if err != nil {
		log.Printf("Invalid path parameter: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbQueries := db.New(dbConn)

	err = dbQueries.DeletePoint(*dbCtx, pointId)
	if err != nil {
		log.Printf("Could not delete point from database: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Executed delete query on databse for point id: %+v\n", pointId)

	w.WriteHeader(http.StatusNoContent)
}
