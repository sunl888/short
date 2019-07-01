package dao

import "demo/internal/model"

func (d *Dao) FindUserByUserName(id int64) (*model.User, error) {
	return &model.User{}, nil
}
