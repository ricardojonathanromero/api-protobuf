package service

import (
	"errors"
	"github.com/ricardojonathanromero/api-protobuf/internal/domain/models"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestService_CreatePost(t *testing.T) {
	tests := []struct {
		title         string
		input         *sma.CreatePostReq
		doc           *models.Post
		now           *time.Time
		insertID      primitive.ObjectID
		insertErr     error
		isErrExpected bool
		expected      interface{}
	}{
		{
			title: "happy_path",
			input: &sma.CreatePostReq{
				Title:       "Example",
				Description: "example",
				UserId:      "1",
				MediaIds:    []string{"hello"},
				ScheduledAt: timestamppb.Now(),
			},
			doc: &models.Post{
				Title:       "Example",
				Description: "example",
				UserID:      "1",
				MediaIDs:    []string{"hello"},
				Status:      sma.PostStatus_POST_STATUS_ACTIVE.String(),
			},
			insertID: primitive.NewObjectID(),
		},
		{
			title: "error_insert_document",
			input: &sma.CreatePostReq{},
			doc: &models.Post{
				Status: sma.PostStatus_POST_STATUS_ACTIVE.String(),
			},
			insertID:      primitive.NilObjectID,
			insertErr:     errors.New("mongo: error insert document"),
			isErrExpected: true,
			expected:      errors.New("mongo: error insert document"),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			mockRepo := new(MockRepo)
			mockUtils := new(MockUtils)

			now := time.Now()

			mockRepo.On("InsertDocument", test.doc).Return(test.insertID, test.insertErr)
			mockUtils.On("Now").Return(&now)

			srv := New(mockRepo, mockUtils)

			if test.doc != nil {
				test.doc.CreatedAt = &now
				test.doc.UpdatedAt = &now
			}

			post, err := srv.CreatePost(test.input)
			if test.isErrExpected {
				eval(t, assert.Error(t, err))
				eval(t, assert.Equal(t, err, test.expected.(error)))
			} else {
				eval(t, assert.NoError(t, err))
				expected := &sma.Post{
					Id:          test.insertID.Hex(),
					Title:       test.doc.Title,
					Description: test.doc.Description,
					UserId:      test.doc.UserID,
					Status:      sma.PostStatus_POST_STATUS_ACTIVE,
					CreatedAt:   timestamppb.New(*test.doc.CreatedAt),
					UpdatedAt:   timestamppb.New(*test.doc.UpdatedAt),
				}
				eval(t, assert.Equal(t, expected, post))
			}
		})
	}
}
