package databases

import (
	"fmt"
	"github.com/ava-cn/trading-central-playlists/app/models"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"net/url"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() (db *gorm.DB) {
	var (
		driverName string
		host       string
		port       int
		database   string
		user       string
		password   string
		charset    string
		local      string

		args string
		err  error
	)

	driverName = viper.GetString("database.driverName")
	host = viper.GetString("database.host")
	port = viper.GetInt("database.port")
	database = viper.GetString("database.database")
	user = viper.GetString("database.user")
	password = viper.GetString("database.password")
	charset = viper.GetString("database.charset")
	local = viper.GetString("database.local")

	args = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=%s",
		user,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(local),
	)

	if db, err = gorm.Open(driverName, args); err != nil {
		panic("failed to connect database,err: " + err.Error())
	}

	// 设置引擎和字符集
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8 auto_increment=1")

	// 自动创建数据表
	db.AutoMigrate(&models.Videos{})

	DB = db
	return
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
