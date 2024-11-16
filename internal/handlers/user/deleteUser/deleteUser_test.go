package deleteUser_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-rest-api-auth/internal/database"
	"go-rest-api-auth/internal/handlers/user/deleteUser"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDeleteUserHandler(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		mockGetResponse database.UserDTO
		mockGetError    error
		mockDeleteError error
		expectedStatus  string
		expectedBody    deleteUser.Response
	}{
		{
			name:   "SuccessfulDeleteUser",
			userID: "1",
			mockGetResponse: database.UserDTO{
				Id:          1,
				Username:    "testuser",
				Password:    "testpass",
				Description: "test description",
				DateJoined:  pgtype.Date{},
			},
			mockGetError:    nil,
			mockDeleteError: nil,
			expectedStatus:  "OK",
			expectedBody: deleteUser.Response{
				Status: "OK",
				UserID: 1,
			},
		},
		{
			name:            "InvalidUserID",
			userID:          "abc", // Невалидный ID
			mockGetResponse: database.UserDTO{},
			mockGetError:    nil,
			mockDeleteError: nil,
			expectedStatus:  "Bad Request",
			expectedBody: deleteUser.Response{
				Status: "Bad Request",
				Error:  "Invalid user id",
			},
		},
		{
			name:            "UserNotFound",
			userID:          "2",
			mockGetResponse: database.UserDTO{},
			mockGetError:    errors.New("user not found"),
			mockDeleteError: nil,
			expectedStatus:  "Bad Request",
			expectedBody: deleteUser.Response{
				Status: "Bad Request",
				Error:  "User not found",
			},
		},
		{
			name:   "ErrorDeletingUser",
			userID: "1",
			mockGetResponse: database.UserDTO{
				Id:          1,
				Username:    "testuser",
				Password:    "testpass",
				Description: "test description",
				DateJoined:  pgtype.Date{},
			},
			mockGetError:    nil,
			mockDeleteError: errors.New("delete error"),
			expectedStatus:  "Bad Request",
			expectedBody: deleteUser.Response{
				Status: "Bad Request",
				Error:  "Error deleting user",
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
				mockService.On("DeleteUser", mock.AnythingOfType("int")).Return(tt.mockDeleteError)
			}
			defer mockService.AssertExpectations(t)

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			handler := deleteUser.New(logger, mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /users/{userID}", handler)

			server := httptest.NewServer(mux)
			defer server.Close()

			url := server.URL + "/users/" + tt.userID

			req, err := http.NewRequest(http.MethodDelete, url, nil)
			assert.NoError(t, err)

			req = req.WithContext(context.WithValue(req.Context(), "user_id", "123"))

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			var responseBody deleteUser.Response
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, responseBody.Status)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
