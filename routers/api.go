package routers

import (
	"log"
	"net/http"

	"github.com/ava-cn/trading-central-playlists/app/http/resources"
	"github.com/ava-cn/trading-central-playlists/app/models"
	"github.com/ava-cn/trading-central-playlists/databases"
	"github.com/gin-gonic/gin"
)

// 路由
func InitRouters(r *gin.Engine) *gin.Engine {

	r.GET("ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/playlists/latest", func(ctx *gin.Context) {
		var (
			video models.Videos
		)
		// 获取最新的一条数据
		databases.GetDB().
			Order("video_id desc").
			Where("synced = ?", true).
			First(&video)

		ctx.JSON(http.StatusOK, gin.H{"data": resources.VideoShow(&video), "message": "success", "code": http.StatusOK})
	})

	r.GET("/playlists", func(ctx *gin.Context) {
		var (
			videos []*models.Videos
			total  uint64
			err    error
		)
		params := struct {
			Page  int `form:"page,default=1"`
			Limit int `form:"limit,default=10"`
		}{}

		if err = ctx.ShouldBind(&params); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "数据验证错误", "code": http.StatusUnprocessableEntity, "data": gin.H{}})
			return
		}

		// 获取已经同步的数据
		if videos, total, err = models.ListVideo(databases.GetDB(), params.Page, params.Limit); err != nil {
			log.Println(err)
		}

		ctx.JSON(http.StatusOK, gin.H{"data": resources.VideoCollection(videos), "current_page": params.Page, "total": total, "message": "success", "code": http.StatusOK})
	})

	return r
}
