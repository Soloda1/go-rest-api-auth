package getAllPosts_test

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"go-rest-api-auth/internal/database"
	"go-rest-api-auth/internal/handlers/post/getAllPosts"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetAllPostsHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   []database.PostDTO
		mockError      error
		expectedStatus string
		expectedBody   getAllPosts.Response
	}{
		{
			name: "SuccessfulGetAllPosts",
			mockResponse: []database.PostDTO{
				{
					Id:      1,
					Title:   "Test Title 1",
					Content: "Test Content 1",
					UserId:  123,
					Tags:    []string{"tag1", "tag2"},
				},
				{
					Id:      2,
					Title:   "Test Title 2",
					Content: "Test Content 2",
					UserId:  124,
					Tags:    []string{"tag3", "tag4"},
				},
			},
			mockError:      nil,
			expectedStatus: "OK",
			expectedBody: getAllPosts.Response{
				Status: "OK",
				Posts: []database.PostDTO{
					{
						Id:      1,
						Title:   "Test Title 1",
						Content: "Test Content 1",
						UserId:  123,
						Tags:    []string{"tag1", "tag2"},
					},
					{
						Id:      2,
						Title:   "Test Title 2",
						Content: "Test Content 2",
						UserId:  124,
						Tags:    []string{"tag3", "tag4"},
					},
				},
			},
		},
		{
			name:           "ErrorGetAllPosts",
			mockResponse:   nil,
			mockError:      errors.New("get all posts failed"),
			expectedStatus: "Bad Request",
			expectedBody: getAllPosts.Response{
				Status: "Bad Request",
				Error:  "get all posts failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.PostService)
			mockService.On("GetALlPosts").Return(tt.mockResponse, tt.mockError)
			defer mockService.AssertExpectations(t)

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			handler := getAllPosts.New(logger, mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /posts", handler)

			server := httptest.NewServer(mux)
			defer server.Close()

			url := server.URL + "/posts"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			var responseBody getAllPosts.Response
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
