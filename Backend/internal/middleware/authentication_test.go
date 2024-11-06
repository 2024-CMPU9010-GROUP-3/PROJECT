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

func TestAccessAuthenticated(t *testing.T) {

	tests := []testutil.MiddlewareTestDefinition{
		{
			Name:               "Missing authentication cookie",
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
			Name:               "Missing JWT secret",
			AuthCookieValue:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJiM2EyMmYxOS01ZjljLTRkY2MtOWU4Zi00NmU4NDUzNmRiMTAiLCJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwiaWF0IjoxNTE2MjM5MDIyfQ.BDNt2077waOG_3l9Nd-pInTs30AAZ30z18KXB495HII",
			EnvJwtSecret:       "",
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedBody: `{
					"error": {
						"errorCode": 1011,
						"errorMsg": "Could not generate JWT, secret not set"
					},
					"response": null
				}`,
		},
		{
			Name:               "Invalid JWT token",
			AuthCookieValue:    "invalid_jwt_token",
			EnvJwtSecret:       "test_secret",
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedBody: `{
					"error": {
						"errorCode": 1013,
						"errorMsg": "Could not parse JWT"
					},
					"response": null
				}`,
		},
		{
			Name:               "Unexpected JWT signing method",
			AuthCookieValue:    "eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.AbVUinMiT3J_03je8WTOIl-VdggzvoFgnOsdouAs-DLOtQzau9valrq-S6pETyi9Q18HH-EuwX49Q7m3KC0GuNBJAc9Tksulgsdq8GqwIqZqDKmG7hNmDzaQG1Dpdezn2qzv-otf3ZZe-qNOXUMRImGekfQFIuH_MjD2e8RZyww6lbZk",
			EnvJwtSecret:       "test_secret",
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedBody: `{
					"error": {
						"errorCode": 1013,
						"errorMsg": "Could not parse JWT"
					},
					"response": null
				}`,
		},
		{
			Name:               "Subject missing",
			AuthCookieValue:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwiaWF0IjoxNTE2MjM5MDIyfQ.Ogb2idCPd36PAimIsia-7hLdmbVxqXFsP74YLQm9KqI",
			EnvJwtSecret:       "test_secret",
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedBody: `{
					"error": {
						"errorCode": 1013,
						"errorMsg": "Could not parse JWT"
					},
					"response": null
				}`,
		},
		{
			Name:                 "Authorized access with valid JWT",
			AuthCookieValue:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJiM2EyMmYxOS01ZjljLTRkY2MtOWU4Zi00NmU4NDUzNmRiMTAiLCJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwiaWF0IjoxNTE2MjM5MDIyfQ.BDNt2077waOG_3l9Nd-pInTs30AAZ30z18KXB495HII",
			EnvJwtSecret:         "test_secret",
			ExpectedStatusCode:   http.StatusOK,
			ExpectedBodyContains: "next handler called",
		},
	}

	testutil.RunMiddlewareTests(t, accessAuthenticated, tests)
}
