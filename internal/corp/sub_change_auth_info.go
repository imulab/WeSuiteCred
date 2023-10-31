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

func NewChangeAuthInfoSubscriber(suiteProps *suite.Properties, service *Service) wt.Subscriber {
	return &changeAuthInfoSubscriber{
		suiteProps: suiteProps,
		service:    service,
	}
}

type changeAuthInfoSubscriber struct {
	suiteProps *suite.Properties
	service    *Service
}

func (s *changeAuthInfoSubscriber) Option() paho.SubscribeOptions {
	return paho.SubscribeOptions{Topic: "T/WeTriage/change_auth_info"}
}

func (s *changeAuthInfoSubscriber) Handle(pub *paho.Publish) error {
	var body wt.Payload[ChangeAuthInfo]
	if err := json.Unmarshal(pub.Payload, &body); err != nil {
		return err
	}

	switch {
	case body.Content.SuiteId != s.suiteProps.Id:
		return errors.New("suite_id mismatch")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.service.OnAuthorizationChanged(ctx, body.Content.AuthCorpId)
}

type ChangeAuthInfo struct {
	SuiteId    string `json:"suite_id"`
	InfoType   string `json:"info_type"`
	Timestamp  int64  `json:"timestamp"`
	AuthCorpId string `json:"auth_corp_id"`
	State      string `json:"state"`
}
