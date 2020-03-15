package httpx

import (
	"net/http"
	"time"
)

type (
	HttpMethod string
)

var (
	HttpClient = &http.Client{
		Timeout: 3 * time.Second,
	}
)

const (
	HttpMethodPost   HttpMethod = "POST"
	HttpMethodGet    HttpMethod = "GET"
	HttpMethodDelete HttpMethod = "DELETE"

	HttpStatus200             string = "200 OK"
	HttpResponseSuccessStatus string = "SUCCESS"
	HttpResponseFailStatus    string = "FAIL"
)
