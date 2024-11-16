package getPost_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-rest-api-auth/internal/database"
	"go-rest-api-auth/internal/handlers/post/getPost"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetPostHandler(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockResponse   database.PostDTO
		mockError      error
		expectedStatus string
		expectedBody   getPost.Response
	}{
		{
			name:   "SuccessfulGetPost",
			postID: "1",
			mockResponse: database.PostDTO{
				Id:      1,
				Title:   "Test Title",
				Content: "Test Content",
				UserId:  123,
				Tags:    []string{"tag1", "tag2"},
			},
			mockError:      nil,
			expectedStatus: "OK",
			expectedBody: getPost.Response{
				Status: "OK",
				Post: database.PostDTO{
					Id:      1,
					Title:   "Test Title",
					Content: "Test Content",
					UserId:  123,
					Tags:    []string{"tag1", "tag2"},
				},
			},
		},
		{
			name:           "InvalidPostID",
			postID:         "abc", // Невалидный ID
			mockResponse:   database.PostDTO{},
			mockError:      nil,
			expectedStatus: "Bad Request",
			expectedBody: getPost.Response{
				Status: "Bad Request",
				Error:  "Invalid post id",
			},
		},
		{
			name:           "PostNotFound",
			postID:         "2",
			mockResponse:   database.PostDTO{},
			mockError:      errors.New("post not found"),
			expectedStatus: "Bad Request",
			expectedBody: getPost.Response{
				Status: "Bad Request",
				Error:  "post not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.PostService)
			if tt.name != "InvalidPostID" {
				mockService.On("GetPost", mock.AnythingOfType("int")).Return(tt.mockResponse, tt.mockError)
			}
			defer mockService.AssertExpectations(t)

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			handler := getPost.New(logger, mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /posts/{postID}", handler)

			server := httptest.NewServer(mux)
			defer server.Close()

			url := server.URL + "/posts/" + tt.postID

			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			req = req.WithContext(context.WithValue(req.Context(), "user_id", "123"))

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			var responseBody getPost.Response
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, responseBody.Status)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
