package listener

import (
	"absurdlab.io/WeSuiteCred/cmd/internal"
	"absurdlab.io/WeSuiteCred/internal/stringx"
	"context"
	"errors"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
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
	return nil
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

func newMqttClient(c *config, logger *zerolog.Logger) (*autopaho.ConnectionManager, error) {
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

	cm, err := autopaho.NewConnection(ctx, autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{brokerUrl},
		KeepAlive:         60,
		ConnectRetryDelay: time.Millisecond,
		ConnectTimeout:    15 * time.Second,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, _ *paho.Connack) {
			if _, subErr := cm.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{
					{Topic: "T/WeTriage/suite_ticket_info"},
					{Topic: "T/WeTriage/reset_permanent_code_info"},
				},
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
			//Router: paho.NewSingleHandlerRouter(func(pub *paho.Publish) {
			//	logger.Info().
			//		Str("topic", pub.Topic).
			//		Int("qos", int(pub.QoS)).
			//		RawJSON("payload", pub.Payload).
			//		Msg("Received message")
			//}),
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
