package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestService_RemovePost(t *testing.T) {
	tests := []struct {
		title         string
		input         primitive.ObjectID
		deleteErr     error
		isErrExpected bool
		expected      interface{}
	}{
		{
			title: "happy_path",
			input: primitive.NewObjectID(),
		},
		{
			title:         "error_remove_document",
			deleteErr:     errors.New("mongo: error remove document"),
			isErrExpected: true,
			expected:      errors.New("mongo: error remove document"),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			mockRepo := new(MockRepo)
			mockUtils := new(MockUtils)

			mockRepo.On("DeleteDocumentByID", test.input).Return(test.deleteErr)

			srv := New(mockRepo, mockUtils)

			err := srv.RemovePost(test.input)
			if test.isErrExpected {
				eval(t, assert.Error(t, err))
				eval(t, assert.Equal(t, err, test.expected.(error)))
			} else {
				eval(t, assert.NoError(t, err))
			}
		})
	}
}
