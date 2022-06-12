package service

import (
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	log "github.com/sirupsen/logrus"
)

func (s *service) ListPosts(in *sma.ListPostsReq) (*sma.ListPostsResp, error) {
	var result *sma.ListPostsResp

	log.Info("finding posts")
	res, err := s.repo.FindDocuments(in)
	if err != nil {
		log.Errorf("error finding documents. reason \n%v", err)
		return result, err
	}

	log.Info("calculating pagination")
	result = &sma.ListPostsResp{
		Posts: modelToProto(res.Posts),
		PageInfo: &sma.PageInfo{
			Page:       uint64(in.Page),
			PageSize:   uint64(in.PerPage),
			TotalItems: uint64(res.Count),
			TotalPages: calculateTotalPages(res.Count, in.PerPage),
		},
	}

	return result, nil
}
