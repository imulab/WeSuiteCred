package listener

import (
	"absurdlab.io/WeSuiteCred/cmd/internal"
	"absurdlab.io/WeSuiteCred/internal/stringx"
	"absurdlab.io/WeSuiteCred/internal/wt"
	"context"
	"errors"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/peterbourgon/diskv/v3"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"net/url"
	"os"
	"strings"
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
			newDiskStore,
			newMqttClient,
			newWtProperties,
			wt.NewCorpSecretDao,
			wt.NewCorpAuthInfoDao,
			wt.NewSuiteAccessTokenSupplier,
			wt.NewCorpService,
			wt.NewSuiteTicketInfoSubscriber,
			wt.NewCreateAuthInfoSubscriber,
			wt.NewResetPermanentCodeInfoSubscriber,
		),
		fx.Invoke(
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
			}
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
			ClientID: fmt.Sprintf("WeSuiteCred@%s", stringx.RandAlphaNumeric(6)),
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

func newDiskStore(c *config) *diskv.Diskv {
	return diskv.New(diskv.Options{
		BasePath: c.StoreDir,
		Transform: func(s string) []string {
			s = strings.TrimSpace(s)
			return strings.Split(s, "/")
		},
		CacheSizeMax: 1024 * 1024, // 1MB
	})
}

func newWtProperties(c *config) *wt.Properties {
	return &wt.Properties{
		SuiteId:                c.SuiteId,
		SuiteSecret:            c.SuiteSecret,
		SuiteAccessTokenLeeway: 30 * time.Second,
	}
}
