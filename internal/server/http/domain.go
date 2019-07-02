package http

import (
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"github.com/bilibili/kratos/pkg/net/http/blademaster/binding"
	"github.com/wq1019/short/internal/model"
)

func createDomain(c *bm.Context) {
	domain := model.CreateDomainRequest{}
	if err := c.BindWith(&domain, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	shortUrl, err := svc.CreateDomain(c.Context, &domain)
	c.JSON(shortUrl, err)
}
