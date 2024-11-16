package updatePost_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-rest-api-auth/internal/database"
	"go-rest-api-auth/internal/handlers/post/updatePost"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUpdatePostHandler(t *testing.T) {
	tests := []struct {
		name            string
		postID          string
		requestBody     updatePost.Request
		mockGetResponse database.PostDTO
		mockGetError    error
		mockUpdateError error
		expectedStatus  string
		expectedBody    updatePost.Response
	}{
		{
			name:   "SuccessfulUpdatePost",
			postID: "1",
			requestBody: updatePost.Request{
				Title:   "Updated Title",
				Content: "Updated Content",
				Tags:    []string{"tag1", "tag2"},
			},
			mockGetResponse: database.PostDTO{
				Id:      1,
				Title:   "Original Title",
				Content: "Original Content",
				UserId:  123,
				Tags:    []string{"tag3"},
			},
			mockGetError:    nil,
			mockUpdateError: nil,
			expectedStatus:  "OK",
			expectedBody: updatePost.Response{
				Status: "OK",
				Post: database.PostDTO{
					Id:      1,
					Title:   "Updated Title",
					Content: "Updated Content",
					UserId:  123,
					Tags:    []string{"tag1", "tag2"},
				},
			},
		},
		{
			name:            "InvalidPostID",
			postID:          "abc", // Невалидный ID
			requestBody:     updatePost.Request{},
			mockGetResponse: database.PostDTO{},
			mockGetError:    nil,
			mockUpdateError: nil,
			expectedStatus:  "Bad Request",
			expectedBody: updatePost.Response{
				Status: "Bad Request",
				Error:  "Invalid post id",
			},
		},
		{
			name:            "PostNotFound",
			postID:          "2",
			requestBody:     updatePost.Request{},
			mockGetResponse: database.PostDTO{},
			mockGetError:    errors.New("post not found"),
			mockUpdateError: nil,
			expectedStatus:  "Bad Request",
			expectedBody: updatePost.Response{
				Status: "Bad Request",
				Error:  "post not found",
			},
		},
		{
			name:   "ErrorUpdatingPost",
			postID: "1",
			requestBody: updatePost.Request{
				Title:   "Updated Title",
				Content: "Updated Content",
				Tags:    []string{"tag1", "tag2"},
			},
			mockGetResponse: database.PostDTO{
				Id:      1,
				Title:   "Original Title",
				Content: "Original Content",
				UserId:  123,
				Tags:    []string{"tag3"},
			},
			mockGetError:    nil,
			mockUpdateError: errors.New("update error"),
			expectedStatus:  "Bad Request",
			expectedBody: updatePost.Response{
				Status: "Bad Request",
				Error:  "failed to update post",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.PostService)
			if tt.name == "PostNotFound" {
				mockService.On("GetPost", mock.AnythingOfType("int")).Return(tt.mockGetResponse, tt.mockGetError)
			} else if tt.name != "InvalidPostID" {
				mockService.On("GetPost", mock.AnythingOfType("int")).Return(tt.mockGetResponse, tt.mockGetError)
				mockService.On("UpdatePost", mock.AnythingOfType("database.PostDTO")).Return(tt.mockUpdateError)
			}
			defer mockService.AssertExpectations(t)

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			handler := updatePost.New(logger, mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("PUT /posts/{postID}", handler)

			server := httptest.NewServer(mux)
			defer server.Close()

			url := server.URL + "/posts/" + tt.postID

			requestBody, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			req = req.WithContext(context.WithValue(req.Context(), "user_id", "123"))

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			var responseBody updatePost.Response
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, responseBody.Status)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
