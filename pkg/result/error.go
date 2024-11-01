package results

import (
	"net/http"
	"time"
)

type ErrType string

const (
	Validation   ErrType = "Validation result"
	NotFound     ErrType = "Not Found"
	Conflict     ErrType = "Conflict"
	UnAuthorized ErrType = "UnAuthorized"
)

type Problem struct {
	Code      int       `response:"code"`
	Title     string    `response:"title"`
	Detail    string    `response:"detail"`
	TimeStamp time.Time `response:"timeStamp"`
	Type      ErrType   `response:"type"`
}

func (e *Problem) Error() string {
	return e.Title
}

func NewErr(code int, title string, detail string) *Problem {
	return &Problem{
		Code:      code,
		Title:     title,
		Detail:    detail,
		TimeStamp: time.Now(),
		Type:      CreateErrType(code),
	}
}

func NewErrTyped(code int, title string, detail string, errType ErrType) *Problem {
	if errType == "" {
		errType = CreateErrType(code)
	}
	return &Problem{
		Code:      code,
		Title:     title,
		Detail:    detail,
		TimeStamp: time.Now(),
		Type:      errType,
	}
}

func CreateErrType(code int) ErrType {
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

func ValidationError(title string, detail string) *Problem {
	return NewErrTyped(400, title, detail, Validation)
}

func NotFoundError(title string, detail string) *Problem {
	return NewErrTyped(400, title, detail, NotFound)
}
