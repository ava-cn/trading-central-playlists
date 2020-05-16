package configs

import (
	"github.com/spf13/viper"
	"os"
)

// 初始化配置文件
func InitConfigs() {
	var (
		workDir string
		err     error
	)
	workDir, err = os.Getwd()
	viper.SetConfigName("app")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/configs")
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
