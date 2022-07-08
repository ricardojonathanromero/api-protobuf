package service

import (
	"github.com/ricardojonathanromero/api-protobuf/internal/domain/models"
	"github.com/ricardojonathanromero/api-protobuf/internal/port"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	repo port.IRepository
	ut   port.IUtils
}

var _ port.IService = (*service)(nil)

func modelToProto(posts []*models.Post) []*sma.Post {
	result := make([]*sma.Post, len(posts))

	for i, post := range posts {
		result[i] = &sma.Post{
			Id:          post.ID.Hex(),
			Title:       post.Title,
			Description: post.Description,
			UserId:      post.UserID,
			Status:      statusStrToInt(post.Status),
			CreatedAt:   timestamppb.New(*post.CreatedAt),
			UpdatedAt:   timestamppb.New(*post.UpdatedAt),
		}
	}

	return result
}

func statusStrToInt(status string) sma.PostStatus {
	val, ok := sma.PostStatus_value[status]
	if !ok {
		return sma.PostStatus_POST_STATUS_UNSPECIFIED
	}

	return sma.PostStatus(val)
}

func calculateTotalPages(total, perPage int64) uint64 {
	if perPage > total {
		return 1
	}

	pages := total / perPage
	// mod per page with total
	residue := total % perPage

	if residue > 0 {
		pages++
	}

	return uint64(pages)
}

func New(repo port.IRepository, ut port.IUtils) port.IService {
	return &service{repo: repo, ut: ut}
}
