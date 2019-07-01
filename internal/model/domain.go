package model

import "time"

// short domain model
type Domain struct {
	Id          int
	RedirectUrl string
	ShortUrl    string
	HitCount    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
