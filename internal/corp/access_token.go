package corp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	getCorpAccessTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
)

// getCorpAccessToken retrieves the access token for the given corp identified by the credentials.
//
// We are NOT reusing the access token here because the frequency at which it is required is very low and will not
// trigger the rate limit. Hence, it is not worth the effort to cache it.
func getCorpAccessToken(corpId string, corpSecret string) (accessToken string, err error) {
	q := url.Values{}
	{
		q.Add("corpid", corpId)
		q.Add("corpsecret", corpSecret)
	}

	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to get corp access token: %w", err)
			return
		}
	}()

	res, err := http.Get(getCorpAccessTokenUrl + "?" + q.Encode())
	if err != nil {
		return "", err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var response getCorpAccessTokenResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", err
	}

	if response.ErrorCode != 0 {
		return "", errors.New(response.ErrorMessage)
	}

	accessToken = response.AccessToken

	return
}

type getCorpAccessTokenResponse struct {
	ErrorCode    int    `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
