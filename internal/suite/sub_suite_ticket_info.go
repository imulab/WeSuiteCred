package suite

import (
	"absurdlab.io/WeSuiteCred/internal/wt"
	"context"
	"encoding/json"
	"errors"
	"github.com/eclipse/paho.golang/paho"
	"github.com/uptrace/bun"
	"time"
)

func NewSuiteTicketInfoSubscriber(props *Properties, db *bun.DB) wt.Subscriber {
	return &suiteTicketInfoSubscriber{props: props, db: db}
}

type suiteTicketInfoSubscriber struct {
	props *Properties
	db    *bun.DB
}

func (s *suiteTicketInfoSubscriber) Option() paho.SubscribeOptions {
	return paho.SubscribeOptions{
		Topic: "T/WeTriage/suite_ticket_info",
	}
}

func (s *suiteTicketInfoSubscriber) Handle(pub *paho.Publish) error {
	var body wt.Payload[suiteTicketInfo]
	if err := json.Unmarshal(pub.Payload, &body); err != nil {
		return err
	}

	switch {
	case body.Content.SuiteId != s.props.Id:
		return errors.New("suite_id mismatch")
	case len(body.Content.SuiteTicket) == 0:
		return errors.New("suite_ticket is empty")
	}

	ticket := Ticket{ID: 1, Ticket: body.Content.SuiteTicket}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := s.db.NewInsert().
		Model(&ticket).
		On("CONFLICT (id) DO UPDATE").
		Set("ticket = EXCLUDED.ticket").
		Exec(ctx); err != nil {
		return err
	}

	return nil
}

type suiteTicketInfo struct {
	SuiteId     string `json:"suite_id"`
	InfoType    string `json:"info_type"`
	Timestamp   int64  `json:"timestamp"`
	SuiteTicket string `json:"suite_ticket"`
}
