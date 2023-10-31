package listener

import (
	"absurdlab.io/WeSuiteCred/cmd/internal"
	"absurdlab.io/WeSuiteCred/internal/corp"
	"absurdlab.io/WeSuiteCred/internal/sqlitedb"
	"absurdlab.io/WeSuiteCred/internal/suite"
	"absurdlab.io/WeSuiteCred/internal/wt"
	"absurdlab.io/WeSuiteCred/internal/x"
	"context"
	"errors"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"net/url"
	"os"
	"time"
)

func Command() *cli.Command {
	conf := new(config)

	return &cli.Command{
		Name:   "listener",
		Usage:  "Start listening for credential related messages from WeTriage",
		Flags:  conf.flags(),
		Action: func(c *cli.Context) error { return runApp(c.Context, conf) },
	}
}

func runApp(ctx context.Context, conf *config) error {
	return fx.New(
		fx.NopLogger,
		fx.Supply(conf),
		fx.Provide(
			newLogger,
			newMqttClient,
			newSuiteProperties,
			sqlitedb.New,
			suite.NewAccessTokenSupplier,
			suite.NewSuiteTicketInfoSubscriber,
			corp.NewService,
			corp.NewCreateAuthInfoSubscriber,
			corp.NewChangeAuthInfoSubscriber,
			corp.NewCancelAuthInfoSubscriber,
			corp.NewResetPermanentCodeInfoSubscriber,
		),
		fx.Invoke(
			func(db *bun.DB) {
				db.RegisterModel(
					(*suite.Ticket)(nil),
					(*corp.AuthInfo)(nil),
				)
			},
			sqlitedb.Migrate,
			func(logger *zerolog.Logger, cm *autopaho.ConnectionManager) {
				logger.Info().Msg("WeSuiteCred waiting for messages.")
				<-cm.Done()
			},
		),
	).Start(ctx)
}

func newLogger(c *config) *zerolog.Logger {
	var lvl = zerolog.InfoLevel
	if c.Debug {
		lvl = zerolog.DebugLevel
	}

	logger := zerolog.New(os.Stderr).
		Level(lvl).
		With().
		Timestamp().
		Logger()

	return &logger
}

func newMqttClient(c *config, logger *zerolog.Logger, subscribers []wt.Subscriber) (*autopaho.ConnectionManager, error) {
	brokerUrl, err := url.Parse(c.MqttUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid mqtt broker url: %w", err)
	}

	pahoConnUp := make(chan struct{})

	ctx := context.Background()

	var (
		errorLogger             = internal.NewPahoZeroLogger(logger)
		debugLogger paho.Logger = paho.NOOPLogger{}
	)
	if c.Debug {
		debugLogger = errorLogger
	}

	router := paho.NewStandardRouter()
	for _, s := range subscribers {
		router.RegisterHandler(s.Option().Topic, func(p *paho.Publish) {
			if handleErr := s.Handle(p); handleErr != nil {
				logger.Err(handleErr).Str("topic", p.Topic).Msg("failed to handle message")
				return
			}

			logger.Info().Str("topic", p.Topic).Msg("handled message")
		})
	}

	cm, err := autopaho.NewConnection(ctx, autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{brokerUrl},
		KeepAlive:         60,
		ConnectRetryDelay: time.Millisecond,
		ConnectTimeout:    15 * time.Second,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, _ *paho.Connack) {
			if _, subErr := cm.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: lo.Map(subscribers, func(s wt.Subscriber, _ int) paho.SubscribeOptions {
					return s.Option()
				}),
			}); subErr != nil {
				panic(subErr)
			}
			pahoConnUp <- struct{}{}
		},
		Debug:      debugLogger,
		PahoDebug:  debugLogger,
		PahoErrors: errorLogger,
		ClientConfig: paho.ClientConfig{
			ClientID: fmt.Sprintf("WeSuiteCred@%s", x.RandAlphaNumeric(6)),
			Router:   router,
		},
	})
	if err != nil {
		return nil, err
	}

	select {
	case <-pahoConnUp:
	case <-time.After(1 * time.Minute):
		return nil, errors.New("timeout exceeded when connecting to mqtt broker")
	}

	return cm, nil
}

func newSuiteProperties(c *config) *suite.Properties {
	return &suite.Properties{
		Id:                c.SuiteId,
		Secret:            c.SuiteSecret,
		AccessTokenLeeway: 30 * time.Second,
	}
}
