package middleware

import (
	"fmt"
	"gateway/pkg/e"
	"gateway/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var code int
		code = 200
		token := ctx.GetHeader("Authorization")
		if token == "" {
			code = 404
		} else {
			fmt.Println("in JWT middleware")
			claim, err := util.ParseToken(token)
			fmt.Println("claim", claim)
			if err != nil {
				fmt.Println("1")
				code = e.ErrorAuthCheckTokenFail
			} else if time.Now().Unix() > claim.ExpiresAt {
				fmt.Println("2")
				code = e.ErrorAuthCheckTokenTimeout
			}
		}

		if code != 200 {
			ctx.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(uint(code)),
			})
			ctx.Abort()
			return
		}
		fmt.Println("in JWT middleware Next")
		ctx.Next()
	}
}
