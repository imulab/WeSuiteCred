package x

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func PostJson(apiUrl string, request interface{}, response interface{}) (err error) {
	var bodyBytes []byte
	if request != nil {
		if bodyBytes, err = json.Marshal(request); err != nil {
			return
		}
	}

	if shouldDebug, _ := strconv.ParseBool(os.Getenv("WSC_DEBUG")); shouldDebug {
		u, _ := url.Parse(apiUrl)

		q := u.Query()
		q.Add("debug", "1")

		u.RawQuery = q.Encode()
		apiUrl = u.String()
	}

	res, err := http.Post(apiUrl, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return
	}

	return
}
