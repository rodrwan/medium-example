package logger

import (
	"time"
)

type Logger struct {
	StatusCode int       `json:"status_code,omitempty"`
	Size       int       `json:"size,omitempty"`
	Method     string    `json:"method,omitempty"`
	TimeStamp  time.Time `json:"time_stamp,omitempty"`
	URL        string    `json:"url,omitempty"`
}
