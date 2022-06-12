package service

import (
	"github.com/ricardojonathanromero/api-protobuf/internal/domain/models"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// CreatePost generates new post document to insert into db
func (s *service) CreatePost(in *sma.CreatePostReq) (*sma.Post, error) {
	log.Info("start createPost service")

	now := time.Now()

	log.Info("activating post")
	document := &models.Post{
		Title:       in.Title,
		Description: in.Description,
		UserID:      in.UserId,
		MediaIDs:    in.MediaIds,
		Status:      sma.PostStatus_POST_STATUS_ACTIVE.String(),
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	id, err := s.repo.InsertDocument(document)
	if err != nil {
		log.Errorf("error inserting document in db. reason\n%v", err)
		return nil, err
	}

	log.Info("document generated!")
	result := &sma.Post{
		Id:          id.Hex(),
		Title:       document.Title,
		Description: document.Description,
		UserId:      document.UserID,
		Status:      sma.PostStatus_POST_STATUS_ACTIVE,
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
	}
	return result, nil
}
