package corp

import (
	"absurdlab.io/WeSuiteCred/internal/suite"
	"absurdlab.io/WeSuiteCred/internal/wt"
	"context"
	"encoding/json"
	"errors"
	"github.com/eclipse/paho.golang/paho"
	"time"
)

func NewCreateAuthInfoSubscriber(props *suite.Properties, service *Service) wt.Subscriber {
	return &createAuthInfoSubscriber{
		suiteProps: props,
		service:    service,
	}
}

type createAuthInfoSubscriber struct {
	suiteProps *suite.Properties
	service    *Service
}

func (s *createAuthInfoSubscriber) Option() paho.SubscribeOptions {
	return paho.SubscribeOptions{Topic: "T/WeTriage/create_auth_info"}
}

func (s *createAuthInfoSubscriber) Handle(pub *paho.Publish) error {
	var body wt.Payload[createAuthInfo]
	if err := json.Unmarshal(pub.Payload, &body); err != nil {
		return err
	}

	switch {
	case body.Content.SuiteId != s.suiteProps.Id:
		return errors.New("suite_id mismatch")
	case len(body.Content.AuthCode) == 0:
		return errors.New("auth_code is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.service.OnNewAuthCode(ctx, body.Content.AuthCode)
}

type createAuthInfo struct {
	SuiteId   string `json:"suite_id"`
	AuthCode  string `json:"auth_code"`
	InfoType  string `json:"info_type"`
	Timestamp int64  `json:"timestamp"`
	State     string `json:"state"`
}
