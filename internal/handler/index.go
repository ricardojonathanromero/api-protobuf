package handler

import (
	"github.com/ricardojonathanromero/api-protobuf/internal/port"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
)

type handler struct {
	sma.UnimplementedPostsServer
	srv port.IService
}

func New(srv port.IService) *handler {
	return &handler{srv: srv}
}
