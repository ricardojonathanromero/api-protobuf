package handler

import (
	"errors"
	"fmt"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestHandler_DeletePost(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected *expectDef
		id       primitive.ObjectID
	}{
		{
			name:     "happy path",
			expected: NewExpected(&sma.PostDeleteResp{Message: "post id %s removed"}),
			id:       primitive.NewObjectID(),
		},
		{
			name:     "post id required error message",
			expected: NewExpected(errors.New("'post_id' field is required")),
		},
		{
			name:     "error deleting post",
			err:      errors.New("error deleting post"),
			expected: NewExpected(errors.New("error deleting post")),
			id:       primitive.NewObjectID(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.expected.Clear()

			payload := &sma.PostIdReq{PostId: ""}
			if test.id != primitive.NilObjectID {
				srv.On("RemovePost", test.id).Return(test.err)
				payload.PostId = test.id.Hex()
			}

			if test.id == primitive.NilObjectID || test.err != nil {
				test.expected.SetValue(nil)
			}

			r, e := constructor.DeletePost(ctx, payload)
			if test.expected.IsError() {
				eval(t, assert.Error(t, e))
				eval(t, assert.EqualError(t, e, test.expected.GetValue().(error).Error()))
			} else {
				m := test.expected.GetValue().(*sma.PostDeleteResp)
				m.Message = fmt.Sprintf(m.Message, payload.PostId)
				eval(t, assert.Equal(t, test.expected.GetValue(), r))
			}
		})
	}
}
