package jwtLogout_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	jwtLogout "go-rest-api-auth/internal/handlers/auth/jwt/logout"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestJwtLogout(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		tokenManagerErr error
		expectedStatus  string
		expectedError   string
	}{
		{
			name:           "TestJwtLogout_NoContextUserID",
			userID:         "",
			expectedStatus: "Bad Request",
			expectedError:  "Invalid user context",
		},
		{
			name:            "TestJwtLogout_TokenManagerError",
			userID:          "123",
			tokenManagerErr: errors.New("token manager error"),
			expectedStatus:  "Bad Request",
			expectedError:   "Error deleting refresh token",
		},
		{
			name:           "TestJwtLogout_TokenDeleted",
			userID:         "123",
			expectedStatus: "OK",
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTokenManager := new(mocks.JwtManager)
			if tt.userID != "" {
				mockTokenManager.On("DeleteRefreshToken", 123).Return(tt.tokenManagerErr)
			}
			defer mockTokenManager.AssertExpectations(t)

			req := httptest.NewRequest(http.MethodPost, "/jwt_logout", nil)
			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "user_id", tt.userID)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			handler := jwtLogout.New(logger, mockTokenManager)
			handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			var respBody jwtLogout.Response
			err := json.NewDecoder(resp.Body).Decode(&respBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, respBody.Status)
			assert.Equal(t, tt.expectedError, respBody.Error)
		})
	}
}
