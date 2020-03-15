package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"consumer/common/httpx"

	"github.com/yakaa/log4g"
)

func HttpRequest(method httpx.HttpMethod, url string, param interface{}) (bool, error) {

	bs, err := json.Marshal(param)
	if err != nil {
		return false, err
	}
	req, err := http.NewRequest(string(method), url, bytes.NewBuffer(bs))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpx.HttpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log4g.ErrorFormat("resp.Body.Close %+v", err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	return resp.Status == httpx.HttpStatus200 || strings.TrimSpace(string(body)) == httpx.HttpResponseSuccessStatus, nil
}
