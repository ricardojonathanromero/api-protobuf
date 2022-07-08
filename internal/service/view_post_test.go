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

func TestService_RetrievePost(t *testing.T) {
	id := primitive.NewObjectID()
	errID := primitive.NewObjectID()
	now := time.Now()

	tests := []struct {
		title         string
		id            primitive.ObjectID
		findDoc       *models.Post
		findErr       error
		isErrExpected bool
		expected      interface{}
	}{
		{
			title: "happy_path",
			id:    id,
			findDoc: &models.Post{
				ID:          id,
				Title:       "example",
				Description: "example",
				UserID:      "1",
				MediaIDs:    []string{"happy"},
				Status:      sma.PostStatus_POST_STATUS_ACTIVE.String(),
				CreatedAt:   &now,
				UpdatedAt:   &now,
			},
			expected: &sma.Post{
				Id:          id.Hex(),
				Title:       "example",
				Description: "example",
				UserId:      "1",
				Status:      sma.PostStatus_POST_STATUS_ACTIVE,
				CreatedAt:   timestamppb.New(now),
				UpdatedAt:   timestamppb.New(now),
			},
		},
		{
			title:         "error_retrieve_document",
			id:            errID,
			findDoc:       &models.Post{},
			findErr:       errors.New("mongo: document not exists"),
			isErrExpected: true,
			expected:      errors.New("mongo: document not exists"),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			mockRepo := new(MockRepo)
			mockUtils := new(MockUtils)

			mockRepo.On("FindDocumentByID", test.id).Return(test.findDoc, test.findErr)

			srv := New(mockRepo, mockUtils)

			doc, err := srv.RetrievePost(test.id)
			if test.isErrExpected {
				eval(t, assert.Error(t, err))
				eval(t, assert.Equal(t, err, test.expected.(error)))
			} else {
				eval(t, assert.NoError(t, err))
				eval(t, assert.Equal(t, test.expected.(*sma.Post), doc))
			}
		})
	}
}
