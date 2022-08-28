package handler

import (
	"net/http"

	"context"
	"gateway/internal/service/pb"
	"gateway/pkg/e"
	"gateway/pkg/res"
	"gateway/pkg/util"

	"github.com/gin-gonic/gin"
)

func ListTask(ctx *gin.Context) {
	var tReq pb.TaskRequest
	PanicIfTaskError(ctx.Bind(&tReq))
	claim, _ := util.ParseToken(ctx.GetHeader("Authorization"))
	tReq.UserID = uint32(claim.UserID)
	taskservice := ctx.Keys["task"].(pb.TaskServiceClient)
	taskResp, err := taskservice.TaskShow(context.Background(), &tReq)
	PanicIfTaskError(err)
	r := res.Response{
		Status: uint(taskResp.Code),
		Data:   taskResp,
		Msg:    e.GetMsg(uint(taskResp.Code)),
	}

	ctx.JSON(http.StatusOK, r)
}

func CreateTask(ctx *gin.Context) {
	var tReq pb.TaskRequest
	PanicIfTaskError(ctx.Bind(&tReq))
	claim, _ := util.ParseToken(ctx.GetHeader("Authorization"))
	tReq.UserID = uint32(claim.UserID)
	taskservice := ctx.Keys["task"].(pb.TaskServiceClient)
	taskResp, err := taskservice.TaskCreate(context.Background(), &tReq)
	PanicIfTaskError(err)
	r := res.Response{
		Status: uint(taskResp.Code),
		Data:   taskResp,
		Msg:    e.GetMsg(uint(taskResp.Code)),
	}

	ctx.JSON(http.StatusOK, r)
}

func UpdateTask(ctx *gin.Context) {
	var tReq pb.TaskRequest
	PanicIfTaskError(ctx.Bind(&tReq))
	claim, _ := util.ParseToken(ctx.GetHeader("Authorization"))
	tReq.UserID = uint32(claim.UserID)
	taskservice := ctx.Keys["task"].(pb.TaskServiceClient)
	taskResp, err := taskservice.TaskUpdate(context.Background(), &tReq)
	PanicIfTaskError(err)
	r := res.Response{
		Status: uint(taskResp.Code),
		Data:   taskResp,
		Msg:    e.GetMsg(uint(taskResp.Code)),
	}

	ctx.JSON(http.StatusOK, r)
}

func DeleteTask(ctx *gin.Context) {
	var tReq pb.TaskRequest
	PanicIfTaskError(ctx.Bind(&tReq))
	claim, _ := util.ParseToken(ctx.GetHeader("Authorization"))
	tReq.UserID = uint32(claim.UserID)
	taskservice := ctx.Keys["task"].(pb.TaskServiceClient)
	taskResp, err := taskservice.TaskDelete(context.Background(), &tReq)
	PanicIfTaskError(err)
	r := res.Response{
		Status: uint(taskResp.Code),
		Data:   taskResp,
		Msg:    e.GetMsg(uint(taskResp.Code)),
	}

	ctx.JSON(http.StatusOK, r)
}
