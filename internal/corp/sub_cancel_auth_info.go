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

func NewCancelAuthInfoSubscriber(suiteProps *suite.Properties, service *Service) wt.Subscriber {
	return &cancelAuthInfoSubscriber{
		suiteProps: suiteProps,
		service:    service,
	}
}

type cancelAuthInfoSubscriber struct {
	suiteProps *suite.Properties
	service    *Service
}

func (s *cancelAuthInfoSubscriber) Option() paho.SubscribeOptions {
	return paho.SubscribeOptions{Topic: "T/WeTriage/cancel_auth_info"}
}

func (s *cancelAuthInfoSubscriber) Handle(pub *paho.Publish) error {
	var body wt.Payload[cancelAuthInfo]
	if err := json.Unmarshal(pub.Payload, &body); err != nil {
		return err
	}

	switch {
	case body.Content.SuiteId != s.suiteProps.Id:
		return errors.New("suite_id mismatch")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.service.OnAuthorizationRemoved(ctx, body.Content.AuthCorpId)
}

type cancelAuthInfo struct {
	SuiteId    string `json:"suite_id"`
	InfoType   string `json:"info_type"`
	Timestamp  int64  `json:"timestamp"`
	AuthCorpId string `json:"auth_corp_id"`
}
