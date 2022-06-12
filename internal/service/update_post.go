package service

import (
	"github.com/ricardojonathanromero/api-protobuf/internal/domain/models"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *service) UpdatePost(id primitive.ObjectID, in *sma.UpdatePost) error {
	log.Infof("updating post %s", id.Hex())

	document := &models.Post{
		Title:       in.Title,
		Description: in.Description,
		MediaIDs:    in.MediaIds,
	}

	log.Info("updating document")
	err := s.repo.UpdateDocumentByID(id, document)
	if err != nil {
		log.Errorf("error updating document. reason - %v", err)
		return err
	}

	log.Info("document updated!")
	return nil
}
