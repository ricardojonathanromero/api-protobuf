package handler

import (
	"context"
	"errors"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *handler) UpdatePost(_ context.Context, req *sma.UpdatePostReq) (*sma.Post, error) {
	var res *sma.Post
	log.Info("start handler updatePost()")

	id, err := primitive.ObjectIDFromHex(req.PostId)
	if err != nil {
		log.Errorf("invalid id %s", req.PostId)
		return res, errors.New("'post_id' is not valid")
	}

	err = h.srv.UpdatePost(id, req.Post)
	if err != nil {
		log.Error("error updating post")
		return res, err
	}

	log.Info("looking for new document updated")
	post, err := h.srv.RetrievePost(id)
	if err != nil {
		log.Errorf("id %s not founded", id)
		return res, err
	}

	return post, nil
}
