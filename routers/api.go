package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 路由
func InitRouters(r *gin.Engine) *gin.Engine {

	r.GET("ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	return r
}
