package main

import (
	"github.com/ava-cn/trading-central-playlists/configs"
	"github.com/ava-cn/trading-central-playlists/databases"
	"github.com/ava-cn/trading-central-playlists/routers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

func main() {
	var (
		r    *gin.Engine
		db   *gorm.DB
		port string
	)

	// 初始化配置文件
	configs.InitConfigs()

	// 初始化数据库连接
	db = databases.InitDB()
	defer db.Close()

	r = gin.Default()

	// 载入路由
	r = routers.InitRouters(r)

	// 自定义http服务端口
	port = viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}