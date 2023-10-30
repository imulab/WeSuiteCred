package suite

import (
	"absurdlab.io/WeSuiteCred/internal/x"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"sync"
	"time"
)

const (
	getSuiteAccessTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/service/get_suite_token"
)

func NewAccessTokenSupplier(
	props *Properties,
	logger *zerolog.Logger,
	db *bun.DB,
) *AccessTokenSupplier {
	return &AccessTokenSupplier{
		logger: logger,
		props:  props,
		db:     db,
	}
}

type AccessTokenSupplier struct {
	sync.RWMutex

	logger *zerolog.Logger
	props  *Properties
	db     *bun.DB

	accessToken string
	expiresAt   time.Time
}

func (s *AccessTokenSupplier) Get() (string, error) {
	token, err := s.doGet()
	if err != nil {
		return "", fmt.Errorf("failed to request suite_access_token: %w", err)
	}

	return token, nil
}

func (s *AccessTokenSupplier) Reset() {
	s.Lock()
	s.accessToken = ""
	s.expiresAt = time.Time{}
	s.Unlock()
}

func (s *AccessTokenSupplier) getLatestTicket() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var ticket Ticket
	if err := s.db.NewSelect().Model(&ticket).Where("id = ?", int64(1)).Scan(ctx); err != nil {
		return "", fmt.Errorf("failed to get suite_ticket: %w", err)
	}

	return ticket.Ticket, nil
}

func (s *AccessTokenSupplier) doGet() (string, error) {
	if token, ok := s.tryReadToken(true); ok {
		return token, nil
	}

	s.Lock()
	defer s.Unlock()

	if token, ok := s.tryReadToken(false); ok {
		return token, nil
	}

	suiteTicket, err := s.getLatestTicket()
	if err != nil {
		return "", err
	}

	var res getSuiteAccessTokenResponse
	if err = x.PostJson(
		getSuiteAccessTokenUrl,
		getSuiteAccessTokenRequest{
			SuiteId:     s.props.Id,
			SuiteSecret: s.props.Secret,
			SuiteTicket: suiteTicket,
		},
		&res,
	); err != nil {
		return "", err
	}

	if res.ErrorCode != 0 {
		s.logger.Error().
			Int("code", res.ErrorCode).
			Str("message", res.ErrorMessage).
			Msg("getSuiteAccessToken error")
		return "", errors.New(res.ErrorMessage)
	}

	s.accessToken = res.SuiteAccessToken
	s.expiresAt = time.Now().Add(time.Duration(res.ExpiresIn) * time.Second)

	return s.accessToken, nil
}

func (s *AccessTokenSupplier) tryReadToken(lock bool) (string, bool) {
	if lock {
		s.RLock()
		defer s.RUnlock()
	}

	switch {
	case len(s.accessToken) == 0:
	case s.expiresAt.IsZero():
	case time.Now().Add(s.props.AccessTokenLeeway).After(s.expiresAt):
	default:
		return s.accessToken, true
	}

	return "", false
}

type getSuiteAccessTokenRequest struct {
	SuiteId     string `json:"suite_id"`
	SuiteSecret string `json:"suite_secret"`
	SuiteTicket string `json:"suite_ticket"`
}

type getSuiteAccessTokenResponse struct {
	ErrorCode        int    `json:"errcode"`
	ErrorMessage     string `json:"errmsg"`
	SuiteAccessToken string `json:"suite_access_token"`
	ExpiresIn        int64  `json:"expires_in"`
}
