package result

import (
	"net/http"
	"testing"
)

func TestNewErr(t *testing.T) {
	tests := []struct {
		status   int
		detail   string
		expected Err
	}{
		{
			status: http.StatusBadRequest,
			detail: "Invalid input data",
			expected: Err{
				Status: http.StatusBadRequest,
				Title:  Validation,
				Detail: "Invalid input data",
			},
		},
		{
			status: http.StatusUnauthorized,
			detail: "Unauthorized access",
			expected: Err{
				Status: http.StatusUnauthorized,
				Title:  UnAuthorized,
				Detail: "Unauthorized access",
			},
		},
		{
			status: http.StatusNotFound,
			detail: "Resource not found",
			expected: Err{
				Status: http.StatusNotFound,
				Title:  NotFound,
				Detail: "Resource not found",
			},
		},
		{
			status: http.StatusConflict,
			detail: "Resource conflict",
			expected: Err{
				Status: http.StatusConflict,
				Title:  Conflict,
				Detail: "Resource conflict",
			},
		},
		{
			status: http.StatusInternalServerError,
			detail: "Unexpected error",
			expected: Err{
				Status: http.StatusInternalServerError,
				Title:  "Internal Server Error",
				Detail: "Unexpected error",
			},
		},
	}

	for _, tt := range tests {
		err := NewErr(tt.status, tt.detail)

		if err.Status != tt.expected.Status {
			t.Errorf("Expected Status %d, got %d", tt.expected.Status, err.Status)
		}
		if err.Title != tt.expected.Title {
			t.Errorf("Expected Title %s, got %s", tt.expected.Title, err.Title)
		}
		if err.Detail != tt.expected.Detail {
			t.Errorf("Expected Detail %s, got %s", tt.expected.Detail, err.Detail)
		}
	}
}

func TestValidationErr(t *testing.T) {
	detail := "Invalid input"
	err := ValidationErr(detail)

	if err.Status != http.StatusBadRequest {
		t.Errorf("Expected Status %d, got %d", http.StatusBadRequest, err.Status)
	}
}

func TestInternalServerErr(t *testing.T) {
	detail := "Server error"
	err := InternalServerErr(detail)

	if err.Status != http.StatusInternalServerError {
		t.Errorf("Expected Status %d, got %d", http.StatusInternalServerError, err.Status)
	}
}

func TestNotFoundErr(t *testing.T) {
	detail := "Item not found"
	err := NotFoundErr(detail)

	if err.Status != http.StatusNotFound {
		t.Errorf("Expected Status %d, got %d", http.StatusNotFound, err.Status)
	}
}

func TestTimeoutErr(t *testing.T) {
	err := TimeoutErr()

	if err.Status != http.StatusGatewayTimeout {
		t.Errorf("Expected Status %d, got %d", http.StatusGatewayTimeout, err.Status)
	}
}

func TestUnauthorizedErr(t *testing.T) {
	detail := "Unauthorized access"
	err := UnauthorizedErr(detail)

	if err.Status != http.StatusUnauthorized {
		t.Errorf("Expected Status %d, got %d", http.StatusUnauthorized, err.Status)
	}
}
