package service

import (
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *service) RetrievePost(id primitive.ObjectID) (*sma.Post, error) {
	var res *sma.Post

	log.Infof("looking for %s", id.Hex())

	doc, err := s.repo.FindDocumentByID(id)
	if err != nil {
		log.Errorf("document not found! reason\n%v", err)
		return res, err
	}

	res = &sma.Post{
		Id:          doc.ID.Hex(),
		Title:       doc.Title,
		Description: doc.Description,
		UserId:      doc.UserID,
		Status:      statusStrToInt(doc.Status),
		CreatedAt:   timestamppb.New(*doc.CreatedAt),
		UpdatedAt:   timestamppb.New(*doc.UpdatedAt),
	}

	return res, nil
}
