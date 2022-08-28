package handler

import (
	"context"

	"user/internal/repository"
	"user/internal/service/pb"
	"user/pkg/e"
)

type UserSerivce struct{}

func NewUserService() *UserSerivce {
	return &UserSerivce{}
}

func (*UserSerivce) UserLogin(ctx context.Context, req *pb.UserRequest) (resp *pb.UserDefaultResponse, err error) {
	var user repository.User
	resp = new(pb.UserDefaultResponse)
	resp.Code = e.Success
	err = user.ShowUserInfo(req)
	if err != nil {
		resp.Code = e.Error
		return resp, err
	}
	resp.UserDetail = repository.BuilderUser(user)
	return resp, nil
}

func (*UserSerivce) UserRegister(ctx context.Context, req *pb.UserRequest) (resp *pb.UserDefaultResponse, err error) {
	var user repository.User
	resp = new(pb.UserDefaultResponse)
	resp.Code = e.Success
	user, err = user.UserCreate(req)
	if err != nil {
		resp.Code = e.Error
		return resp, err
	}
	resp.UserDetail = repository.BuilderUser(user)
	return resp, nil
}
