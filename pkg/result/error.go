package result

import (
	"net/http"
	"time"
)

type ErrTitle string

const (
	Validation   ErrTitle = "Validation result"
	NotFound     ErrTitle = "Not Found"
	Conflict     ErrTitle = "Conflict"
	UnAuthorized ErrTitle = "UnAuthorized"
)

type Problem struct {
	Status    int       `json:"code"`
	Title     ErrTitle  `json:"title"`
	Detail    string    `json:"detail"`
	TimeStamp time.Time `json:"timeStamp"`
	Type      string    `json:"type"`
}

func (e *Problem) Error() string {
	return e.Detail
}

func NewErr(status int, detail string) *Problem {
	return &Problem{
		Status:    status,
		Title:     getErrTitle(status),
		Detail:    detail,
		Type:      getErrType(status),
		TimeStamp: time.Now(),
	}
}

func getErrTitle(code int) ErrTitle {
	switch code {
	case http.StatusBadRequest:
		return Validation
	case http.StatusUnauthorized:
		return UnAuthorized
	case http.StatusNotFound:
		return NotFound
	case http.StatusConflict:
		return Conflict
	}
	return "Internal Server Error"
}

func getErrType(code int) string {
	switch code {
	case http.StatusBadRequest:
		return "https://tools.ietf.org/html/rfc7807#section-3.1"
	case http.StatusUnauthorized:
		return "https://tools.ietf.org/html/rfc7235#section-3.1"
	case http.StatusNotFound:
		return "https://tools.ietf.org/html/rfc7231#section-6.5.4"
	case http.StatusConflict:
		return "https://tools.ietf.org/html/rfc7231#section-6.5.8"
	}
	return "https://tools.ietf.org/html/rfc7231#section-6.6.1"
}

func ValidationErr(detail string) *Problem {
	return NewErr(http.StatusBadRequest, detail)
}

func NotFoundErr(detail string) *Problem {
	return NewErr(http.StatusNotFound, detail)
}