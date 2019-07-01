package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/wq1019/short/internal/lib"
	"github.com/wq1019/short/internal/model"
)

func (s *Service) CreateDomain(c context.Context, domain *model.Domain) (err error) {

	urls := lib.GenerateShortUrl(domain.RedirectUrl)
	shortUrl := ""
	for _, url := range urls {
		if exist := s.dao.ShortUrlHasBeenUsed(c, url); !exist {
			shortUrl = url
			break
		}
	}
	if shortUrl == "" {
		err = errors.New("短链接生成失败")
	}
	err = s.dao.CacheShortUrl(c, shortUrl)
	if err != nil {
		err = errors.WithMessage(err, "保存到 cache 失败")
		return err
	}
	return s.dao.CreateDomain(c, domain)
}
