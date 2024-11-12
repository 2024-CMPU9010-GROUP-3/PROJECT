//go:build public

package routes

import (
	"context"
	"net/http"
	"testing"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/handlers"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/testutil"
	"github.com/pashagolub/pgxmock/v4"
)

func TestPublicRoutes(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	handlers.RegisterDatabaseConnection(&ctx, mock)

	tests := []testutil.RouterTestDefinition{
		{
			Name:               "Test /points/inRadius with GET",
			Path:               "/points/inRadius?long=0.0&lat=0.0&radius=1000",
			Method:             "GET",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no db calls expected
			},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
		{
			Name:               "Test /points/{id} with GET",
			Path:               "/points/12345",
			Method:             "GET",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no db calls expected
			},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
		{
			Name:               "Test /auth/User/login with POST",
			Path:               "/auth/User/login",
			Method:             "POST",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no db calls expected
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Test /auth/User/ with POST",
			Path:               "/auth/User/",
			Method:             "POST",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no db calls expected
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Test /auth/User/{id} with GET",
			Path:               "/auth/User/63275AAE-C901-4537-9517-1C9B6F19264A",
			Method:             "GET",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no db calls expected
			},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
		{
			Name:               "Test /auth/User/{id} with PUT",
			Path:               "/auth/User/63275AAE-C901-4537-9517-1C9B6F19264A",
			Method:             "PUT",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no db calls expected
			},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
		{
			Name:               "Test /auth/User/{id} with DELETE",
			Path:               "/auth/User/63275AAE-C901-4537-9517-1C9B6F19264A",
			Method:             "DELETE",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no db calls expected
			},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
		{
			Name:               "Test non-existing route",
			Path:               "/",
			Method:             "GET",
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no db calls expected
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	testutil.RunRouterTests(t, tests, public(), mock)
}
