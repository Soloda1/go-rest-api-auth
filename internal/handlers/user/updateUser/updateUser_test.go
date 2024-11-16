package updateUser_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-rest-api-auth/internal/database"
	"go-rest-api-auth/internal/handlers/user/updateUser"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUpdateUserHandler(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		requestBody     updateUser.Request
		mockGetResponse database.UserDTO
		mockGetError    error
		mockUpdateError error
		expectedStatus  string
		expectedBody    updateUser.Response
	}{
		{
			name:   "SuccessfulUpdateUser",
			userID: "1",
			requestBody: updateUser.Request{
				Username:    "updateduser",
				Password:    "updatedpass",
				Description: "updated description",
			},
			mockGetResponse: database.UserDTO{
				Id:          1,
				Username:    "testuser",
				Password:    "testpass",
				Description: "test description",
				DateJoined:  pgtype.Date{},
			},
			mockGetError:    nil,
			mockUpdateError: nil,
			expectedStatus:  "OK",
			expectedBody: updateUser.Response{
				Status: "OK",
				User: database.UserDTO{
					Id:          1,
					Username:    "updateduser",
					Password:    "updatedpass",
					Description: "updated description",
					DateJoined:  pgtype.Date{},
				},
			},
		},
		{
			name:            "InvalidUserID",
			userID:          "abc", // Невалидный ID
			requestBody:     updateUser.Request{},
			mockGetResponse: database.UserDTO{},
			mockGetError:    nil,
			mockUpdateError: nil,
			expectedStatus:  "Bad Request",
			expectedBody: updateUser.Response{
				Status: "Bad Request",
				Error:  "Invalid user id",
			},
		},
		{
			name:   "UserNotFound",
			userID: "2",
			requestBody: updateUser.Request{
				Username:    "updateduser",
				Password:    "updatedpass",
				Description: "updated description",
			},
			mockGetResponse: database.UserDTO{},
			mockGetError:    errors.New("user not found"),
			mockUpdateError: nil,
			expectedStatus:  "Bad Request",
			expectedBody: updateUser.Response{
				Status: "Bad Request",
				Error:  "User not found",
			},
		},
		{
			name:   "ErrorUpdatingUser",
			userID: "1",
			requestBody: updateUser.Request{
				Username:    "updateduser",
				Password:    "updatedpass",
				Description: "updated description",
			},
			mockGetResponse: database.UserDTO{
				Id:          1,
				Username:    "testuser",
				Password:    "testpass",
				Description: "test description",
				DateJoined:  pgtype.Date{},
			},
			mockGetError:    nil,
			mockUpdateError: errors.New("update error"),
			expectedStatus:  "Bad Request",
			expectedBody: updateUser.Response{
				Status: "Bad Request",
				Error:  "failed to update user",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.UserService)
			if tt.name == "UserNotFound" {
				mockService.On("GetUserById", mock.AnythingOfType("int")).Return(tt.mockGetResponse, tt.mockGetError)
			} else if tt.name != "InvalidUserID" {
				mockService.On("GetUserById", mock.AnythingOfType("int")).Return(tt.mockGetResponse, tt.mockGetError)
				mockService.On("UpdateUser", mock.AnythingOfType("database.UserDTO")).Return(tt.mockUpdateError)
			}
			defer mockService.AssertExpectations(t)

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			handler := updateUser.New(logger, mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("PUT /users/{userID}", handler)

			server := httptest.NewServer(mux)
			defer server.Close()

			url := server.URL + "/users/" + tt.userID

			requestBody, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			req = req.WithContext(context.WithValue(req.Context(), "user_id", "123"))

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			var responseBody updateUser.Response
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, responseBody.Status)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
