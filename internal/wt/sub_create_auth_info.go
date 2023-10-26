package wt

import (
	"encoding/json"
	"errors"
	"github.com/eclipse/paho.golang/paho"
)

func NewCreateAuthInfoSubscriber(props *Properties, corpService *CorpService) Subscriber {
	return &createAuthInfoSubscriber{
		props:       props,
		corpService: corpService,
	}
}

type createAuthInfoSubscriber struct {
	props       *Properties
	corpService *CorpService
}

func (s *createAuthInfoSubscriber) Option() paho.SubscribeOptions {
	return paho.SubscribeOptions{Topic: "T/WeTriage/create_auth_info"}
}

func (s *createAuthInfoSubscriber) Handle(pub *paho.Publish) error {
	var body payload[createAuthInfo]
	if err := json.Unmarshal(pub.Payload, &body); err != nil {
		return err
	}

	switch {
	case body.Content.SuiteId != s.props.SuiteId:
		return errors.New("suite_id mismatch")
	case len(body.Content.AuthCode) == 0:
		return errors.New("auth_code is empty")
	}

	return s.corpService.UpdateSecret(body.Content.AuthCode)
}

type createAuthInfo struct {
	SuiteId   string `json:"suite_id"`
	AuthCode  string `json:"auth_code"`
	InfoType  string `json:"info_type"`
	Timestamp int64  `json:"timestamp"`
	State     string `json:"state"`
}
