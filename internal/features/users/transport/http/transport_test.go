package features_users_transport

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	transport_user_mock "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/users/transport/http/mock"
	"go.uber.org/zap"
)

type DataTestLogin struct {
	name               string
	body               string
	expectedStatusCode int
	withError          bool
}

func TestCreateUser(t *testing.T) {
	mockService := transport_user_mock.NewMockService()
	handler := NewUserHTTPHandler(mockService, false, 3*time.Second)
	logger := &core_logger.Logger{
		Logger: zap.NewNop(),
	}

	tests := []DataTestLogin{
		{
			name:               "test empty body",
			body:               `{}`,
			expectedStatusCode: http.StatusBadRequest,
			withError:          true,
		},
		{
			name: "password too short",
			body: `{
			"full_name": "John Doe",
			"email": "test@gmail.com",
			"password": "a",
			"phone": "+71124234312"
		}`,
			expectedStatusCode: http.StatusBadRequest,
			withError:          true,
		},
		{
			name: "valid json body",
			body: `{
			"full_name": "John Doe",
			"email": "test@gmail.com",
			"password": "password123",
			"phone": "+71124234312"
		}`,
			expectedStatusCode: http.StatusCreated,
			withError:          false,
		},
		{
			name: "broken json body",
			body: `{
			"full_name": "John Doe",
			"email": "test@gmail.com",
			"password": "password123",
			"phone": "+71124234312"`,
			expectedStatusCode: http.StatusBadRequest,
			withError:          true,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(test.body))
		w := httptest.NewRecorder()
		ctx := core_logger.ToContext(req.Context(), logger)

		handler.CreateUser(w, req.WithContext(ctx))

		if w.Code != test.expectedStatusCode {
			t.Errorf("expected status code %d, got %d", test.expectedStatusCode, w.Code)
		}

		if !test.withError {
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected application/json, got %s", contentType)
			}
		}
	}
}
