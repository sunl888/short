package model

import (
	xtime "github.com/bilibili/kratos/pkg/time"
)

// short domain model
type Domain struct {
	Id          int    // 自增长ID
	RedirectUrl string // 重定向的链接
	ShortUrl    string // 短链接
	HitCount    int64  // 访问量
	IsPublish   bool   // 是否公开访问数据
	CreatedAt   xtime.Time
	UpdatedAt   xtime.Time
}

type QueryDomainResponse struct {
	DomainList []*Domain `json:"domain_list"`
	Pagination
}

// QueryApplyRequest request model
type QueryDomainRequest struct {
	Domain
	Pagination
}

// TableName get table name model
func (w Domain) TableName() string {
	return "domain"
}
