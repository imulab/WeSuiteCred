package wt

import (
	"encoding/json"
	"errors"
	"github.com/eclipse/paho.golang/paho"
)

func NewResetPermanentCodeInfoSubscriber(props *Properties, corpService *CorpService) Subscriber {
	return &resetPermanentCodeInfoSubscriber{
		props:       props,
		corpService: corpService,
	}
}

type resetPermanentCodeInfoSubscriber struct {
	props       *Properties
	corpService *CorpService
}

func (s *resetPermanentCodeInfoSubscriber) Option() paho.SubscribeOptions {
	return paho.SubscribeOptions{Topic: "T/WeTriage/reset_permanent_code_info"}
}

func (s *resetPermanentCodeInfoSubscriber) Handle(pub *paho.Publish) error {
	var body payload[resetPermanentCodeInfo]
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

type resetPermanentCodeInfo struct {
	SuiteId   string `json:"suite_id"`
	AuthCode  string `json:"auth_code"`
	InfoType  string `json:"info_type"`
	Timestamp int64  `json:"timestamp"`
}
