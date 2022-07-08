package service

import (
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	log "github.com/sirupsen/logrus"
)

func (s *service) ListPosts(in *sma.ListPostsReq) (*sma.ListPostsResp, error) {
	var result *sma.ListPostsResp

	log.Info("counting documents")
	count, err := s.repo.CountDocuments(in)
	if err != nil {
		log.Errorf("error counting documents. reason \n%v", err)
		return result, err
	}

	log.Info("finding posts")
	res, err := s.repo.FindDocuments(in)
	if err != nil {
		log.Errorf("error finding documents. reason \n%v", err)
		return result, err
	}

	log.Info("calculating pagination")
	result = &sma.ListPostsResp{
		Posts: modelToProto(res),
		PageInfo: &sma.PageInfo{
			Page:       uint64(in.Page),
			PageSize:   uint64(in.PerPage),
			TotalItems: uint64(count),
			TotalPages: calculateTotalPages(count, in.PerPage),
		},
	}

	return result, nil
}
