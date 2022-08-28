package handlers

import (
	"context"
	"gateway/pkg/utils"
	"net/http"

	"gateway/services/pb"

	"github.com/gin-gonic/gin"
)

func UserRegister(ctx *gin.Context) {
	var userReq pb.UserRequest
	PanicIfUserError(ctx.Bind(&userReq))
	userService := ctx.Keys["userService"].(pb.UserService)
	userResp, err := userService.UserRegister(context.Background(), &userReq)
	PanicIfUserError(err)
	ctx.JSON(http.StatusOK, gin.H{"data": userResp})
}

func UserLogin(ctx *gin.Context) {
	var userReq pb.UserRequest
	PanicIfUserError(ctx.Bind(&userReq))
	userService := ctx.Keys["userService"].(pb.UserService)
	userResp, err := userService.UserLogin(context.Background(), &userReq)
	PanicIfUserError(err)
	token, _ := utils.GenerateToken(uint(userResp.UserDetail.UserId))
	ctx.JSON(http.StatusOK, gin.H{
		"code": userResp.Code,
		"data": gin.H{
			"user":  userResp.UserDetail,
			"token": token,
		},
	})
}
