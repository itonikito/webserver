package entities

import "time"

type BatRequest struct {
	Action string
	Param  string
}

type BatResponse struct {
	Success    bool      `json:"success"`
	Message    string    `json:"message"`
	Output     string    `json:"output"`
	Action     string    `json:"action,omitempty"`
	Parameter  string    `json:"parameter,omitempty"`
	DurationMs int64     `json:"duration_ms,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// Now возвращает текущее время - утилитарная функция
func Now() time.Time {
	return time.Now()
}
