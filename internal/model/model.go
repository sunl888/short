package model

import "github.com/pkg/errors"

// Verify verify the value of pageNum and pageSize.
func (p *Pagination) Verify() error {
	if p.PageNum < 0 {
		return errors.New("分页页码不合法")
	} else if p.PageNum == 0 {
		p.PageNum = 1
	}
	if p.PageSize < 0 {
		return errors.New("分页大小不合法")
	} else if p.PageSize == 0 {
		p.PageSize = 10
	}
	return nil
}

// Pagination page num
type Pagination struct {
	PageNum   int32 `form:"page_num" json:"page_num"`
	PageSize  int32 `form:"page_size" json:"page_size"`
	TotalSize int32 `form:"total_size" json:"total_size"`
}
