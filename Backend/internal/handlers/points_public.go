//go:build public

package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	db "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/db/public"
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/twpayne/go-geom/encoding/geojson"
)

const floatPrecision = 32

type pointDto struct {
	Id      int64
	Longlat geojson.Geometry
	Type    db.PointType
}

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

	// construct envelope
	x1 := long - radius
	y1 := lat - radius
	x2 := long + radius
	y2 := lat + radius

	points, err := db.New(dbConn).GetPointsInEnvelope(*dbCtx, db.GetPointsInEnvelopeParams{X1: x1, Y1: y1, X2: x2, Y2: y2})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			log.Printf("Could not get points from database, unknown error: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	pointDtos := []pointDto{}

	for _, p := range points {
		longlat, err := geojson.Encode(p.Longlat)
		if err != nil {
			log.Printf("Could not encode point: %+v\n", err)
		} else {
			pointDtos = append(pointDtos, pointDto{
				p.ID,
				*longlat,
				p.Type,
			})
		}
	}
	err = json.NewEncoder(w).Encode(pointDtos)
	util.CheckResponseError(err, w)
}

func (p *PointsHandler) HandleGetPointDetails(w http.ResponseWriter, r *http.Request) {
	pointIdPathParam := r.PathValue("id")
	pointId, err := strconv.ParseInt(pointIdPathParam, 10, 64)

	// bad request if id can't be parsed to int
	if err != nil {
		log.Printf("Invalid path parameter: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pointDetailsBytes, err := db.New(dbConn).GetPointDetails(*dbCtx, pointId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			log.Printf("Could not get point details from database: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// check if bytes received from database are valid json
	var decodedDetails json.RawMessage
	err = json.Unmarshal(pointDetailsBytes, &decodedDetails)
	if err != nil {
		log.Printf("Received invalid json from database: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(decodedDetails)
	if err != nil {
		log.Printf("Could not send bearer token as response: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
