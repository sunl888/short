package service

import (
	"context"
	"demo/internal/model"
)

func (s *Service) QueryUserInfo(c context.Context, id int64) (userInfo *model.User, err error) {
	return s.dao.FindUserByUserName(id)
}
