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

func TestHandler_ShowPost(t *testing.T) {
	id := primitive.NewObjectID().Hex()
	now := timestamppb.Now()
	tests := []struct {
		name              string
		req               *sma.PostIdReq
		isRetrievedPostOn bool
		retrieveErr       error
		retrieveRes       *sma.Post
		expected          *expectDef
	}{
		{
			name:              "happy path ShowPost",
			req:               &sma.PostIdReq{PostId: id},
			isRetrievedPostOn: true,
			retrieveRes: &sma.Post{
				Id:          id,
				Title:       "ExamplePost",
				Description: "This is an example",
				UserId:      uuid.New().String(),
				Status:      sma.PostStatus_POST_STATUS_ACTIVE,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			expected: NewExpected(&sma.Post{
				Id:          id,
				Title:       "ExamplePost",
				Description: "This is an example",
				UserId:      uuid.New().String(),
				Status:      sma.PostStatus_POST_STATUS_ACTIVE,
				CreatedAt:   now,
				UpdatedAt:   now,
			}),
		},
		{
			name:     "invalid post id",
			req:      &sma.PostIdReq{},
			expected: NewExpected(errors.New("'post_id' is not valid")),
		},
		{
			name:              "RetrievePost error",
			req:               &sma.PostIdReq{PostId: primitive.NewObjectID().Hex()},
			isRetrievedPostOn: true,
			retrieveErr:       errors.New("post not found"),
			expected:          NewExpected(errors.New("post not found")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.expected.Clear()

			test.expected.SetValue(test.retrieveRes)
			if test.isRetrievedPostOn {
				testID, _ := primitive.ObjectIDFromHex(test.req.PostId)
				srv.On("RetrievePost", testID).Return(test.retrieveRes, test.retrieveErr)
			}

			r, e := constructor.ShowPost(ctx, test.req)

			if test.expected.IsError() {
				eval(t, assert.Error(t, e))
				eval(t, assert.EqualError(t, e, test.expected.GetValue().(error).Error()))
			} else {
				eval(t, assert.Equal(t, test.expected.GetValue(), r))
			}
		})
	}
}
