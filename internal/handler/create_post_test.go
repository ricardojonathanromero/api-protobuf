package handler

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestHandler_CreatePost(t *testing.T) {
	tests := []struct {
		name     string
		payload  *sma.CreatePostReq
		returned *protoDef
		err      error
		expected *expectDef
	}{
		{
			name: "happy path",
			payload: &sma.CreatePostReq{
				Title:       "ExampleOneOK",
				Description: "Happy path correct flow",
				UserId:      uuid.New().String(),
				MediaIds:    []string{"one"},
				ScheduledAt: timestamppb.New(time.Now()),
			},
			returned: NewProtoDef(false),
			expected: NewExpected(nil),
		},
		{
			name:     "user id required error message",
			payload:  &sma.CreatePostReq{},
			returned: NewProtoDef(true),
			expected: NewExpected(errors.New("'user_id' field is required")),
		},
		{
			name:     "title required error message",
			payload:  &sma.CreatePostReq{UserId: uuid.New().String()},
			returned: NewProtoDef(true),
			expected: NewExpected(errors.New("'title' field is required")),
		},
		{
			name:     "description required error message",
			payload:  &sma.CreatePostReq{UserId: uuid.New().String(), Title: "error description field"},
			returned: NewProtoDef(true),
			expected: NewExpected(errors.New("'description' field is required")),
		},
		{
			name: "error creating post",
			payload: &sma.CreatePostReq{
				Title:       "ErrorCreatingPost",
				Description: "service post create",
				UserId:      uuid.New().String(),
				MediaIds:    []string{"one", "two"},
				ScheduledAt: timestamppb.New(time.Now()),
			},
			returned: NewProtoDef(true),
			err:      errors.New("error creating post"),
			expected: NewExpected(errors.New("error creating post")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.returned.Clear()
			defer test.expected.Clear()
			test.returned.SetValue()

			returned := test.returned.GetValue()
			test.expected.SetValue(returned)

			srv.On("CreatePost", test.payload).Return(returned, test.err)
			r, e := constructor.CreatePost(ctx, test.payload)
			if test.expected.IsError() {
				eval(t, assert.Error(t, e))
				eval(t, assert.EqualError(t, e, test.expected.GetValue().(error).Error()))
			} else {
				eval(t, assert.Equal(t, test.expected.GetValue(), r))
			}
		})
	}
}
