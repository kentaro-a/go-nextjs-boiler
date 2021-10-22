package error

import (
	"strings"
)

type Error struct {
	StatusCode int      `json:"status_code"`
	Messages   []string `json:"messages"`
}

func (err Error) Error() string {
	message := strings.Join(err.Messages, "\n")
	return message
}
func New(status_code int, messages []string) Error {
	return Error{
		status_code,
		messages,
	}
}
