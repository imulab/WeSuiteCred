package wt

import (
	"absurdlab.io/WeSuiteCred/internal/httpx"
	"errors"
	"fmt"
	"net/url"
)

const (
	getPermanentCodeUrl = "https://qyapi.weixin.qq.com/cgi-bin/service/get_permanent_code"
	getAuthInfoUrl      = "https://qyapi.weixin.qq.com/cgi-bin/service/get_auth_info"
)

func NewCorpService(
	accessToken *SuiteAccessTokenSupplier,
	corpSecretDao *CorpSecretDao,
	corpAuthInfoDao *CorpAuthInfoDao,
) *CorpService {
	return &CorpService{
		accessToken: accessToken,
		secretDao:   corpSecretDao,
		authInfoDao: corpAuthInfoDao,
	}
}

type CorpService struct {
	accessToken *SuiteAccessTokenSupplier
	secretDao   *CorpSecretDao
	authInfoDao *CorpAuthInfoDao
}

func (s *CorpService) UpdateSecret(authCode string) error {
	suiteAccessToken, err := s.accessToken.Get()
	if err != nil {
		return err
	}

	var res getPermanentCodeResponse
	if err = httpx.PostJson(
		getPermanentCodeUrl+"?"+url.Values{"suite_access_token": []string{suiteAccessToken}}.Encode(),
		getPermanentCodeRequest{AuthCode: authCode},
		&res,
	); err != nil {
		return fmt.Errorf("failed to get permanent_code: %w", err)
	}

	switch {
	case res.ErrorCode != 0:
		return fmt.Errorf("failed to get permanent_code: %s", res.ErrorMessage)
	case len(res.PermanentCode) == 0:
		return errors.New("permanent_code is empty")
	}

	if err = s.secretDao.Write(res.AuthCorpInfo.CorpId, res.PermanentCode); err != nil {
		return err
	}

	if err = s.authInfoDao.Write(&res.AuthInfo); err != nil {
		return err
	}

	return nil
}

func (s *CorpService) UpdateAuthInfo(authCorpId string) error {
	suiteAccessToken, err := s.accessToken.Get()
	if err != nil {
		return err
	}

	code, err := s.secretDao.Get(authCorpId)
	if err != nil {
		return err
	}

	var res getAuthInfoResponse
	if err = httpx.PostJson(
		getAuthInfoUrl+"?"+url.Values{"suite_access_token": []string{suiteAccessToken}}.Encode(),
		getAuthInfoRequest{AuthCorpId: authCorpId, PermanentCode: code},
		&res,
	); err != nil {
		return fmt.Errorf("failed to get auth info: %w", err)
	}

	return s.authInfoDao.Write(&res.AuthInfo)
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
