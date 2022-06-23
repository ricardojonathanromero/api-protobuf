package handler

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestHandler_ListPosts(t *testing.T) {
	tests := []struct {
		name     string
		req      *sma.ListPostsReq
		res      *sma.ListPostsResp
		err      error
		expected *expectDef
	}{
		{
			name: "happy path",
			req: &sma.ListPostsReq{
				UserId:  uuid.New().String(),
				Page:    1,
				PerPage: 10,
				Filter:  sma.Filters_FILTER_ACTIVE,
			},
			res: &sma.ListPostsResp{
				Posts: []*sma.Post{{
					Id:          "1",
					Title:       "Example",
					Description: "Example",
					UserId:      "user1",
					Status:      sma.PostStatus_POST_STATUS_ACTIVE,
					CreatedAt:   timestamppb.Now(),
					UpdatedAt:   timestamppb.Now(),
				}},
				PageInfo: &sma.PageInfo{
					Page:       1,
					PageSize:   10,
					TotalItems: 1,
					TotalPages: 1,
				},
			},
			expected: NewExpected(&sma.ListPostsResp{
				Posts: []*sma.Post{{
					Id:          "1",
					Title:       "Example",
					Description: "Example",
					UserId:      "user1",
					Status:      sma.PostStatus_POST_STATUS_ACTIVE,
					CreatedAt:   timestamppb.Now(),
					UpdatedAt:   timestamppb.Now(),
				}},
				PageInfo: &sma.PageInfo{
					Page:       1,
					PageSize:   10,
					TotalItems: 1,
					TotalPages: 1,
				},
			}),
		},
		{
			name: "happy path pagination",
			req: &sma.ListPostsReq{
				UserId: uuid.New().String(),
				Filter: sma.Filters_FILTER_ACTIVE,
			},
			res: &sma.ListPostsResp{
				Posts: []*sma.Post{{
					Id:          "1",
					Title:       "Example",
					Description: "Example",
					UserId:      "user1",
					Status:      sma.PostStatus_POST_STATUS_ACTIVE,
					CreatedAt:   timestamppb.Now(),
					UpdatedAt:   timestamppb.Now(),
				}},
				PageInfo: &sma.PageInfo{
					Page:       1,
					PageSize:   20,
					TotalItems: 1,
					TotalPages: 1,
				},
			},
			expected: NewExpected(&sma.ListPostsResp{
				Posts: []*sma.Post{{
					Id:          "1",
					Title:       "Example",
					Description: "Example",
					UserId:      "user1",
					Status:      sma.PostStatus_POST_STATUS_ACTIVE,
					CreatedAt:   timestamppb.Now(),
					UpdatedAt:   timestamppb.Now(),
				}},
				PageInfo: &sma.PageInfo{
					Page:       1,
					PageSize:   20,
					TotalItems: 1,
					TotalPages: 1,
				},
			}),
		},
		{
			name:     "user id empty",
			req:      &sma.ListPostsReq{},
			res:      &sma.ListPostsResp{},
			expected: NewExpected(errors.New("'user id' field is required")),
		},
		{
			name:     "list post error",
			req:      &sma.ListPostsReq{UserId: "2"},
			res:      &sma.ListPostsResp{},
			err:      errors.New("error retrieving list post"),
			expected: NewExpected(errors.New("error retrieving list post")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.expected.Clear()

			test.expected.SetValue(test.res)
			srv.On("ListPosts", test.req).Return(test.res, test.err)
			r, e := constructor.ListPosts(ctx, test.req)

			if test.expected.IsError() {
				eval(t, assert.Error(t, e))
				eval(t, assert.EqualError(t, e, test.expected.GetValue().(error).Error()))
			} else {
				eval(t, assert.Equal(t, test.expected.GetValue(), r))
			}
		})
	}
}
