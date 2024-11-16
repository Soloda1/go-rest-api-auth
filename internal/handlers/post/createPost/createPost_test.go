package createPost_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-rest-api-auth/internal/database"
	createPost "go-rest-api-auth/internal/handlers/post/createPost"
	"go-rest-api-auth/testing/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestCreatePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    createPost.Request
		mockResponse   database.PostDTO
		mockError      error
		expectedStatus string
		expectedBody   createPost.Response
	}{
		{
			name:   "SuccessfulPostCreation",
			userID: "123",
			requestBody: createPost.Request{
				Title:   "Test Title",
				Content: "Test Content",
				Tags:    []string{"tag1", "tag2"},
			},
			mockResponse: database.PostDTO{
				Id:      1,
				Title:   "Test Title",
				Content: "Test Content",
				UserId:  123,
				Tags:    []string{"tag1", "tag2"},
				CreatedAt: pgtype.Timestamp{
					Time:  time.Date(2024, 10, 10, 12, 0, 0, 0, time.UTC),
					Valid: true,
				},
			},
			mockError:      nil,
			expectedStatus: "OK",
			expectedBody: createPost.Response{
				Status: "OK",
				Post: database.PostDTO{
					Id:      1, // JSON decodes numbers as float64
					Title:   "Test Title",
					Content: "Test Content",
					UserId:  123,
					CreatedAt: pgtype.Timestamp{
						Time:  time.Date(2024, 10, 10, 12, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Tags: []string{"tag1", "tag2"},
				},
			},
		},
		{
			name:           "InvalidRequestBody",
			userID:         "123",
			requestBody:    createPost.Request{}, // Empty request
			mockResponse:   database.PostDTO{},
			mockError:      nil,
			expectedStatus: "Bad Request",
			expectedBody: createPost.Response{
				Status: "Bad Request",
				Error:  "failed to validate request",
			},
		},
		{
			name:   "PostCreationFailure",
			userID: "123",
			requestBody: createPost.Request{
				Title:   "Test Title",
				Content: "Test Content",
				Tags:    []string{"tag1", "tag2"},
			},
			mockResponse:   database.PostDTO{},
			mockError:      errors.New("database error"),
			expectedStatus: "Bad Request",
			expectedBody: createPost.Response{
				Status: "Bad Request",
				Error:  "failed to create post",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.PostService)
			if tt.mockError == nil && tt.name != "InvalidRequestBody" {
				mockService.On("CreatePost", mock.Anything).Return(tt.mockResponse, nil)
			} else if tt.mockError != nil {
				mockService.On("CreatePost", mock.Anything).Return(database.PostDTO{}, tt.mockError)
			}
			defer mockService.AssertExpectations(t)

			// Создаем запрос
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/create_post", bytes.NewReader(body))
			req = req.WithContext(context.WithValue(req.Context(), "user_id", tt.userID))
			w := httptest.NewRecorder()

			// Логгер
			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			// Хендлер
			handler := createPost.New(logger, mockService)
			handler(w, req)

			// Проверка ответа
			resp := w.Result()
			defer resp.Body.Close()
			var responseBody createPost.Response
			err := json.NewDecoder(resp.Body).Decode(&responseBody)

			assert.Equal(t, tt.expectedStatus, responseBody.Status)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
