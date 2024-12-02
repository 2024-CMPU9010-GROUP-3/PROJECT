//go:build public

package handlers

import (
	"errors"
	"net/http"

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

	err := userId.Scan(userIdPathParam)
	if err != nil {
		resp.SendError(customErrors.Parameter.InvalidUUIDError, w)
		return
	}

	rows, err := db.New(dbConn).GetLocationHistory(*dbCtx, userId)
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

		typeRows, err := db.New(dbConn).GetAmenityTypeCount(*dbCtx, row.ID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
			return
		}
		var amenityCounts []dtos.AmenityTypeWithCount

		for _, typeRow := range typeRows {
			amenityCounts = append(amenityCounts, dtos.AmenityTypeWithCount{AmenityType: typeRow.Type, Count: int(typeRow.Count)})
		}

		dto := dtos.LocationHistoryEntryDto{
			ID: row.ID, Datecreated: row.Datecreated, Amenitytypes: amenityCounts, Longlat: *longlat, Radius: row.Radius, DisplayName: row.Displayname.String,
		}
		entries = append(entries, dto)
	}

	resp.SendResponse(dtos.ResponseContentDto{Content: dtos.GetLocationHistoryListDto{HistoryEntries: entries}, HttpStatus: http.StatusOK}, w)
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

	err = db.New(dbConn).DeleteLocationHistoryEntries(*dbCtx, idListDto.IdList)
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
		Userid:      userId,
		Longlat:     longlat,
		Radius:      historyEntryDto.Radius,
		Displayname: historyEntryDto.DisplayName,
	}

	tx, err := dbConn.Begin(*dbCtx)
	if err != nil {
		resp.SendError(customErrors.Database.TransactionStartError, w)
		return
	}
	defer func() {
		// potential error from rollback is not fatal, ignoring for now
		if tx != nil {
			_ = tx.Rollback(*dbCtx)
		}
	}()

	id, err := db.New(dbConn).WithTx(tx).CreateLocationHistoryEntry(*dbCtx, createLocationHistoryEntryParam)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	for _, entry := range historyEntryDto.Amenitytypes {
		err = db.New(dbConn).WithTx(tx).CreateAmenityCountEntry(*dbCtx, db.CreateAmenityCountEntryParams{Historyentryid: id, Type: entry.AmenityType, Count: int32(entry.Count)})
		if err != nil {
			resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
			return
		}
	}

	err = tx.Commit(*dbCtx)
	if err != nil {
		resp.SendError(customErrors.Database.TransactionCommitError, w)
		return
	}

	resp.SendResponse(dtos.ResponseContentDto{Content: dtos.IntIdDto{Id: id}, HttpStatus: http.StatusAccepted}, w)
}
