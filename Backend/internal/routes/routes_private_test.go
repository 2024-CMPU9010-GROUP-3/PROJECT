//go:build private

package routes

import (
	"context"
	"net/http"
	"testing"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/handlers"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/testutil"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
)

func TestPrivateRoutes(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	handlers.RegisterDatabaseConnection(&ctx, mock)

	tests := []testutil.RouterTestDefinition{
		{
			Name:   "Test POST /points/",
			Path:   "/points/",
			Method: "POST",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no db calls expected
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:   "Test PUT /points/{id}",
			Path:   "/points/123",
			Method: "PUT",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no db calls expected
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:   "Test DELETE /points/{id}",
			Path:   "/points/123",
			Method: "DELETE",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`DELETE FROM points WHERE id = \$1`).
					WithArgs(int64(123)).
					WillReturnResult(pgconn.NewCommandTag("DELETED"))
			},
			ExpectedStatusCode: http.StatusAccepted,
		},
		{
			Name:   "Test non-existing private route",
			Path:   "/points/nonexistent",
			Method: "GET",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no db calls expected
			},
			ExpectedStatusCode: http.StatusMethodNotAllowed,
		},
	}

	testutil.RunRouterTests(t, tests, private(), mock)
}
