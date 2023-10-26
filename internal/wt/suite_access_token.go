package wt

import (
	"absurdlab.io/WeSuiteCred/internal/httpx"
	"errors"
	"fmt"
	"github.com/peterbourgon/diskv/v3"
	"github.com/rs/zerolog"
	"sync"
	"time"
)

const (
	getSuiteAccessTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/service/get_suite_token"
)

func NewSuiteAccessTokenSupplier(
	props *Properties,
	logger *zerolog.Logger,
	store *diskv.Diskv,
) *SuiteAccessTokenSupplier {
	return &SuiteAccessTokenSupplier{
		logger: logger,
		props:  props,
		store:  store,
	}
}

type SuiteAccessTokenSupplier struct {
	sync.RWMutex

	logger *zerolog.Logger
	props  *Properties
	store  *diskv.Diskv

	accessToken string
	expiresAt   time.Time
}

func (s *SuiteAccessTokenSupplier) Get() (string, error) {
	token, err := s.doGet()
	if err != nil {
		return "", fmt.Errorf("failed to request suite_access_token: %w", err)
	}

	return token, nil
}

func (s *SuiteAccessTokenSupplier) doGet() (string, error) {
	if token, ok := s.tryReadToken(true); ok {
		return token, nil
	}

	s.Lock()
	defer s.Unlock()

	if token, ok := s.tryReadToken(false); ok {
		return token, nil
	}

	suiteTicket := s.store.ReadString(keySuiteTicket)
	if len(suiteTicket) == 0 {
		return "", errors.New("suite_ticket is empty")
	}

	var res getSuiteAccessTokenResponse
	if err := httpx.PostJson(getSuiteAccessTokenUrl, getSuiteAccessTokenRequest{}, &res); err != nil {
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

func (s *SuiteAccessTokenSupplier) tryReadToken(lock bool) (string, bool) {
	if lock {
		s.RLock()
		defer s.RUnlock()
	}

	switch {
	case len(s.accessToken) == 0:
	case s.expiresAt.IsZero():
	case time.Now().Add(s.props.SuiteAccessTokenLeeway).After(s.expiresAt):
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
