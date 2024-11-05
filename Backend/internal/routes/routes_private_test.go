//go:build private

package routes

import (
	"net/http"
	"testing"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/testutil"
)


func TestPrivateRoutes(t *testing.T) {
	tests := []testutil.RouterTestDefinition{
		{
			Name:               "Test POST /points/",
			Path:               "/points/",
			Method:             "POST",
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Test PUT /points/{id}",
			Path:               "/points/123",
			Method:             "PUT",
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Test DELETE /points/{id}",
			Path:               "/points/123",
			Method:             "DELETE",
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Test non-existing private route",
			Path:               "/points/nonexistent",
			Method:             "GET",
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	testutil.RunRouterTests(t, tests, private())
}
