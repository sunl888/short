package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/wq1019/short/internal/model"
)

const (
	_shortDomainRedisKey       = "short_domain:%s"
	_shortUrlListRedisKey      = "short_url_list"
	_originUrlHashListRedisKey = "origin_url_hash_list"
)

// 创建短网址记录
func (d *Dao) CreateDomain(c context.Context, domain *model.Domain) (err error) {
	err = d.CacheDomainByShortUrl(c, domain)
	if err != nil {
		return err
	}
	return d.db.Create(&domain).Error
}

// 已经录入的网址列表
func (d *Dao) QueryOrder(domain *model.Domain, pn int32, ps int32) (qor *model.QueryDomainResponse, err error) {
	qor = &model.QueryDomainResponse{}
	err = d.db.Table(model.Domain{}.TableName()).Where(model.Domain{
		Id: domain.Id, RedirectUrl: domain.RedirectUrl, ShortUrl: domain.ShortUrl}).
		Count(&qor.TotalSize).Offset((pn - 1) * ps).Limit(ps).Find(&qor.DomainList).Error
	qor.PageSize = ps
	qor.PageNum = pn
	return
}

// 为短网址创建缓存
func (d *Dao) CacheShortUrl(c context.Context, shortUrl string) (err error) {
	_, err = d.SAdd(c, _shortUrlListRedisKey, shortUrl)
	return
}

// 为源网址创建缓存
func (d *Dao) CacheOriginUrlHash(c context.Context, originUrlHash string) (err error) {
	_, err = d.SAdd(c, _originUrlHashListRedisKey, originUrlHash)
	return
}

// 缓存短网址 model
func (d *Dao) CacheDomainByShortUrl(c context.Context, domain *model.Domain) (err error) {
	var body []byte
	if body, err = json.Marshal(&domain); err != nil {
		err = errors.WithStack(err)
		return
	}
	if ok, err := d.Set(c, fmt.Sprintf(_shortDomainRedisKey, domain.ShortUrl), body); !ok || err != nil {
		err = errors.WithMessage(err, "缓存短网址 model 失败")
		return
	}
	return
}

// 判断这个短网址是否被使用
func (d *Dao) ShortUrlHasBeenUsed(c context.Context, shortUrl string) (hasBeenUsed bool) {
	hasBeenUsed, err := d.SIsMember(c, _shortUrlListRedisKey, shortUrl)
	if err != nil {
		panic(err)
	}
	return
}

// 判断这个源网址是否存在
func (d *Dao) OriginUrlHashHasExist(c context.Context, originUrl string) (hasBeenUsed bool) {
	hasBeenUsed, err := d.SIsMember(c, _originUrlHashListRedisKey, originUrl)
	if err != nil {
		panic(err)
	}
	return
}
