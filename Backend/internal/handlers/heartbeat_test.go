package handlers

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/testutil"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
)

func TestHeartbeat(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	tests := []testutil.RawRequestTestDefinition{
		{
			Name:   "Valid",
			Method: "GET",
			Route:  "/heartbeat",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("-- ping").WillReturnResult(pgconn.NewCommandTag("EXEC")).WillDelayFor(time.Duration(1 * time.Second))
			},
			ExpectedStatus:   http.StatusOK,
			ExpectedResponse: "true",
		},
		{
			Name:   "Database timeout",
			Method: "GET",
			Route:  "/heartbeat",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("-- ping").WillReturnResult(pgconn.NewCommandTag("EXEC")).WillDelayFor(time.Duration(6 * time.Second))
			},
			ExpectedStatus:   http.StatusInternalServerError,
			ExpectedResponse: "false",
		},
	}

	testutil.RunRawRequestTests(t, HandleHeartbeat, mock, tests)
}
