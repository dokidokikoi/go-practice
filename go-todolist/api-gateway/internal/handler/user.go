package handler

import (
	"context"
	"fmt"
	"net/http"

	"gateway/internal/service/pb"
	"gateway/pkg/e"
	"gateway/pkg/res"
	"gateway/pkg/util"

	"github.com/gin-gonic/gin"
)

func UserRegister(ctx *gin.Context) {
	var userReq pb.UserRequest
	PanicIfUserError(ctx.Bind(&userReq))
	// gin Key 中获取服务实例
	userService := ctx.Keys["user"].(pb.UserServiceClient)
	userResp, err := userService.UserRegister(context.Background(), &userReq)
	PanicIfUserError(err)
	r := res.Response{
		Data:   userResp,
		Status: uint(userResp.Code),
		Msg:    e.GetMsg(uint(userResp.Code)),
		// Error:  err.Error(),
	}
	ctx.JSON(http.StatusOK, r)
}

func UserLogin(ctx *gin.Context) {
	var userReq pb.UserRequest
	PanicIfUserError(ctx.Bind(&userReq))
	// gin Key 中获取服务实例
	userService := ctx.Keys["user"].(pb.UserServiceClient)
	userResp, err := userService.UserLogin(context.Background(), &userReq)
	fmt.Println("resp", userResp)
	fmt.Println("userResp", userResp)
	PanicIfUserError(err)
	token, err := util.GenerateToken(uint(userResp.UserDetail.UserId))
	r := res.Response{
		Data: res.TokenData{
			User:  userResp.UserDetail,
			Token: token,
		},
		Status: uint(userResp.Code),
		Msg:    e.GetMsg(uint(userResp.Code)),
		// Error:  err.Error(),
	}
	ctx.JSON(http.StatusOK, r)
}
