package sessionLogin_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"go-rest-api-auth/internal/database"
	sessionLogin "go-rest-api-auth/internal/handlers/auth/session/login"
	"go-rest-api-auth/internal/utils"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestSessionLogin(t *testing.T) {
	tests := []struct {
		name              string
		reqBody           string
		userServiceErr    error
		sessionManagerErr error
		checkPasswordHash bool
		sessionExists     bool
		createSessionErr  error
		expectedStatus    string
		expectedError     string
	}{
		{
			name:           "TestSessionLogin_DecodeRequestBodyError",
			reqBody:        "{ invalid json }",
			expectedStatus: "Bad Request",
			expectedError:  "failed to decode request body",
		},
		{
			name:           "TestSessionLogin_ValidateRequestBodyError",
			reqBody:        "{\"username\":\"\",\"password\":\"\"}",
			expectedStatus: "Bad Request",
			expectedError:  "failed to validate request",
		},
		{
			name:           "TestSessionLogin_GetUserByNameError",
			reqBody:        "{\"username\":\"testuser\",\"password\":\"testpassword\"}",
			userServiceErr: errors.New("get user by name error"),
			expectedStatus: "Bad Request",
			expectedError:  "failed to get user",
		},
		{
			name:              "TestSessionLogin_InvalidPassword",
			reqBody:           "{\"username\":\"testuser\",\"password\":\"wrongpassword\"}",
			checkPasswordHash: false,
			expectedStatus:    "Bad Request",
			expectedError:     "invalid password",
		},
		{
			name:           "TestSessionLogin_SessionAlreadyExists",
			reqBody:        "{\"username\":\"testuser\",\"password\":\"testpassword\"}",
			sessionExists:  true,
			expectedStatus: "Bad Request",
			expectedError:  "session already exists",
		},
		{
			name:             "TestSessionLogin_CreateSessionError",
			reqBody:          "{\"username\":\"testuser\",\"password\":\"testpassword\"}",
			createSessionErr: errors.New("create session error"),
			expectedStatus:   "Bad Request",
			expectedError:    "create session error",
		},
		{
			name:           "TestSessionLogin_Success",
			reqBody:        "{\"username\":\"testuser\",\"password\":\"testpassword\"}",
			expectedStatus: "OK",
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionManager := new(mocks.SessionManager)
			mockUserService := new(mocks.UserService)
			log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			req, err := http.NewRequest(http.MethodPost, "/session_login", bytes.NewBuffer([]byte(tt.reqBody)))
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()

			handler := sessionLogin.New(log, mockSessionManager, mockUserService)

			mockSessionManager.On("GetterErrSessionNotFound").Return(errors.New("session not found"))
			mockSessionManager.On("GetterTtl").Return(time.Minute)

			if tt.userServiceErr != nil {
				mockUserService.On("GetUserByName", "testuser").Return(database.UserDTO{}, tt.userServiceErr)
			} else {
				hashedPassword, _ := utils.HashPassword("testpassword")
				mockUserService.On("GetUserByName", "testuser").Return(database.UserDTO{Id: 1, Username: "testuser", Password: hashedPassword}, nil)
			}

			if tt.sessionExists {
				mockSessionManager.On("GetSessionByUserID", "1").Return("", nil)
			} else {
				mockSessionManager.On("GetSessionByUserID", "1").Return("", mockSessionManager.GetterErrSessionNotFound())
			}

			if tt.createSessionErr != nil {
				mockSessionManager.On("CreateSession", "1").Return("", tt.createSessionErr)
			} else {
				mockSessionManager.On("CreateSession", "1").Return("session123", nil)
			}

			handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			var respBody sessionLogin.Response
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedStatus, respBody.Status)
			assert.Equal(t, tt.expectedError, respBody.Error)
		})
	}
}
