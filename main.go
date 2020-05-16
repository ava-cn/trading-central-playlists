package main

import (
	"github.com/ava-cn/trading-central-playlists/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	var (
		r *gin.Engine
	)

	r = gin.Default()

	// 载入路由
	r = routers.InitRouters(r)

	panic(r.Run())
}
