package main

import (
	"log"

	"github.com/ava-cn/trading-central-playlists/app/consoles"
	"github.com/ava-cn/trading-central-playlists/configs"
	"github.com/ava-cn/trading-central-playlists/databases"
	"github.com/ava-cn/trading-central-playlists/routers"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	var (
		config configs.ConfYaml
		r      *gin.Engine
		db     *gorm.DB
		port   string
		err    error
	)

	// 初始化配置文件
	if config, err = configs.InitConf(); err != nil {
		log.Panicf("failed to init config,err: %s", err)
		return
	}

	// 初始化数据库连接
	db = databases.InitDB()
	defer db.Close()

	// 执行定时任务
	go consoles.InitCorn()

	// 启动Gin框架
	r = gin.Default()

	// 载入路由
	r = routers.InitRouters(r)

	// 自定义http服务端口
	port = config.Server.Port
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}
