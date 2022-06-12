package handler

import (
	"context"
	"errors"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	log "github.com/sirupsen/logrus"
)

func (h *handler) ListPosts(_ context.Context, req *sma.ListPostsReq) (*sma.ListPostsResp, error) {
	var res *sma.ListPostsResp
	log.Info("listing posts")

	if len(req.UserId) <= 0 {
		log.Errorf("user id %s is empty", req.UserId)
		return res, errors.New("'user id' field is required")
	}

	if req.PerPage == 0 {
		req.PerPage = 20
	}

	if req.Page == 0 {
		req.Page = 1
	}

	posts, err := h.srv.ListPosts(req)
	if err != nil {
		log.Errorf("error produces from ListPost(). reason \n%v", err)
		return res, err
	}

	return posts, nil
}
