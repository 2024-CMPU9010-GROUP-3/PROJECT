//go:build public

package routes

import (
	"net/http"
	"testing"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/testutil"
)

func TestPublicRoutes(t *testing.T) {
	tests := []testutil.RouterTestDefinition{
		{Name: "Test /points/inRadius with GET", Path: "/points/inRadius?long=0.0&lat=0.0&radius=1000", Method: "GET", ExpectedStatusCode: http.StatusUnauthorized},
		{Name: "Test /points/{id} with GET", Path: "/points/12345", Method: "GET", ExpectedStatusCode: http.StatusUnauthorized},
		{Name: "Test /auth/User/login with POST", Path: "/auth/User/login", Method: "POST", ExpectedStatusCode: http.StatusBadRequest},
		{Name: "Test /auth/User/ with POST", Path: "/auth/User/", Method: "POST", ExpectedStatusCode: http.StatusBadRequest},
		{Name: "Test /auth/User/{id} with GET", Path: "/auth/User/63275AAE-C901-4537-9517-1C9B6F19264A", Method: "GET", ExpectedStatusCode: http.StatusUnauthorized},
		{Name: "Test /auth/User/{id} with PUT", Path: "/auth/User/63275AAE-C901-4537-9517-1C9B6F19264A", Method: "PUT", ExpectedStatusCode: http.StatusUnauthorized},
		{Name: "Test /auth/User/{id} with DELETE", Path: "/auth/User/63275AAE-C901-4537-9517-1C9B6F19264A", Method: "DELETE", ExpectedStatusCode: http.StatusUnauthorized},
		{Name: "Test non-existing route", Path: "/", Method: "GET", ExpectedStatusCode: http.StatusNotFound},
	}

	testutil.RunRouterTests(t, tests, public())
}
