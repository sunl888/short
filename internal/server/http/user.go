package http

import (
	bm "go-common/library/net/http/blademaster"
)

func queryUserInfo(c *bm.Context) {
	req := new(struct {
		Id int64 `form:"id"`
	})
	if err := c.Bind(&req); err != nil {
		return
	}
	c.JSON(svc.QueryUserInfo(c, 1))
}
