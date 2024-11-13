package jwtLogin_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"go-rest-api-auth/internal/database"
	jwtLogin "go-rest-api-auth/internal/handlers/auth/jwt/login"
	"go-rest-api-auth/internal/utils"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestJwtLogin(t *testing.T) {
	tests := []struct {
		name                string
		reqBody             string
		userServiceErr      error
		tokenGenAccessErr   error
		tokenGenRefreshErr  error
		saveRefreshTokenErr error
		checkPasswordHash   bool
		expectedStatus      string
		expectedError       string
	}{
		{
			name:           "TestJwtLogin_DecodeRequestBodyError",
			reqBody:        "{ invalid json }",
			expectedStatus: "Bad Request",
			expectedError:  "failed to decode request body",
		},
		{
			name:           "TestJwtLogin_ValidateRequestBodyError",
			reqBody:        "{\"username\":\"\",\"password\":\"\"}",
			expectedStatus: "Bad Request",
			expectedError:  "failed to validate request",
		},
		{
			name:           "TestJwtLogin_GetUserByNameError",
			reqBody:        "{\"username\":\"testuser\",\"password\":\"testpassword\"}",
			userServiceErr: errors.New("get user by name error"),
			expectedStatus: "Bad Request",
			expectedError:  "failed to get user",
		},
		{
			name:              "TestJwtLogin_InvalidPassword",
			reqBody:           "{\"username\":\"testuser\",\"password\":\"wrongpassword\"}",
			checkPasswordHash: false,
			expectedStatus:    "Bad Request",
			expectedError:     "invalid password",
		},
		{
			name:              "TestJwtLogin_AccessTokenGenerationError",
			reqBody:           "{\"username\":\"testuser\",\"password\":\"testpassword\"}",
			tokenGenAccessErr: errors.New("generate access token error"),
			expectedStatus:    "Bad Request",
			expectedError:     "failed to generate access token",
		},
		{
			name:               "TestJwtLogin_RefreshTokenGenerationError",
			reqBody:            "{\"username\":\"testuser\",\"password\":\"testpassword\"}",
			tokenGenRefreshErr: errors.New("generate refresh token error"),
			expectedStatus:     "Bad Request",
			expectedError:      "failed to generate refresh token",
		},
		{
			name:                "TestJwtLogin_SaveRefreshTokenError",
			reqBody:             "{\"username\":\"testuser\",\"password\":\"testpassword\"}",
			saveRefreshTokenErr: errors.New("save refresh token error"),
			expectedStatus:      "Bad Request",
			expectedError:       "failed to save refresh token",
		},
		{
			name:           "TestJwtLogin_Success",
			reqBody:        "{\"username\":\"testuser\",\"password\":\"testpassword\"}",
			expectedStatus: "OK",
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTokenManager := new(mocks.JwtManager)
			mockUserService := new(mocks.UserService)
			log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			req, err := http.NewRequest(http.MethodPost, "/jwt_login", bytes.NewBuffer([]byte(tt.reqBody)))
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			handler := jwtLogin.New(log, mockTokenManager, mockUserService)

			mockTokenManager.On("GetterAccessExpiresAt").Return(time.Minute)
			mockTokenManager.On("GetterRefreshExpiresAt").Return(time.Hour)

			if tt.userServiceErr != nil {
				mockUserService.On("GetUserByName", "testuser").Return(database.UserDTO{}, tt.userServiceErr)
			} else {
				hashedPassword, _ := utils.HashPassword("testpassword")
				mockUserService.On("GetUserByName", "testuser").Return(database.UserDTO{Id: 1, Username: "testuser", Password: hashedPassword}, nil)
			}

			if tt.checkPasswordHash == false {
				mockUserService.On("CheckPasswordHash", "wrongpassword", "testpassword").Return(false)
			} else {
				mockUserService.On("CheckPasswordHash", "testpassword", "testpassword").Return(true)
			}

			if tt.tokenGenAccessErr != nil {
				mockTokenManager.On("GenerateJWT", "1", "access", mockTokenManager.GetterAccessExpiresAt()).Return("", tt.tokenGenAccessErr)
			} else {
				mockTokenManager.On("GenerateJWT", "1", "access", mockTokenManager.GetterAccessExpiresAt()).Return("access123", nil)
			}

			if tt.tokenGenRefreshErr != nil {
				mockTokenManager.On("GenerateJWT", "1", "refresh", mockTokenManager.GetterRefreshExpiresAt()).Return("", tt.tokenGenRefreshErr)
			} else {
				mockTokenManager.On("GenerateJWT", "1", "refresh", mockTokenManager.GetterRefreshExpiresAt()).Return("refresh123", nil)
			}

			if tt.saveRefreshTokenErr != nil {
				mockTokenManager.On("SaveRefreshToken", "refresh123").Return(tt.saveRefreshTokenErr)
			} else {
				mockTokenManager.On("SaveRefreshToken", "refresh123").Return(nil)
			}

			handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			var respBody jwtLogin.Response
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedStatus, respBody.Status)
			assert.Equal(t, tt.expectedError, respBody.Error)
		})
	}
}
