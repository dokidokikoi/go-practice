package repository

import (
	"errors"
	"fmt"
	"user/internal/service/pb"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	UserId   uint   `gorm:"primarykey"`
	UserName string `gorm:"unique"`
	NickName string
	Password string
}

const (
	passwordCost = 12 // 密码加密难度
)

func (user *User) CheckUserExist(req *pb.UserRequest) bool {
	if err := DB.Where("user_name=?", req.UserName).First(&user).Error; err == gorm.ErrRecordNotFound {
		return false
	}

	return true
}

func (user *User) ShowUserInfo(req *pb.UserRequest) error {
	if exist := user.CheckUserExist(req); exist {
		return nil
	}
	return errors.New("User Not Found")
}

func (user *User) UserCreate(req *pb.UserRequest) (User, error) {
	if exist := user.CheckUserExist(req); exist {
		return User{}, errors.New("UserName Already Exist")
	}

	u := User{
		UserName: req.UserName,
		NickName: req.NickName,
	}

	_ = u.SetPassword(req.Password)
	err := DB.Create(&u).Error
	fmt.Println("user", u)
	return u, err
}

func (user *User) SetPassword(pwd string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), passwordCost)

	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
	return err == nil
}

func BuilderUser(item User) *pb.UserModel {
	return &pb.UserModel{
		UserId:   uint32(item.UserId),
		UserName: item.UserName,
		NickName: item.NickName,
	}
}
