package handler

import (
	"context"
	"errors"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *handler) ShowPost(_ context.Context, req *sma.PostIdReq) (*sma.Post, error) {
	var res *sma.Post
	log.Info("start ShowPost")

	id, err := primitive.ObjectIDFromHex(req.PostId)
	if err != nil {
		log.Errorf("error invalid post if %s", req.PostId)
		return res, errors.New("'post_id' is not valid")
	}

	res, err = h.srv.RetrievePost(id)
	if err != nil {
		log.Errorf("erro finding post. reason \n%v", err)
		return res, err
	}

	return res, nil
}
