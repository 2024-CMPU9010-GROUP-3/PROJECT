//go:build private

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "github.com/2024-CMPU9010-GROUP-3/magpie/internal/db/private"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/dtos"
	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	resp "github.com/2024-CMPU9010-GROUP-3/magpie/internal/responses"
	"github.com/twpayne/go-geom"
)

func (p *PointsHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	dbQueries := db.New(dbConn)

	var point dtos.CreatePointDto

  err := point.Decode(r.Body)
  if err != nil {
		e, ok := err.(customErrors.CustomError)
		if ok {
			resp.SendError(e, w)
			return
		} else {
			resp.SendError(customErrors.Internal.UnknownError.WithCause(err), w)
			return
		}
	}

	// cannot result in error if decoding above was successful
	encodedJson, _ := json.Marshal(point.Details)

	geometry, err := point.Longlat.Decode()
	if err != nil {
		resp.SendError(customErrors.Payload.InvalidPayloadPointError.WithCause(err), w)
		return
	}

	pt, ok := geometry.(*geom.Point)
	if !ok {
		resp.SendError(customErrors.Payload.InvalidPayloadPointError, w)
		return
	}

	pointType := db.PointType(point.Type)
	if !pointType.IsValid() {
		resp.SendError(customErrors.Payload.InvalidPayloadPointError, w)
		return
	}

	createPointParams := db.CreatePointParams{
		Longlat: pt,
		Type:    db.PointType(point.Type),
		Details: encodedJson,
	}

	pointId, err := dbQueries.CreatePoint(*dbCtx, createPointParams)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	resp.SendResponse(dtos.ResponseContentDto{Content: dtos.PointIdDto{Id: pointId}, HttpStatus: http.StatusCreated}, w)
}

func (p *PointsHandler) HandlePut(w http.ResponseWriter, r *http.Request) {

	// careful: if an ID was given that doesn't exist, no changes are made to the database
	// there is currently no mechanism that notifies the frontend if any actual changes were made

	pointIdPathParam := r.PathValue("id")
	pointId, err := strconv.ParseInt(pointIdPathParam, 10, 64)

	// bad request if id can't be parsed to int
	if err != nil {
		resp.SendError(customErrors.Parameter.InvalidIntError, w)
		return
	}

	dbQueries := db.New(dbConn)

	var point dtos.CreatePointDto
	err = point.Decode(r.Body)
  if err != nil {
		e, ok := err.(customErrors.CustomError)
		if ok {
			resp.SendError(e, w)
			return
		} else {
			resp.SendError(customErrors.Internal.UnknownError.WithCause(err), w)
			return
		}
	}

	// cannot result in error if decoding above was successful
	encodedJson, _ := json.Marshal(point.Details)

	geometry, err := point.Longlat.Decode()
	if err != nil {
		resp.SendError(customErrors.Payload.InvalidPayloadPointError.WithCause(err), w)
		return
	}

	pt, ok := geometry.(*geom.Point)
	if !ok {
		resp.SendError(customErrors.Payload.InvalidPayloadPointError, w)
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
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	resp.SendResponse(dtos.ResponseContentDto{Content: dtos.PointIdDto{Id: pointId}, HttpStatus: http.StatusAccepted}, w)
}

func (p *PointsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {

	pointIdPathParam := r.PathValue("id")
	pointId, err := strconv.ParseInt(pointIdPathParam, 10, 64)

	// bad request if id can't be parsed to int
	if err != nil {
		resp.SendError(customErrors.Parameter.InvalidIntError, w)
		return
	}

	dbQueries := db.New(dbConn)

	err = dbQueries.DeletePoint(*dbCtx, pointId)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	resp.SendResponse(dtos.ResponseContentDto{Content: dtos.PointIdDto{Id: pointId}, HttpStatus: http.StatusAccepted}, w)
}
