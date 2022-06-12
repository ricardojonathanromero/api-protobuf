package service

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *service) RemovePost(id primitive.ObjectID) error {
	log.Info("removing document...")

	err := s.repo.DeleteDocumentByID(id)
	if err != nil {
		log.Errorf("document not deleted. reason \n%v", err)
		return err
	}

	log.Info("removed!")
	return nil
}
