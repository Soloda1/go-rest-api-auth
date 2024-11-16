package createUser_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-rest-api-auth/internal/database"
	"go-rest-api-auth/internal/handlers/user/createUser"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCreateUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    createUser.Request
		mockResponse   database.UserDTO
		mockError      error
		expectedStatus string
		expectedBody   createUser.Response
	}{
		{
			name: "SuccessfulCreateUser",
			requestBody: createUser.Request{
				Username:    "testuser",
				Password:    "testpass",
				Description: "test description",
			},
			mockResponse: database.UserDTO{
				Id:          1,
				Username:    "testuser",
				Password:    "testpass",
				Description: "test description",
				DateJoined:  pgtype.Date{},
			},
			mockError:      nil,
			expectedStatus: "OK",
			expectedBody: createUser.Response{
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
			name: "InvalidRequestBody",
			requestBody: createUser.Request{
				Username: "", // Невалидные данные
				Password: "testpass",
			},
			mockResponse:   database.UserDTO{},
			mockError:      nil,
			expectedStatus: "Bad Request",
			expectedBody: createUser.Response{
				Status: "Bad Request",
				Error:  "failed to validate request",
			},
		},
		{
			name: "CreateUserError",
			requestBody: createUser.Request{
				Username:    "testuser",
				Password:    "testpass",
				Description: "test description",
			},
			mockResponse:   database.UserDTO{},
			mockError:      errors.New("create user failed"),
			expectedStatus: "Bad Request",
			expectedBody: createUser.Response{
				Status: "Bad Request",
				Error:  "failed to create user",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.UserService)
			if tt.name != "InvalidRequestBody" {
				mockService.On("CreateUser", mock.AnythingOfType("database.UserDTO")).Return(tt.mockResponse, tt.mockError)
			}
			defer mockService.AssertExpectations(t)

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			handler := createUser.New(logger, mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /users", handler)

			server := httptest.NewServer(mux)
			defer server.Close()

			url := server.URL + "/users"

			requestBody, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			req = req.WithContext(context.WithValue(req.Context(), "user_id", "123"))

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			var responseBody createUser.Response
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, responseBody.Status)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
