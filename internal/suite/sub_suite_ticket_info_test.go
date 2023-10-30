package suite

import (
	"absurdlab.io/WeSuiteCred/internal/sqlitedb"
	"absurdlab.io/WeSuiteCred/internal/wt"
	"context"
	"github.com/eclipse/paho.golang/paho"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
	"testing"
	"time"
)

func TestSuiteTicketInfoSubscriber(t *testing.T) {
	logger := zerolog.New(zerolog.NewTestWriter(t))

	err := fx.New(
		fx.NopLogger,
		fx.Supply(
			&logger,
			&Properties{
				Id:                "wwddddccc7775555aaa",
				Secret:            "ldAE_H9anCRN21GKXVfdAAAAAAAAAAAAAAAAAA",
				AccessTokenLeeway: 30 * time.Second,
			},
		),
		fx.Provide(
			sqlitedb.NewMemory,
			NewSuiteTicketInfoSubscriber,
		),
		fx.Invoke(
			func(db *bun.DB) { db.RegisterModel((*Ticket)(nil)) },
			sqlitedb.Migrate,
			func(sub wt.Subscriber, db *bun.DB) {
				err := sub.Handle(&paho.Publish{
					QoS:   2,
					Topic: "T/WeTriage/suite_ticket_info",
					Payload: []byte(`
{
	"id": "test",
	"create_at": 0,
	"topic": "suite_ticket_info",
	"content": {
		"suite_id": "wwddddccc7775555aaa",
		"info_type": "suite_ticket",
		"timestamp": 0,
		"suite_ticket": "asdfasdfasdfa"
	}
}
`),
				})

				if !assert.NoError(t, err) {
					return
				}

				var ticket Ticket
				if !assert.NoError(t, db.NewSelect().Model(&ticket).Where("id = ?", 1).Scan(context.Background())) {
					return
				}

				assert.Equal(t, "asdfasdfasdfa", ticket.Ticket)
			},
		),
	).Start(context.TODO())

	require.NoError(t, err)
}
