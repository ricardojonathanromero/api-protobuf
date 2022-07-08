package service

import (
	"errors"
	"github.com/ricardojonathanromero/api-protobuf/internal/domain/models"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestService_UpdatePost(t *testing.T) {
	id := primitive.NewObjectID()
	errID := primitive.NewObjectID()
	tests := []struct {
		title         string
		inputID       primitive.ObjectID
		inputUpd      *sma.UpdatePost
		updID         primitive.ObjectID
		updDoc        *models.Post
		updateErr     error
		isErrExpected bool
		expected      interface{}
	}{
		{
			title:   "happy_path",
			inputID: id,
			inputUpd: &sma.UpdatePost{
				Title:       "updated",
				Description: "updated",
				MediaIds:    []string{"upd"},
			},
			updID: id,
			updDoc: &models.Post{
				Title:       "updated",
				Description: "updated",
				MediaIDs:    []string{"upd"},
			},
		},
		{
			title:   "error_update_document",
			inputID: errID,
			inputUpd: &sma.UpdatePost{
				Title:       "error",
				Description: "error",
				MediaIds:    []string{"err"},
			},
			updID: errID,
			updDoc: &models.Post{
				Title:       "error",
				Description: "error",
				MediaIDs:    []string{"err"},
			},
			updateErr:     errors.New("mongo: document not updated"),
			isErrExpected: true,
			expected:      errors.New("mongo: document not updated"),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			now := time.Now()

			mockRepo := new(MockRepo)
			mockUtils := new(MockUtils)

			if test.updDoc != nil {
				test.updDoc.UpdatedAt = &now
			}

			mockUtils.On("Now").Return(&now)
			mockRepo.On("UpdateDocumentByID", test.updID, test.updDoc).Return(test.updateErr)

			srv := New(mockRepo, mockUtils)

			err := srv.UpdatePost(test.inputID, test.inputUpd)
			if test.isErrExpected {
				eval(t, assert.Error(t, err))
				eval(t, assert.Equal(t, err, test.expected.(error)))
			} else {
				eval(t, assert.NoError(t, err))
			}
		})
	}
}
