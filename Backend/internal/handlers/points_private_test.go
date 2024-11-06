//go:build private

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	db "github.com/2024-CMPU9010-GROUP-3/magpie/internal/db/private"
	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	testutil "github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/testutil"
	"github.com/pashagolub/pgxmock/v4"
	go_geom "github.com/twpayne/go-geom"
)

func TestPointsHandlerHandlePost(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	handler := &PointsHandler{}

	tests := []testutil.HandlerTestDefinition{
		{
			Name:   "Valid input",
			Method: "POST",
			Route:  "/points",
			InputJSON: `{
				"longlat": {
					"type": "Point",
					"coordinates": [11, 12]
				},
				"type": "parking",
				"details": {
					"test": 1234
				}
			}`,
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("INSERT INTO points").
					WithArgs(
						pgxmock.AnyArg(),
						db.PointType("parking"),
						[]byte(`{"test":1234}`),
					).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(1)))
			},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:   "Invalid input",
			Method: "POST",
			Route:  "/points",
			InputJSON: `{
				"type1": "parking",
				"details": {
					"test": 1234
				}
			}`,
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  customErrors.Payload.InvalidPayloadPointError.ErrorMsg,
		},
		{
			Name:   "Invalid geometry",
			Method: "POST",
			Route:  "/points",
			InputJSON: `{
				"longlat": {
					"type": "InvalidType",
					"coordinates": [11, 12]
				},
				"type": "parking",
				"details": {
					"test": 1234
				}
			}`,
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  customErrors.Payload.InvalidPayloadPointError.ErrorMsg,
		},
		{
			Name:   "Invalid type",
			Method: "POST",
			Route:  "/points",
			InputJSON: `{
				"longlat": {
					"type": "Point",
					"coordinates": [11, 12]
				},
				"type": "invalid",
				"details": {
					"test": 1234
				}
			}`,
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  customErrors.Payload.InvalidPayloadPointError.ErrorMsg,
		},
		{
			Name:   "Valid geometry, but not point",
			Method: "POST",
			Route:  "/points",
			InputJSON: `{
				"longlat": {
					"type": "Polygon",
					"coordinates": [
						[
							[100.0, 0.0],
							[101.0, 0.0],
							[101.0, 1.0],
							[100.0, 1.0],
							[100.0, 0.0]
						]
					]
				},
				"type": "parking",
				"details": {
					"test": 1234
				}
			}`,
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  customErrors.Payload.InvalidPayloadPointError.ErrorMsg,
		},
		{
			Name:   "Database error on insert",
			Method: "POST",
			Route:  "/points",
			InputJSON: `{
				"longlat": {
					"type": "Point",
					"coordinates": [11, 12]
				},
				"type": "parking",
				"details": {
					"test": 1234
				}
			}`,
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("INSERT INTO points").
					WithArgs(
						pgxmock.AnyArg(),
						db.PointType("parking"),
						[]byte(`{"test":1234}`),
					).
					// Simulate a database error
					WillReturnError(fmt.Errorf("Unable to connect to database"))
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  customErrors.Database.UnknownDatabaseError.ErrorMsg,
		},
	}

	testutil.RunHandlerTests(t, handler.HandlePost, mock, tests)
}

func TestPointsHandlerHandlePut(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	handler := &PointsHandler{}

	tests := []testutil.HandlerTestDefinition{
		{
			Name:   "Valid input",
			Method: "PUT",
			Route:  "/points",
			InputJSON: `{
				"longlat": {
					"type": "Point",
					"coordinates": [11, 12]
				},
				"type": "parking",
				"details": {
					"test": 1234
				}
			}`,
			PathParams: map[string]string{
				"id": "123456",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				point := go_geom.NewPoint(go_geom.XY)
				_, _ = point.SetCoords(go_geom.Coord{11, 12})
				mock.ExpectExec("UPDATE points").
					WithArgs(
						int64(123456),
						point,
						db.PointType("parking"),
						[]byte(`{"test":1234}`),
					).WillReturnResult(pgxmock.NewResult("UPDATED", 1))
			},
			ExpectedStatus: http.StatusAccepted,
		},
		{
			Name:   "Invlid id",
			Method: "PUT",
			Route:  "/points",
			InputJSON: `{
				"longlat": {
					"type": "Point",
					"coordinates": [11, 12]
				},
				"type": "parking",
				"details": {
					"test": 1234
				}
			}`,
			PathParams: map[string]string{
				"id": "abdcd",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  customErrors.Parameter.InvalidIntError.ErrorMsg,
		},
		{
			Name:   "Invalid input",
			Method: "PUT",
			Route:  "/points",
			InputJSON: `{
				"type1": "parking",
				"details": {
					"test": 1234
				}
			}`,
			PathParams: map[string]string{
				"id": "123456",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  customErrors.Payload.InvalidPayloadPointError.ErrorMsg,
		},
		{
			Name:   "Parameter Missing",
			Method: "PUT",
			Route:  "/points",
			InputJSON: `{
				"longlat": {
					"type": "Point",
					"coordinates": [11, 12]
				},
				"details": {
					"test": 1234
				}
			}`,
			PathParams: map[string]string{
				"id": "123456",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  customErrors.Parameter.RequiredParameterMissingError.ErrorMsg,
		},
		{
			Name:   "Invalid geometry",
			Method: "PUT",
			Route:  "/points",
			InputJSON: `{
				"longlat": {
					"type": "InvalidType",
					"coordinates": [11, 12]
				},
				"type": "parking",
				"details": {
					"test": 1234
				}
			}`,
			PathParams: map[string]string{
				"id": "123456",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  customErrors.Payload.InvalidPayloadPointError.ErrorMsg,
		},
		{
			Name:   "Valid geometry, but not point",
			Method: "POST",
			Route:  "/points",
			InputJSON: `{
				"longlat": {
					"type": "Polygon",
					"coordinates": [
						[
							[100.0, 0.0],
							[101.0, 0.0],
							[101.0, 1.0],
							[100.0, 1.0],
							[100.0, 0.0]
						]
					]
				},
				"type": "parking",
				"details": {
					"test": 1234
				}
			}`,
			PathParams: map[string]string{
				"id": "123456",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  customErrors.Payload.InvalidPayloadPointError.ErrorMsg,
		},
		{
			Name:   "Database error on insert",
			Method: "PUT",
			Route:  "/points",
			InputJSON: `{
				"longlat": {
					"type": "Point",
					"coordinates": [11, 12]
				},
				"type": "parking",
				"details": {
					"test": 1234
				}
			}`,
			PathParams: map[string]string{
				"id": "123456",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				point := go_geom.NewPoint(go_geom.XY)
				_, _ = point.SetCoords(go_geom.Coord{11, 12})
				mock.ExpectExec("UPDATE points").
					WithArgs(
						int64(123456),
						point,
						db.PointType("parking"),
						[]byte(`{"test":1234}`),
					).WillReturnResult(pgxmock.NewResult("UPDATED", 1)).
					// Simulate a database error
					WillReturnError(fmt.Errorf("Unable to connect to database"))
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  customErrors.Database.UnknownDatabaseError.ErrorMsg,
		},
	}

	testutil.RunHandlerTests(t, handler.HandlePut, mock, tests)
}

func TestPointsHandlerHandleDelete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	handler := &PointsHandler{}

	tests := []testutil.HandlerTestDefinition{
		{
			Name:   "Valid input",
			Method: "DELETE",
			Route:  "/points",
			PathParams: map[string]string{
				"id": "123456",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				point := go_geom.NewPoint(go_geom.XY)
				_, _ = point.SetCoords(go_geom.Coord{11, 12})
				mock.ExpectExec("DELETE FROM points").
					WithArgs(
						int64(123456),
					).WillReturnResult(pgxmock.NewResult("DELETED", 1))
			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON: `{
				"error": null,
				"response": {
					"content": {
						"id": 123456
					}
				}
			}`,
		},
		{
			Name:   "Invalid id",
			Method: "DELETE",
			Route:  "/points",
			PathParams: map[string]string{
				"id": "abdcd",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  customErrors.Parameter.InvalidIntError.ErrorMsg,
			ExpectedJSON: `{
				"error": {
					"errorCode": 1204,
					"errorMsg": "Parameter invalid, expected type Int"
				},
				"response": null
			}`,
		},
		{
			Name: "Database error on delete",
			PathParams: map[string]string{
				"id": "123456",
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				point := go_geom.NewPoint(go_geom.XY)
				_, _ = point.SetCoords(go_geom.Coord{11, 12})
				mock.ExpectExec("DELETE FROM points").
					WithArgs(
						int64(123456),
					).
					// Simulate a database error
					WillReturnError(fmt.Errorf("Unable to connect to database"))
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  customErrors.Database.UnknownDatabaseError.ErrorMsg,
			ExpectedJSON: `{
				"error": {
					"errorCode": 1104,
					"errorMsg": "Unknown database error",
					"cause": "Unable to connect to database"
				},
				"response": null
			}`,
		},
	}

	testutil.RunHandlerTests(t, handler.HandleDelete, mock, tests)
}
