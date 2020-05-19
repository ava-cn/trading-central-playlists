package configs

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/viper"
)

var defaultConfig = []byte(`
server:
  port: 80

database:
  driverName: mysql
  host: mysql
  port: 3306
  database: trading_central_playlists
  user: root
  password: root
  charset: utf8
  local: Asia/Shanghai

app:
  xmlURL: https://video.tradingcentral.com/playlists/23125.xml

qiniu:
  enabled: false
  bucket:
  privateBucket:
  accessKey:
  secretKey:
  domain:
  useHTTPS:
  useCdnDomains:
`)

// ConfYaml is config structure.
type ConfYaml struct {
	Server   SectionServer   `yaml:"server"`
	Database SectionDatabase `yaml:"database"`
	App      SectionApp      `yaml:"app"`
	Qiniu    SectionQiNiu    `yaml:"qiniu"`
}

// SectionServer config for HTTP
type SectionServer struct {
	Port string `yaml:"port"`
}

// SectionDatabases config for database
type SectionDatabase struct {
	DriverName string `yaml:"driverName"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Database   string `yaml:"database"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Charset    string `yaml:"charset"`
	Local      string `yaml:"local"`
}

// AppSection config for app custom config
type SectionApp struct {
	XmlURL string `yaml:"xmlURL"`
}

// SectionQiniu config for qiniu storage
type SectionQiNiu struct {
	Enabled       bool   `yaml:"enabled"`
	Bucket        string `yaml:"bucket"`
	PrivateBucket bool   `yaml:"privateBucket"`
	AccessKey     string `yaml:"accessKey"`
	SecretKey     string `yaml:"secretKey"`
	Domain        string `yaml:"domain"`
	UseHTTPS      bool   `yaml:"useHttps"`
	UseCdnDomains bool   `yaml:"useCdnDomains"`
}

// LoadConf load config from file and read in environment variables that match
func LoadConf(confPath string) (ConfYaml, error) {
	var (
		config ConfYaml
		err    error
	)

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()                  // read in environment variables that match
	viper.SetEnvPrefix("trading_central") // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)

		if err != nil {
			return config, err
		}

		if err = viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return config, err
		}
	} else {
		// Search config in home directory with name ".gorush" (without extension).
		viper.AddConfigPath("/etc/trading-central-playlists/")
		viper.AddConfigPath("$HOME/.trading-central-playlists")
		viper.AddConfigPath(".")
		viper.SetConfigName("app")

		// If a config file is found, read it in.
		if err = viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			// load default config
			if err = viper.ReadConfig(bytes.NewBuffer(defaultConfig)); err != nil {
				return config, err
			}
		}
	}

	// Server
	config.Server.Port = viper.GetString("server.port")

	// Database
	config.Database.DriverName = viper.GetString("database.driverName")
	config.Database.Host = viper.GetString("database.host")
	config.Database.Port = viper.GetInt("database.port")
	config.Database.Database = viper.GetString("database.database")
	config.Database.User = viper.GetString("database.user")
	config.Database.Password = viper.GetString("database.password")
	config.Database.Charset = viper.GetString("database.charset")
	config.Database.Charset = viper.GetString("database.local")

	// App
	config.App.XmlURL = viper.GetString("app.xmlURL")

	// Qiniu
	config.Qiniu.Enabled = viper.GetBool("qiniu.enabled")
	config.Qiniu.Bucket = viper.GetString("qiniu.bucket")
	config.Qiniu.PrivateBucket = viper.GetBool("qiniu.privateBucket")
	config.Qiniu.AccessKey = viper.GetString("qiniu.AccessKey")
	config.Qiniu.SecretKey = viper.GetString("qiniu.secretKey")
	config.Qiniu.Domain = viper.GetString("qiniu.domain")
	config.Qiniu.UseHTTPS = viper.GetBool("qiniu.useHTTPS")
	config.Qiniu.UseCdnDomains = viper.GetBool("qiniu.useCdnDomains")

	// Default Value
	if config.Server.Port == "" {
		config.Server.Port = "80"
	}

	return config, nil
}

// InitConf
func InitConf() (config ConfYaml, err error) {
	var (
		configFile string
	)

	flag.StringVar(&configFile, "c", "", "Configuration file path.")
	flag.StringVar(&configFile, "config", "", "Configuration file path.")
	flag.Parse()

	if config, err = LoadConf(configFile); err != nil {
		log.Printf("Load yaml config file error: '%v'", err)
		return
	}
	return config, nil
}
