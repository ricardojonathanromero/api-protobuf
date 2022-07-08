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

func TestService_ListPosts(t *testing.T) {
	id := primitive.NewObjectID()
	now := time.Now()
	tests := []struct {
		title         string
		input         *sma.ListPostsReq
		countRes      int64
		countErr      error
		findRes       []*models.Post
		findErr       error
		isErrExpected bool
		expected      interface{}
	}{
		{
			title: "happy_path",
			input: &sma.ListPostsReq{
				UserId:  "1",
				Page:    1,
				PerPage: 10,
			},
			countRes: 1,
			findRes: []*models.Post{{
				ID:          id,
				Title:       "success",
				Description: "success",
				UserID:      "1",
				MediaIDs:    []string{"hello"},
				Status:      sma.PostStatus_POST_STATUS_ACTIVE.String(),
				CreatedAt:   &now,
				UpdatedAt:   &now,
			}},
			expected: &sma.ListPostsResp{
				Posts: []*sma.Post{
					{
						Id:          id.Hex(),
						Title:       "success",
						Description: "success",
						UserId:      "1",
						Status:      sma.PostStatus_POST_STATUS_ACTIVE,
						CreatedAt:   timestamppb.New(now),
						UpdatedAt:   timestamppb.New(now),
					},
				},
				PageInfo: &sma.PageInfo{
					Page:       1,
					PageSize:   10,
					TotalItems: 1,
					TotalPages: 1,
				},
			},
		},
		{
			title: "happy_path_unspecified",
			input: &sma.ListPostsReq{
				UserId:  "1",
				Page:    1,
				PerPage: 10,
			},
			countRes: 1,
			findRes: []*models.Post{{
				ID:          id,
				Title:       "success",
				Description: "success",
				UserID:      "1",
				MediaIDs:    []string{"hello"},
				Status:      "unspecified",
				CreatedAt:   &now,
				UpdatedAt:   &now,
			}},
			expected: &sma.ListPostsResp{
				Posts: []*sma.Post{
					{
						Id:          id.Hex(),
						Title:       "success",
						Description: "success",
						UserId:      "1",
						Status:      sma.PostStatus_POST_STATUS_UNSPECIFIED,
						CreatedAt:   timestamppb.New(now),
						UpdatedAt:   timestamppb.New(now),
					},
				},
				PageInfo: &sma.PageInfo{
					Page:       1,
					PageSize:   10,
					TotalItems: 1,
					TotalPages: 1,
				},
			},
		},
		{
			title: "happy_path_pagination",
			input: &sma.ListPostsReq{
				UserId:  "1",
				Page:    1,
				PerPage: 10,
			},
			countRes: 15,
			findRes: []*models.Post{{
				ID:          id,
				Title:       "success",
				Description: "success",
				UserID:      "1",
				MediaIDs:    []string{"hello"},
				Status:      sma.PostStatus_POST_STATUS_ACTIVE.String(),
				CreatedAt:   &now,
				UpdatedAt:   &now,
			}},
			expected: &sma.ListPostsResp{
				Posts: []*sma.Post{
					{
						Id:          id.Hex(),
						Title:       "success",
						Description: "success",
						UserId:      "1",
						Status:      sma.PostStatus_POST_STATUS_ACTIVE,
						CreatedAt:   timestamppb.New(now),
						UpdatedAt:   timestamppb.New(now),
					},
				},
				PageInfo: &sma.PageInfo{
					Page:       1,
					PageSize:   10,
					TotalItems: 15,
					TotalPages: 2,
				},
			},
		},
		{
			title:         "error_count_documents",
			input:         &sma.ListPostsReq{UserId: "0"},
			countRes:      0,
			countErr:      errors.New("mongo: error documents"),
			findRes:       []*models.Post{},
			isErrExpected: true,
			expected:      errors.New("mongo: error documents"),
		},
		{
			title:         "error_count_documents",
			input:         &sma.ListPostsReq{UserId: "0"},
			countRes:      1,
			findRes:       []*models.Post{},
			findErr:       errors.New("mongo: no documents"),
			isErrExpected: true,
			expected:      errors.New("mongo: no documents"),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			mockRepo := new(MockRepo)
			mockUtils := new(MockUtils)

			mockRepo.On("CountDocuments", test.input).Return(test.countRes, test.countErr)
			mockRepo.On("FindDocuments", test.input).Return(test.findRes, test.findErr)

			srv := New(mockRepo, mockUtils)

			posts, err := srv.ListPosts(test.input)
			if test.isErrExpected {
				eval(t, assert.Error(t, err))
				eval(t, assert.Equal(t, err, test.expected.(error)))
			} else {
				eval(t, assert.NoError(t, err))
				eval(t, assert.Equal(t, test.expected.(*sma.ListPostsResp), posts))
			}
		})
	}
}
