package handler

import (
	"context"
	"errors"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	log "github.com/sirupsen/logrus"
)

func (h *handler) CreatePost(_ context.Context, req *sma.CreatePostReq) (*sma.Post, error) {
	log.Infof("request CreatePost has been received")

	log.Info("validating request")
	if len(req.UserId) == 0 {
		return nil, errors.New("'user_id' field is required")
	}

	if len(req.Title) == 0 {
		return nil, errors.New("'title' field is required")
	}

	if len(req.Description) == 0 {
		return nil, errors.New("'description' filed is required")
	}

	log.Info("all works fine, next step will be executed")
	post, err := h.srv.CreatePost(req)
	if err != nil {
		return nil, err
	}

	log.Info("post created!")
	return post, nil
}
