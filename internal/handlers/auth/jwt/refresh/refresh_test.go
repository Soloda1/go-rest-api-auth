package refresh_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-rest-api-auth/internal/handlers/auth/jwt/refresh"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestRefreshHandler(t *testing.T) {
	tests := []struct {
		name                string
		requestBody         interface{}
		validateJWTErr      error
		isRefreshTokenErr   error
		isRefreshTokenValid bool
		deleteTokenErr      error
		generateAccessErr   error
		generateRefreshErr  error
		saveTokenErr        error
		expectedStatus      string
		expectedResponse    refresh.Response
	}{
		{
			name:             "InvalidRequestBody",
			requestBody:      "invalid_json",
			expectedStatus:   "Bad Request",
			expectedResponse: refresh.Response{Status: "Bad Request", Error: "failed to decode request body"},
		},
		{
			name: "InvalidToken",
			requestBody: refresh.Request{
				RefreshToken: "invalid_token",
			},
			validateJWTErr:   errors.New("invalid token"),
			expectedStatus:   "Bad Request",
			expectedResponse: refresh.Response{Status: "Bad Request", Error: "failed to validate token"},
		},
		{
			name: "TokenNotFound",
			requestBody: refresh.Request{
				RefreshToken: "valid_token",
			},
			isRefreshTokenValid: false,
			expectedStatus:      "Bad Request",
			expectedResponse:    refresh.Response{Status: "Bad Request", Error: "refresh token not found"},
		},
		{
			name: "Success",
			requestBody: refresh.Request{
				RefreshToken: "valid_token",
			},
			isRefreshTokenValid: true,
			expectedStatus:      "OK",
			expectedResponse: refresh.Response{
				Status:       "OK",
				AccessToken:  "new_access_token",
				RefreshToken: "new_refresh_token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTokenManager := new(mocks.JwtManager)

			mockTokenManager.On("GetterAccessExpiresAt").Return(time.Minute)
			mockTokenManager.On("GetterRefreshExpiresAt").Return(time.Hour)

			mockTokenManager.On("ValidateJWT", mock.Anything, "refresh").Return(jwt.MapClaims{"sub": "123"}, tt.validateJWTErr)
			mockTokenManager.On("IsRefreshTokenValid", mock.Anything).Return(tt.isRefreshTokenValid, tt.isRefreshTokenErr)
			mockTokenManager.On("DeleteRefreshToken", mock.Anything).Return(tt.deleteTokenErr)
			mockTokenManager.On("GenerateJWT", mock.Anything, "access", mock.Anything).Return("new_access_token", tt.generateAccessErr)
			mockTokenManager.On("GenerateJWT", mock.Anything, "refresh", mock.Anything).Return("new_refresh_token", tt.generateRefreshErr)
			mockTokenManager.On("SaveRefreshToken", mock.Anything).Return(tt.saveTokenErr)

			var body []byte
			if tt.requestBody != nil {
				body, _ = json.Marshal(tt.requestBody)
			}
			req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			handler := refresh.New(logger, mockTokenManager)
			handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			var respBody refresh.Response
			json.NewDecoder(resp.Body).Decode(&respBody)

			assert.Equal(t, tt.expectedStatus, respBody.Status)
			assert.Equal(t, tt.expectedResponse, respBody)
		})
	}
}
