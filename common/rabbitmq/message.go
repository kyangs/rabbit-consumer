package rabbitmq

import "time"

type (
	Message struct {
		Data      interface{} `json:"data"`
		Timestamp time.Time   `json:"timestamp"`
		Type      string      `json:"type"`
		Url       string      `json:"url"`
		Delay     int64       `json:"delay"`
		Retry     bool        `json:"retry"`
		RetryTime int64       `json:"retryTime"`
		IsDelay   bool        `json:"isDelay"`
	}
)

const (
	ContentType      string = "Content-Type:application/json"
	TypePong         string = "pong"
	TypeMessage      string = "message"
	DefaultRetryTime int64  = 5000
	DefaultDelayKey  string = "x-delay"
)
