package wt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eclipse/paho.golang/paho"
	"github.com/peterbourgon/diskv/v3"
	"sync"
)

const (
	keySuiteTicket = "suite_ticket"
)

func NewSuiteTicketInfoSubscriber(props *Properties, store *diskv.Diskv) Subscriber {
	return &suiteTicketInfoSubscriber{suiteProps: props, store: store}
}

type suiteTicketInfoSubscriber struct {
	sync.Mutex

	suiteProps *Properties
	store      *diskv.Diskv
}

func (s *suiteTicketInfoSubscriber) Option() paho.SubscribeOptions {
	return paho.SubscribeOptions{
		Topic: "T/WeTriage/suite_ticket_info",
	}
}

func (s *suiteTicketInfoSubscriber) Handle(pub *paho.Publish) error {
	var body payload[suiteTicketInfo]
	if err := json.Unmarshal(pub.Payload, &body); err != nil {
		return err
	}

	switch {
	case body.Content.SuiteId != s.suiteProps.SuiteId:
		return errors.New("suite_id mismatch")
	case len(body.Content.SuiteTicket) == 0:
		return errors.New("suite_ticket is empty")
	}

	s.Lock()
	defer s.Unlock()

	if s.store.Has(keySuiteTicket) {
		if err := s.store.Erase(keySuiteTicket); err != nil {
			return fmt.Errorf("failed to erase suite_ticket: %w", err)
		}
	}

	if err := s.store.WriteString(keySuiteTicket, body.Content.SuiteTicket); err != nil {
		return fmt.Errorf("failed to write suite_ticket: %w", err)
	}

	return nil
}

type suiteTicketInfo struct {
	SuiteId     string `json:"suite_id"`
	InfoType    string `json:"info_type"`
	Timestamp   int64  `json:"timestamp"`
	SuiteTicket string `json:"suite_ticket"`
}
