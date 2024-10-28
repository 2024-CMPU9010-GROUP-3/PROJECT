//go:build public

package handlers

import (
	"context"
	"net/http"
	"testing"

	db "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/db/public"
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/errors"
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/util/testutil"
	"github.com/pashagolub/pgxmock/v4"
	go_geom "github.com/twpayne/go-geom"
)

func TestPointsHandlerHandleGetByRadius(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	handler := &PointsHandler{}

	validPoint, err := go_geom.NewPoint(go_geom.XY).SetSRID(4326).SetCoords([]float64{-6.268726783812387, 53.3484472329815})
	if err != nil {
		t.Fatal(err)
	}

	tests := []testutil.HandlerTestDefinition{
		{
			Name:   "Valid input",
			Method: "GET",
			Route:  "/points/inRadius",
			QueryParams: map[string]string{
				"long":   "-6.269925",
				"lat":    "53.345474",
				"radius": "5000",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT Id, LongLat::geometry, Type from points
                      WHERE ST_DWithin\(
                          LongLat::geography,
                          ST_SetSRID\(ST_MakePoint\(\$1::float, \$2::float\), 4326\)::geography,
                          \$3::float
                      \)`).
					WithArgs(
						float64(-6.269925),
						float64(53.345474),
						float64(5000)).
					WillReturnRows(pgxmock.NewRows([]string{"id", "longlat", "type"}).
						AddRow(int64(236), validPoint, db.PointTypeParking))
			},
			ExpectedStatus: http.StatusOK,
			ExpectedJSON: `{
				"error": null,
				"response": {
					"content": [
						{
							"Id": 236,
							"Longlat": {
								"type": "Point",
								"coordinates": [
									-6.268726783812387,
									53.3484472329815
								]
							},
							"Type": "parking"
						}]
					}
				}`,
		},
		{
			Name:   "Non-parsable Longitude",
			Method: "GET",
			Route:  "/points/inRadius",
			QueryParams: map[string]string{
				"long":   "abcd",
				"lat":    "53.345474",
				"radius": "5000",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// Handler should return error before db call is made
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Parameter.InvalidFloatError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1203,
						"errorMsg": "Parameter invalid, expected type Float"
					},
					"response": null
				}`,
		},
		{
			Name:   "Non-parsable Latitude",
			Method: "GET",
			Route:  "/points/inRadius",
			QueryParams: map[string]string{
				"long":   "-6.269925",
				"lat":    "abcd",
				"radius": "5000",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// Handler should return error before db call is made
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Parameter.InvalidFloatError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1203,
						"errorMsg": "Parameter invalid, expected type Float"
					},
					"response": null
				}`,
		},
		{
			Name:   "Non-parsable Radius",
			Method: "GET",
			Route:  "/points/inRadius",
			QueryParams: map[string]string{
				"long":   "-6.269925",
				"lat":    "53.345474",
				"radius": "xyz",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// Handler should return error before db call is made
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Parameter.InvalidFloatError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1203,
						"errorMsg": "Parameter invalid, expected type Float"
					},
					"response": null
				}`,
		},
	}

	testutil.RunTests(t, handler.HandleGetByRadius, mock, tests)
}

func TestPointsHandlerHandleGetPointDetails(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	handler := &PointsHandler{}

	// validPoint, err := go_geom.NewPoint(go_geom.XY).SetSRID(4326).SetCoords([]float64{-6.268726783812387, 53.3484472329815})
	// if err != nil {
	// 	t.Fatal(err)
	// }

	tests := []testutil.HandlerTestDefinition{
		{
			Name:   "Valid input",
			Method: "GET",
			Route:  "/points/inRadius",
			PathParams: map[string]string{
				"id": "236",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT Details::jsonb FROM points WHERE id = \$1 LIMIT 1`).
					WithArgs(int64(236)).
					WillReturnRows(pgxmock.NewRows([]string{"details"}).
						AddRow([]byte(`{"example": "test"}`)))
			},
			ExpectedStatus: http.StatusOK,
			ExpectedJSON: `{
				"error": null,
				"response": {
					"content": {
						"example": "test"
						}
					}
				}`,
		},
	}

	testutil.RunTests(t, handler.HandleGetPointDetails, mock, tests)
}
