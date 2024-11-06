//go:build public

package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	db "github.com/2024-CMPU9010-GROUP-3/magpie/internal/db/public"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/dtos"
	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	resp "github.com/2024-CMPU9010-GROUP-3/magpie/internal/responses"
	"github.com/jackc/pgx/v5"
	"github.com/twpayne/go-geom/encoding/geojson"
)

const floatPrecision = 64

func (p *PointsHandler) HandleGetByRadius(w http.ResponseWriter, r *http.Request) {
	// parameters long, lat, radius are required
	params := r.URL.Query()
	long, err_long := strconv.ParseFloat(params.Get("long"), floatPrecision)
	lat, err_lat := strconv.ParseFloat(params.Get("lat"), floatPrecision)
	radius, err_radius := strconv.ParseFloat(params.Get("radius"), floatPrecision)
	typesString := params.Get("types")

	// bad request if any parameters can't be parsed to float
	if err_long != nil || err_lat != nil || err_radius != nil {
		resp.SendError(customErrors.Parameter.InvalidFloatError, w)
		return
	}

	var types []db.PointType
	if len(typesString) != 0 {
		typesSplit := strings.Split(typesString, ",")
		for _, t := range typesSplit {
			parsedType := db.PointType(t)
			if parsedType.IsValid() {
				types = append(types, parsedType)
			} else {
				resp.SendError(customErrors.Parameter.InvalidPointTypeError.WithCause(fmt.Errorf("Type %s is not supported", t)), w)
				return
			}
		}
	}

	points, err := db.New(dbConn).GetPointsInRadius(*dbCtx, db.GetPointsInRadiusParams{Latitude: lat, Longitude: long, Radius: radius, Types: types})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
			return
		}
	}

	pointDtos := []dtos.GetPointDto{}

	for _, p := range points {
		longlat, err := geojson.Encode(p.Longlat)
		if err != nil {
			resp.SendError(customErrors.Internal.GeoJsonEncodingError.WithCause(err), w)
			return
		} else {
			pointDtos = append(pointDtos, dtos.GetPointDto{
				Id:      p.ID,
				Longlat: *longlat,
				Type:    p.Type,
			})
		}
	}
	resp.SendResponse(dtos.ResponseContentDto{Content: pointDtos, HttpStatus: http.StatusOK}, w)
}

func (p *PointsHandler) HandleGetPointDetails(w http.ResponseWriter, r *http.Request) {
	pointIdPathParam := r.PathValue("id")
	pointId, err := strconv.ParseInt(pointIdPathParam, 10, 64)

	// bad request if id can't be parsed to int
	if err != nil {
		resp.SendError(customErrors.Parameter.InvalidIntError, w)
		return
	}

	pointDetailsBytes, err := db.New(dbConn).GetPointDetails(*dbCtx, pointId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			resp.SendError(customErrors.NotFound.PointNotFoundError, w)
			return
		} else {
			resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
			return
		}
	}

	// check if bytes received from database are valid json
	var decodedDetails json.RawMessage
	err = json.Unmarshal(pointDetailsBytes, &decodedDetails)
	if err != nil {
		resp.SendError(customErrors.Internal.JsonDecodingError, w)
		return
	}

	resp.SendResponse(dtos.ResponseContentDto{Content: decodedDetails, HttpStatus: http.StatusOK}, w)
}
