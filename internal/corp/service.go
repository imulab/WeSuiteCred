package corp

import (
	"absurdlab.io/WeSuiteCred/internal/sqlitedb"
	"absurdlab.io/WeSuiteCred/internal/suite"
	"absurdlab.io/WeSuiteCred/internal/x"
	"context"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"net/url"
)

const (
	getPermanentCodeUrl = "https://qyapi.weixin.qq.com/cgi-bin/service/get_permanent_code"
	getAuthInfoUrl      = "https://qyapi.weixin.qq.com/cgi-bin/service/get_auth_info"
	getPermissionsUrl   = "https://qyapi.weixin.qq.com/cgi-bin/agent/get_permissions"
)

func NewService(db *bun.DB, suiteAccessToken *suite.AccessTokenSupplier) *Service {
	return &Service{db: db, suiteAccessToken: suiteAccessToken}
}

type Service struct {
	db               *bun.DB
	suiteAccessToken *suite.AccessTokenSupplier
}

func (s *Service) OnNewAuthCode(ctx context.Context, authCode string) error {
	codeWithAuthInfo, err := s.exchangeAuthCodeForPermanentCode(authCode)
	if err != nil {
		return err
	}

	permissions, err := s.getCorpPermissions(codeWithAuthInfo.AuthCorpInfo.CorpId, codeWithAuthInfo.PermanentCode)
	if err != nil {
		return err
	}

	record := Authorization{
		CorpID:        codeWithAuthInfo.AuthCorpInfo.CorpId,
		CorpName:      codeWithAuthInfo.AuthCorpInfo.CorpName,
		PermanentCode: codeWithAuthInfo.PermanentCode,
		AuthInfo:      sqlitedb.WrapJSON(codeWithAuthInfo.AuthInfo),
		Permissions:   sqlitedb.WrapJSON(*permissions),
	}

	if _, err = s.db.NewInsert().
		Model(&record).
		On("CONFLICT (corp_id) DO UPDATE").
		Set("corp_name = EXCLUDED.corp_name").
		Set("perm_code = EXCLUDED.perm_code").
		Set("auth_info = EXCLUDED.auth_info").
		Set("perm = EXCLUDED.perm").
		Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) OnAuthorizationChanged(ctx context.Context, corpId string) error {
	var record Authorization
	if err := s.db.NewSelect().
		Model(&record).
		Where("corp_id = ?", corpId).
		Limit(1).
		Scan(ctx, &record); err != nil {
		return fmt.Errorf("failed to get corp authz: %w", err)
	}

	authInfo, err := s.getAuthInfo(corpId, record.PermanentCode)
	if err != nil {
		return err
	}

	permissions, err := s.getCorpPermissions(corpId, record.PermanentCode)
	if err != nil {
		return err
	}

	record.AuthInfo = sqlitedb.WrapJSON(*authInfo)
	record.Permissions = sqlitedb.WrapJSON(*permissions)

	if _, err = s.db.NewUpdate().Model(&record).WherePK().Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) OnAuthorizationRemoved(ctx context.Context, corpId string) error {
	if _, err := s.db.NewDelete().
		Model(&Authorization{}).
		Where("corp_id = ?", corpId).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete corp authz: %w", err)
	}

	return nil
}

func (s *Service) getCorpPermissions(corpId string, corpSecret string) (*Permissions, error) {
	corpAccessToken, err := getCorpAccessToken(corpId, corpSecret)
	if err != nil {
		return nil, err
	}

	var res getPermissionsResponse
	if err = x.PostJson(
		getPermissionsUrl+"?"+url.Values{"access_token": []string{corpAccessToken}}.Encode(),
		nil,
		&res,
	); err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	if res.ErrorCode != 0 {
		return nil, fmt.Errorf("failed to get permissions: %s", res.ErrorMessage)
	}

	return &res.Permissions, nil
}

func (s *Service) exchangeAuthCodeForPermanentCode(authCode string) (*getPermanentCodeResponse, error) {
	suiteAccessToken, err := s.suiteAccessToken.Get()
	if err != nil {
		return nil, err
	}

	var res getPermanentCodeResponse
	if err = x.PostJson(
		getPermanentCodeUrl+"?"+url.Values{"suite_access_token": []string{suiteAccessToken}}.Encode(),
		getPermanentCodeRequest{AuthCode: authCode},
		&res,
	); err != nil {
		return nil, fmt.Errorf("failed to get permanent_code: %w", err)
	}

	switch {
	case res.ErrorCode != 0:
		return nil, fmt.Errorf("failed to get permanent_code: %s", res.ErrorMessage)
	case len(res.PermanentCode) == 0:
		return nil, errors.New("permanent_code is empty")
	}

	return &res, nil
}

func (s *Service) getAuthInfo(corpId string, permanentCode string) (*AuthInfo, error) {
	suiteAccessToken, err := s.suiteAccessToken.Get()
	if err != nil {
		return nil, err
	}

	var res getAuthInfoResponse
	if err = x.PostJson(
		getAuthInfoUrl+"?"+url.Values{"suite_access_token": []string{suiteAccessToken}}.Encode(),
		getAuthInfoRequest{AuthCorpId: corpId, PermanentCode: permanentCode},
		&res,
	); err != nil {
		return nil, fmt.Errorf("failed to get auth info: %w", err)
	}

	return &res.AuthInfo, nil
}

type getPermanentCodeRequest struct {
	AuthCode string `json:"auth_code_value"`
}

type getPermanentCodeResponse struct {
	AuthInfo

	ErrorCode     int    `json:"errcode"`
	ErrorMessage  string `json:"errmsg"`
	PermanentCode string `json:"permanent_code"`
	State         string `json:"state"`
}

type getAuthInfoRequest struct {
	AuthCorpId    string `json:"auth_corpid"`
	PermanentCode string `json:"permanent_code"`
}

type getAuthInfoResponse struct {
	AuthInfo

	ErrorCode    int    `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}

type getPermissionsResponse struct {
	Permissions

	ErrorCode    int    `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}
