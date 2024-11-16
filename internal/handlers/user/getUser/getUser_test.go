package getUser_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-rest-api-auth/internal/database"
	"go-rest-api-auth/internal/handlers/user/getUser"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockResponse   database.UserDTO
		mockError      error
		expectedStatus string
		expectedBody   getUser.Response
	}{
		{
			name:   "SuccessfulGetUser",
			userID: "1",
			mockResponse: database.UserDTO{
				Id:          1,
				Username:    "testuser",
				Password:    "testpass",
				Description: "test description",
				DateJoined:  pgtype.Date{},
			},
			mockError:      nil,
			expectedStatus: "OK",
			expectedBody: getUser.Response{
				Status: "OK",
				User: database.UserDTO{
					Id:          1,
					Username:    "testuser",
					Password:    "testpass",
					Description: "test description",
					DateJoined:  pgtype.Date{},
				},
			},
		},
		{
			name:           "InvalidUserID",
			userID:         "abc", // Невалидный ID
			mockResponse:   database.UserDTO{},
			mockError:      nil,
			expectedStatus: "Bad Request",
			expectedBody: getUser.Response{
				Status: "Bad Request",
				Error:  "Invalid user id",
			},
		},
		{
			name:           "UserNotFound",
			userID:         "2",
			mockResponse:   database.UserDTO{},
			mockError:      errors.New("user not found"),
			expectedStatus: "Bad Request",
			expectedBody: getUser.Response{
				Status: "Bad Request",
				Error:  "User not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.UserService)
			if tt.name != "InvalidUserID" {
				mockService.On("GetUserById", mock.AnythingOfType("int")).Return(tt.mockResponse, tt.mockError)
			}
			defer mockService.AssertExpectations(t)

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			handler := getUser.New(logger, mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /users/{userID}", handler)

			server := httptest.NewServer(mux)
			defer server.Close()

			url := server.URL + "/users/" + tt.userID

			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			req = req.WithContext(context.WithValue(req.Context(), "user_id", "123"))

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			var responseBody getUser.Response
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, responseBody.Status)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
