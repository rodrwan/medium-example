package logger

import (
	"encoding/json"
	"time"
)

type Logger struct {
	StatusCode int       `json:"status_code,omitempty"`
	Size       int       `json:"size,omitempty"`
	Method     string    `json:"method,omitempty"`
	TimeStamp  time.Time `json:"time_stamp,omitempty"`
	URL        string    `json:"url,omitempty"`
}

func (l *Logger) ToJSON() ([]byte, error) {
	return json.Marshal(l)
}
