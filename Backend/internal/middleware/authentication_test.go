package middleware

import (
	"net/http"
	"testing"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/testutil"
)

const (
	userIdString    = `abcdeeb9-31cd-4c62-b4fe-c67ef8621cd4`
	userIdStringAlt = `8c95319a-5192-484d-bbbb-81cd65793788`
)

func TestAccessOwnerOnly(t *testing.T) {

	tests := []testutil.MiddlewareTestDefinition{
		{
			Name:               "Invalid UUID in path",
			IdPathParam:        "invalid-uuid",
			TokenUserId:        "some-token-user-id",
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody: `{
					"error": {
						"errorCode": 1202,
						"errorMsg": "Parameter invalid, expected type UUIDv4"
					},
					"response": null
				}`,
		},
		{
			Name:               "Missing token_user_id in context",
			IdPathParam:        userIdString,
			TokenUserId:        nil,
			ExpectedStatusCode: http.StatusUnauthorized,
			ExpectedBody: `{
					"error": {
						"errorCode": 1404,
						"errorMsg": "UserId from token missing in context"
					},
					"response": null
				}`,
		},
		{
			Name:               "Unauthorized access - user ID mismatch",
			IdPathParam:        userIdString,
			TokenUserId:        userIdStringAlt,
			ExpectedStatusCode: http.StatusUnauthorized,
			ExpectedBody: `{
					"error": {
						"errorCode": 1401,
						"errorMsg": "Unauthorized"
					},
					"response": null
				}`,
		},
		{
			Name:                 "Authorized access",
			IdPathParam:          userIdString,
			TokenUserId:          userIdString,
			ExpectedStatusCode:   http.StatusOK,
			ExpectedBodyContains: "next handler called",
		},
	}

	testutil.RunMiddlewareTests(t, accessOwnerOnly, tests)
}

func TestPublic(t *testing.T) {
	tests := []testutil.MiddlewareTestDefinition{
		{
			Name:                 "Any access",
			ExpectedStatusCode:   http.StatusOK,
			ExpectedBodyContains: "next handler called",
		},
	}
	testutil.RunMiddlewareTests(t, accessPublic, tests)
}
