package getAllUsers_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go-rest-api-auth/internal/database"
	"go-rest-api-auth/internal/handlers/user/getAllUsers"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetAllUsersHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   []database.UserDTO
		mockError      error
		expectedStatus string
		expectedBody   getAllUsers.Response
	}{
		{
			name: "SuccessfulGetAllUsers",
			mockResponse: []database.UserDTO{
				{
					Id:          1,
					Username:    "testuser1",
					Password:    "testpass1",
					Description: "test description 1",
					DateJoined:  pgtype.Date{},
				},
				{
					Id:          2,
					Username:    "testuser2",
					Password:    "testpass2",
					Description: "test description 2",
					DateJoined:  pgtype.Date{},
				},
			},
			mockError:      nil,
			expectedStatus: "OK",
			expectedBody: getAllUsers.Response{
				Status: "OK",
				Users: []database.UserDTO{
					{
						Id:          1,
						Username:    "testuser1",
						Password:    "testpass1",
						Description: "test description 1",
						DateJoined:  pgtype.Date{},
					},
					{
						Id:          2,
						Username:    "testuser2",
						Password:    "testpass2",
						Description: "test description 2",
						DateJoined:  pgtype.Date{},
					},
				},
			},
		},
		{
			name:           "GetAllUsersError",
			mockResponse:   nil,
			mockError:      errors.New("get all users failed"),
			expectedStatus: "Bad Request",
			expectedBody: getAllUsers.Response{
				Status: "Bad Request",
				Error:  "get all users failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.UserService)
			mockService.On("GetALlUsers").Return(tt.mockResponse, tt.mockError)
			defer mockService.AssertExpectations(t)

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			handler := getAllUsers.New(logger, mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /users", handler)

			server := httptest.NewServer(mux)
			defer server.Close()

			url := server.URL + "/users"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			req = req.WithContext(context.WithValue(req.Context(), "user_id", "123"))

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			var responseBody getAllUsers.Response
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, responseBody.Status)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
