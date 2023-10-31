package utils

import (
	"absurdlab.io/WeSuiteCred/internal/corp"
	"absurdlab.io/WeSuiteCred/internal/wt"
	"absurdlab.io/WeSuiteCred/internal/x"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"net/url"
	"time"
)

func simulateChangeAuth() *cli.Command {
	var (
		mqttUrl string
		suiteId string
		corpId  string
	)

	newMqttClient := func(ctx context.Context) (*autopaho.ConnectionManager, error) {
		brokerUrl, err := url.Parse(mqttUrl)
		if err != nil {
			return nil, fmt.Errorf("invalid mqtt broker url: %w", err)
		}

		pahoConnUp := make(chan struct{})

		cm, err := autopaho.NewConnection(ctx, autopaho.ClientConfig{
			BrokerUrls:        []*url.URL{brokerUrl},
			KeepAlive:         60,
			ConnectRetryDelay: time.Millisecond,
			ConnectTimeout:    15 * time.Second,
			OnConnectionUp:    func(cm *autopaho.ConnectionManager, _ *paho.Connack) { pahoConnUp <- struct{}{} },
			Debug:             paho.NOOPLogger{},
			PahoDebug:         paho.NOOPLogger{},
			PahoErrors:        paho.NOOPLogger{},
			ClientConfig: paho.ClientConfig{
				ClientID: fmt.Sprintf("WeSuiteCred_SimChangeAuth@%s", x.RandAlphaNumeric(6)),
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
	return &cli.Command{
		Name:  "simulate-change-auth",
		Usage: "Refresh app permission for a corporation by simulating a change_auth_info event",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "mqtt-url",
				Usage:       "MQTT broker url",
				Destination: &mqttUrl,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "suite-id",
				Usage:       "Suite id",
				Destination: &suiteId,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "corp-id",
				Usage:       "Corporation id",
				Destination: &corpId,
				Required:    true,
			},
		},
		Action: func(cc *cli.Context) error {
			cm, err := newMqttClient(cc.Context)
			if err != nil {
				return err
			}

			topic := "T/WeTriage/change_auth_info"

			payload := wt.Payload[corp.ChangeAuthInfo]{
				Id:        uuid.New().String(),
				CreatedAt: time.Now().Unix(),
				Topic:     topic,
				Content: corp.ChangeAuthInfo{
					SuiteId:    suiteId,
					InfoType:   "change_auth",
					Timestamp:  time.Now().Unix(),
					AuthCorpId: corpId,
					State:      "simulated",
				},
			}

			payloadBytes, _ := json.Marshal(payload)

			if _, err = cm.Publish(cc.Context, &paho.Publish{
				QoS:     2,
				Retain:  false,
				Topic:   topic,
				Payload: payloadBytes,
			}); err != nil {
				return err
			}

			fmt.Println("Simulated message sent.")

			return nil
		},
	}
}
