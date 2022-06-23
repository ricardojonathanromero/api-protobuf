package handler

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestHandler_UpdatePost(t *testing.T) {
	id := primitive.NewObjectID().Hex()
	uID := uuid.New().String()
	now := timestamppb.Now()
	tests := []struct {
		name              string
		req               *sma.UpdatePostReq
		isUpdateOn        bool
		isRetrievedPostOn bool
		updatePostErr     error
		retrieveErr       error
		retrieveRes       *sma.Post
		expected          *expectDef
	}{
		{
			name: "happy path",
			req: &sma.UpdatePostReq{PostId: id, Post: &sma.UpdatePost{
				Title:       "ExampleUpdated",
				Description: "updated",
				MediaIds:    []string{"one", "two", "three"},
				ScheduledAt: now,
			}},
			isUpdateOn:        true,
			isRetrievedPostOn: true,
			retrieveRes: &sma.Post{
				Id:          id,
				Title:       "ExampleUpdated",
				Description: "updated",
				UserId:      uID,
				Status:      sma.PostStatus_POST_STATUS_ACTIVE,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			expected: NewExpected(&sma.Post{
				Id:          id,
				Title:       "ExampleUpdated",
				Description: "updated",
				UserId:      uID,
				Status:      sma.PostStatus_POST_STATUS_ACTIVE,
				CreatedAt:   now,
				UpdatedAt:   now,
			}),
		},
		{
			name:     "invalid post id",
			req:      &sma.UpdatePostReq{},
			expected: NewExpected(errors.New("'post_id' is not valid")),
		},
		{
			name:          "UpdatedPost error",
			req:           &sma.UpdatePostReq{PostId: id, Post: &sma.UpdatePost{Title: "errorUpd"}},
			isUpdateOn:    true,
			updatePostErr: errors.New("error updating post"),
			expected:      NewExpected(errors.New("error updating post")),
		},
		{
			name: "RetrievePost error",
			req: &sma.UpdatePostReq{PostId: primitive.NewObjectID().Hex(), Post: &sma.UpdatePost{
				Title:       "ExampleUpdated",
				Description: "updated",
				MediaIds:    []string{"one", "two", "three"},
				ScheduledAt: now,
			}},
			isUpdateOn:        true,
			isRetrievedPostOn: true,
			retrieveErr:       errors.New("post not found"),
			expected:          NewExpected(errors.New("post not found")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.expected.Clear()

			test.expected.SetValue(test.retrieveRes)
			if test.isUpdateOn {
				testID, _ := primitive.ObjectIDFromHex(test.req.PostId)
				srv.On("UpdatePost", testID, test.req.Post).Return(test.updatePostErr)
			}
			if test.isRetrievedPostOn {
				testID, _ := primitive.ObjectIDFromHex(test.req.PostId)
				srv.On("RetrievePost", testID).Return(test.retrieveRes, test.retrieveErr)
			}

			r, e := constructor.UpdatePost(ctx, test.req)

			if test.expected.IsError() {
				eval(t, assert.Error(t, e))
				eval(t, assert.EqualError(t, e, test.expected.GetValue().(error).Error()))
			} else {
				eval(t, assert.Equal(t, test.expected.GetValue(), r))
			}
		})
	}
}
