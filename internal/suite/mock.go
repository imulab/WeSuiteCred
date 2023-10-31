package suite

import (
	"absurdlab.io/WeSuiteCred/internal/x"
	"github.com/jarcoal/httpmock"
	"net/http"
)

func MockGetSuiteAccessTokenEndpoint() {
	httpmock.RegisterResponder(
		http.MethodPost,
		"=~^"+getSuiteAccessTokenUrl+".*",
		func(r *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponderOrPanic(http.StatusOK, map[string]interface{}{
				"errcode":            0,
				"errmsg":             "ok",
				"suite_access_token": "61W3mEpU66027wgNZ_MhGHNQDHnFATkDa9-2llMBjUwxRSNPbVsMmyD-yq8wZETSoE5NQgecigDrSHkPtIYA" + x.RandAlphaNumeric(4),
				"expires_in":         7200,
			})(r)
		},
	)
}
