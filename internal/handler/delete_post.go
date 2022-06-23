package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *handler) DeletePost(_ context.Context, req *sma.PostIdReq) (*sma.PostDeleteResp, error) {
	var res *sma.PostDeleteResp
	log.Infof("deleting post %s", req.PostId)

	id, err := primitive.ObjectIDFromHex(req.PostId)
	if err != nil {
		log.Errorf("id %s is not valid. reason\n%v", req.PostId, err)
		return res, errors.New("'post_id' field is required")
	}

	err = h.srv.RemovePost(id)
	if err != nil {
		log.Errorf("error returned from RemovePost(). reason\n%v", err)
		return res, err
	}

	log.Info("post deleted")
	res = &sma.PostDeleteResp{Message: fmt.Sprintf("post id %s removed", req.PostId)}
	return res, nil
}
