package sessionLogout_test

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	sessionLogout "go-rest-api-auth/internal/handlers/auth/session/logout"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSessionLogout(t *testing.T) {
	tests := []struct {
		name              string
		sessionID         string
		sessionManagerErr error
		expectedStatus    string
		expectedError     string
	}{
		{
			name:           "TestSessionLogout_NoSessionCookie",
			sessionID:      "",
			expectedStatus: "Bad Request",
			expectedError:  "Error getting cookie",
		},
		{
			name:              "TestSessionLogout_SessionManagerError",
			sessionID:         "session123",
			sessionManagerErr: errors.New("session manager error"),
			expectedStatus:    "Bad Request",
			expectedError:     "Error deleting session",
		},
		{
			name:           "TestSessionLogout_SessionDeleted",
			sessionID:      "session123",
			expectedStatus: "OK",
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Создание моков
			mockSessionManager := new(mocks.SessionManager)
			if tt.sessionID != "" {
				mockSessionManager.On("DeleteSession", tt.sessionID).Return(tt.sessionManagerErr)
			}
			defer mockSessionManager.AssertExpectations(t)

			req := httptest.NewRequest(http.MethodGet, "/session_logout", nil)
			if tt.sessionID != "" {
				req.AddCookie(&http.Cookie{
					Name:  "session_id",
					Value: tt.sessionID,
				})
			}
			w := httptest.NewRecorder()

			handler := sessionLogout.New(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})), mockSessionManager)
			handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			var respBody sessionLogout.Response
			err := json.NewDecoder(resp.Body).Decode(&respBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, respBody.Status)
			assert.Equal(t, tt.expectedError, respBody.Error)
		})
	}
}
