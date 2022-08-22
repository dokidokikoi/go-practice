package core

import (
	"context"
	"errors"

	"user/model"
	"user/service/pb"

	"github.com/jinzhu/gorm"
)

func BUildeUser(user *model.User) *pb.UserModel {
	return &pb.UserModel{
		UserId:   uint32(user.ID),
		UserName: user.UserName,
	}
}

func IsUserExist(userName string) bool {
	count := 0
	err := model.DB.Model(&model.User{}).Where("user_name=?", userName).Count(&count).Error
	return count > 0 && err == nil
}

func (*UserService) UserLogin(ctx context.Context, req *pb.UserRequest, resp *pb.UserDefaultResponse) error {
	var user model.User
	resp.Code = 200
	if err := model.DB.Where("user_name=?", req.UserName).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			resp.Code = 400
			return nil
		}
		resp.Code = 500
		return nil
	}
	if !user.CheckPassword(req.Password) {
		resp.Code = 400
		return nil
	}

	resp.UserDetail = BUildeUser(&user)
	return nil
}

func (*UserService) UserRegister(ctx context.Context, req *pb.UserRequest, resp *pb.UserDefaultResponse) error {
	if req.Password != req.PasswordConfirm {
		return errors.New("两次密码输入不一致")
	}

	if IsUserExist(req.UserName) {
		return errors.New("用户已存在")
	}

	user := &model.User{
		UserName: req.UserName,
	}
	if err := user.SetPassword(req.Password); err != nil {
		return err
	}

	if err := model.DB.Create(&user).Error; err != nil {
		return err
	}
	resp.UserDetail = BUildeUser(user)
	return nil
}
