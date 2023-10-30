package suite

import (
	"absurdlab.io/WeSuiteCred/internal/sqlitedb"
	"absurdlab.io/WeSuiteCred/internal/x"
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
	"net/http"
	"testing"
	"time"
)

func TestAccessTokenSupplier(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

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

	logger := zerolog.New(zerolog.NewTestWriter(t))

	err := fx.New(
		fx.NopLogger,
		fx.Supply(
			&logger,
			&Properties{
				Id:                "wwddddccc7775555aaa",
				Secret:            "ldAE_H9anCRN21GKXVfdAAAAAAAAAAAAAAAAAA",
				AccessTokenLeeway: 30 * time.Second,
			}),
		fx.Provide(
			sqlitedb.NewMemory,
			NewAccessTokenSupplier,
		),
		fx.Invoke(
			func(db *bun.DB) { db.RegisterModel((*Ticket)(nil)) },
			sqlitedb.Migrate,
			loadTestTicket,
			func(supplier *AccessTokenSupplier) {
				accessToken, err := supplier.Get()
				if assert.NoError(t, err) {
					assert.NotEmpty(t, accessToken)
				}

				accessTokenTake2, err := supplier.Get()
				if assert.NoError(t, err) {
					assert.Equal(t, accessToken, accessTokenTake2)
					assert.Equal(t, 1, httpmock.GetTotalCallCount())
				}

				supplier.Reset()

				accessTokenTake3, err := supplier.Get()
				if assert.NoError(t, err) {
					assert.NotEmpty(t, accessTokenTake3)
					assert.NotEqual(t, accessToken, accessTokenTake3)
					assert.Equal(t, 2, httpmock.GetTotalCallCount())
				}
			},
		),
	).Start(context.TODO())

	require.NoError(t, err)
}

func loadTestTicket(db *bun.DB) error {
	ticket := Ticket{
		ID:     1,
		Ticket: "Cfp0_givEagXcYJIztF6sfbdmIZCmpaR8ZBsvJEFFNBrWmnD5-CGYJ3_NhYexMyw",
	}

	_, err := db.NewInsert().Model(&ticket).Exec(context.TODO())

	return err
}
