package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/wq1019/short/internal/contants"
	"github.com/wq1019/short/internal/lib"
	"github.com/wq1019/short/internal/model"
)

func (s *Service) CreateDomain(c context.Context, request *model.CreateDomainRequest) (shortUrl string, err error) {
	// 原链接hash
	hash := lib.Md5(request.RedirectUrl)
	// 如果数据库中存在则直接返回短链接
	if s.dao.OriginUrlHashHasExist(c, hash) {
		domain, err1 := s.dao.LoadDomain(hash)
		if err1 != nil {
			err = errors.WithMessage(err1, "系统异常")
			return
		}
		return domain.ShortUrl, nil
	} else {
		// 生成4个短网址
		urls := lib.GenerateShortUrl(hash)
		for _, url := range urls {
			if exist := s.dao.ShortUrlHasBeenUsed(c, url); !exist {
				shortUrl = url
				break
			}
		}
		if shortUrl == "" {
			err = errors.New("短链接生成失败, 生成了重复的短链接")
			return
		}
		err = s.dao.CacheDomainUrl(c, hash, shortUrl)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		reply, err := s.dao.Publish(c, contants.CreateDomainChannel, request)
		if err != nil {
			return "", err
		}
		fmt.Println(reply)
	}
	return shortUrl, err
}
