package httpx

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func PostJson(url string, request interface{}, response interface{}) (err error) {
	var bodyBytes []byte
	if request != nil {
		if bodyBytes, err = json.Marshal(request); err != nil {
			return
		}
	}

	res, err := http.Post(url, "application/json", bytes.NewReader(bodyBytes))
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
