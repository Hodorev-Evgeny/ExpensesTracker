package core_middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	core_repository_redis_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/redis/pool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ValueSessionDb map[string]string = map[string]string{}

type MockRedisClient struct {
}

type TestStringCmd struct {
	ans string
}

func (s *TestStringCmd) Result() (string, error) {
	if s.ans == "" {
		return "", redis.Nil
	}
	return s.ans, nil
}

func (r *MockRedisClient) Ping(ctx context.Context) core_repository_redis_pool.CustomStatusCmd {
	return nil
}

func (r *MockRedisClient) Get(ctx context.Context, key string) core_repository_redis_pool.CustomStringCmd {
	if val, ok := ValueSessionDb[key]; ok {
		return &TestStringCmd{ans: val}
	}
	return &TestStringCmd{ans: ""}
}

func (r *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) core_repository_redis_pool.CustomStatusCmd {
	return nil
}

func (r *MockRedisClient) Del(ctx context.Context, key string) core_repository_redis_pool.CustomIntCmd {
	return nil
}

func (r *MockRedisClient) GetTTL() time.Duration {
	return time.Second
}

func (r *MockRedisClient) GetTimeout() time.Duration {
	return time.Second
}

func (r *MockRedisClient) Close() error {
	return nil
}

func TestAuthenticator(t *testing.T) {
	redisMock := &MockRedisClient{}
	ValueSessionDb["sessionID:testParam"] = "123"
	logger := &core_logger.Logger{
		Logger: zap.NewNop(),
	}

	t.Run("Correct authenticate", func(t *testing.T) {
		var receivedUserid int
		var CheckInCtx bool

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			receivedUserid, CheckInCtx = core_domain.UserIDFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		middleware := Authenticator(redisMock)
		handlerTest := middleware(nextHandler)

		r := httptest.NewRequest(http.MethodGet, "/api", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(core_logger.ToContext(r.Context(), logger))

		r.AddCookie(&http.Cookie{
			Name:  "sessionID",
			Value: "testParam",
		})

		handlerTest.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
		}

		expectedUserid, _ := strconv.Atoi(ValueSessionDb["sessionID:testParam"])
		if !CheckInCtx {
			t.Errorf("UserID from context was not set")
		}

		if expectedUserid != receivedUserid {
			t.Errorf("UserID from context did not match: got %v want %v", receivedUserid, expectedUserid)
		}
	})

	t.Run("Authenticate without cookie", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := Authenticator(redisMock)
		handlerTest := middleware(nextHandler)

		r := httptest.NewRequest(http.MethodGet, "/api", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(core_logger.ToContext(r.Context(), logger))

		handlerTest.ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Got status code %v want %v", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("Authenticate with cookie for which there is no session in Redis", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := Authenticator(redisMock)
		handlerTest := middleware(nextHandler)

		r := httptest.NewRequest(http.MethodGet, "/api", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(core_logger.ToContext(r.Context(), logger))

		r.AddCookie(&http.Cookie{
			Name:  "sessionID",
			Value: "NotParam",
		})

		handlerTest.ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Got status code %v want %v", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("Test whilelist in authenticate without cookie for", func(t *testing.T) {
		tests := []struct {
			name   string
			method string
			path   string
			want   int
		}{
			{"register", http.MethodGet, "/register", http.StatusOK},
			{"login", http.MethodGet, "/login", http.StatusOK},
			{"css", http.MethodGet, "/css/style.css", http.StatusOK},
			{"js", http.MethodGet, "/js/app.js", http.StatusOK},
			{"swagger", http.MethodGet, "/swagger/index.html", http.StatusOK},
			{"health", http.MethodGet, "/api/v1/health", http.StatusOK},
			{"get users", http.MethodGet, "/api/v1/users", http.StatusUnauthorized},
			{"register user", http.MethodPost, "/api/v1/users/register", http.StatusOK},
		}

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := Authenticator(redisMock)
		handlerTest := middleware(nextHandler)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := httptest.NewRequest(tt.method, tt.path, nil)
				r = r.WithContext(core_logger.ToContext(r.Context(), logger))

				w := httptest.NewRecorder()

				handlerTest.ServeHTTP(w, r)

				if w.Code != tt.want {
					t.Errorf("Got status code %v want %v", w.Code, tt.want)
				}
			})
		}
	})
}
