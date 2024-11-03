package result

import (
	"net/http"
)

type ErrTitle string

const (
	Validation     ErrTitle = "Validation problem"
	NotFound       ErrTitle = "Not Found"
	Conflict       ErrTitle = "Conflict"
	UnAuthorized   ErrTitle = "Unauthorized"
	GatewayTimeout ErrTitle = "Request Timeout"
)

type Err struct {
	Status int      `json:"code"`
	Title  ErrTitle `json:"title"`
	Detail string   `json:"detail"`
	Type   string   `json:"type"`
}

func (e *Err) Error() string {
	return e.Detail
}

func NewErr(status int, detail string) *Err {
	return &Err{
		Status: status,
		Title:  getErrTitle(status),
		Detail: detail,
		Type:   getErrType(status),
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
	case http.StatusGatewayTimeout:
		return GatewayTimeout

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

func ValidationErr(detail string) *Err {
	return NewErr(http.StatusBadRequest, detail)
}

func NotFoundErr(detail string) *Err {
	return NewErr(http.StatusNotFound, detail)
}

func InternalServerErr(detail string) *Err {
	return NewErr(http.StatusInternalServerError, detail)
}

func TimeoutErr() *Err {
	return NewErr(http.StatusGatewayTimeout, "request timeout")
}

func UnauthorizedErr(detail string) error {
	return NewErr(http.StatusUnauthorized, detail)
}
