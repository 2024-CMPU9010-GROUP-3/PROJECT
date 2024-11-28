//go:build public

package handlers

import (
	"errors"
	"net/http"
	"strconv"

	db "github.com/2024-CMPU9010-GROUP-3/magpie/internal/db/public"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/dtos"
	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	resp "github.com/2024-CMPU9010-GROUP-3/magpie/internal/responses"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

func (handler *LocationHistoryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	var userId pgtype.UUID

	userIdPathParam := r.PathValue("id")
	limitQueryParam := r.URL.Query().Get("limit")
	offsetQueryParam := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitQueryParam)
	if err != nil {
		resp.SendError(customErrors.Parameter.InvalidIntError, w)
		return
	}

	offset, err := strconv.Atoi(offsetQueryParam)
	if err != nil {
		resp.SendError(customErrors.Parameter.InvalidIntError, w)
		return
	}

	err = userId.Scan(userIdPathParam)
	if err != nil {
		resp.SendError(customErrors.Parameter.InvalidUUIDError, w)
		return
	}

	getLocationHistoryParams := db.GetLocationHistoryParams{
		Userid: userId,
		Lim:    int32(limit),
		Off:    int32(offset),
	}

	rows, err := PublicDb().GetLocationHistory(*dbCtx, getLocationHistoryParams)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	var entries []dtos.LocationHistoryEntryDto

	for _, row := range rows {
		longlat, err := geojson.Encode(row.Longlat)
		if err != nil {
			resp.SendError(customErrors.Internal.GeoJsonEncodingError.WithCause(err), w)
			return
		}
		dto := dtos.LocationHistoryEntryDto{
			ID: row.ID, Datecreated: row.Datecreated, Amenitytypes: row.Amenitytypes, Longlat: *longlat, Radius: row.Radius,
		}
		entries = append(entries, dto)
	}

	resp.SendResponse(dtos.ResponseContentDto{Content: dtos.GetLocationHistoryListDto{HistoryEntries: entries, NextOffset: int32(offset + limit)}, HttpStatus: http.StatusAccepted}, w)
}

func (handler *LocationHistoryHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	var idListDto dtos.IntIdListDto
	var userId pgtype.UUID

	userIdPathParam := r.PathValue("id")

	err := userId.Scan(userIdPathParam)
	if err != nil {
		resp.SendError(customErrors.Parameter.InvalidUUIDError, w)
		return
	}

	err = idListDto.Decode(r.Body)
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

	err = PublicDb().DeleteLocationHistoryEntries(*dbCtx, idListDto.IdList)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}
	resp.SendResponse(dtos.ResponseContentDto{Content: struct {
		Deleted bool `json:"deleted"`
	}{Deleted: true}, HttpStatus: http.StatusAccepted}, w)
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

	id, err := PublicDb().CreateLocationHistoryEntry(*dbCtx, createLocationHistoryEntryParam)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	resp.SendResponse(dtos.ResponseContentDto{Content: dtos.IntIdDto{Id: id}, HttpStatus: http.StatusAccepted}, w)
}
