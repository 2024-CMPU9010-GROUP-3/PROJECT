//go:build public

package handlers

import (
	"net/http"

	db "github.com/2024-CMPU9010-GROUP-3/magpie/internal/db/public"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/dtos"
	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	resp "github.com/2024-CMPU9010-GROUP-3/magpie/internal/responses"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twpayne/go-geom"
)

func (handler *LocationHistoryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	return
}

func (handler *LocationHistoryHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	return
}

func (handler *LocationHistoryHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var historyEntryDto dtos.CreateLocationHistoryEntryDto
	var userId pgtype.UUID

	userIdPathParam := r.PathValue("id")

	err := userId.Scan(userIdPathParam)
	if err != nil {
		resp.SendError(customErrors.Parameter.InvalidUUIDError, w)
		return
	}

	err = historyEntryDto.Decode(r.Body)
	if err != nil {
		customError, ok := err.(customErrors.CustomError)
		if !ok {
			resp.SendError(customErrors.Internal.UnknownError.WithCause(err), w)
			return
		} else {
			resp.SendError(customError, w)
			return
		}
	}

	geometry, err := historyEntryDto.Longlat.Decode()
	if err != nil {
		resp.SendError(customErrors.Payload.InvalidPayloadPointError.WithCause(err), w)
		return
	}

	longlat, ok := geometry.(*geom.Point)
	if !ok {
		resp.SendError(customErrors.Payload.InvalidPayloadPointError, w)
		return
	}

	createLocationHistoryEntryParam := db.CreateLocationHistoryEntryParams{
		Userid:       userId,
		Amenitytypes: historyEntryDto.Amenitytypes,
		Longlat:      longlat,
		Radius:       historyEntryDto.Radius,
	}

	id, err := db.New(dbConn).CreateLocationHistoryEntry(*dbCtx, createLocationHistoryEntryParam)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	resp.SendResponse(dtos.ResponseContentDto{Content: dtos.IntIdDto{Id: id}, HttpStatus: http.StatusAccepted}, w)
}
