package consoles

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// UploadToQiniu 通过给定远程资源地址上传到七牛
func UploadToQiniu(remoteURL string, distName string) (key string, err error) {
	var (
		upToken      string
		putRet       storage.PutRet
		formUploader *storage.FormUploader
		putExtra     storage.PutExtra
		readByte     []byte
		transCfg     *http.Transport
		response     *http.Response
		client       *http.Client
	)

	if formUploader, upToken, err = InitQiniu(); err != nil {
		return
	}

	putRet = storage.PutRet{}
	putExtra = storage.PutExtra{}

	// Create New http Transport
	transCfg = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // disable verify
	}
	// Create Http Client
	client = &http.Client{Transport: transCfg}

	if response, err = client.Get(remoteURL); err != nil {
		log.Println(err)
		return
	}

	defer response.Body.Close()
	if readByte, err = ioutil.ReadAll(response.Body); err != nil {
		log.Println(err)
		return
	}

	if err = formUploader.Put(context.Background(), &putRet, upToken, distName, bytes.NewReader(readByte), int64(len(readByte)), &putExtra); err != nil {
		log.Println(err)
		return
	}

	key = putRet.Key

	return
}

// InitQiniu 初始化七牛云存储
func InitQiniu() (formUploader *storage.FormUploader, upToken string, err error) {
	var (
		bucket        string
		accessKey     string
		secretKey     string
		useHTTPS      bool
		useCdnDomains bool
		putPolicy     storage.PutPolicy
		mac           *qbox.Mac
		zone          *storage.Zone
		config        *storage.Config
	)
	bucket = viper.GetString("qiniu.bucket")
	accessKey = viper.GetString("qiniu.accessKey")
	secretKey = viper.GetString("qiniu.secretKey")
	useHTTPS = viper.GetBool("qiniu.useHTTPS")
	useCdnDomains = viper.GetBool("qiniu.useCdnDomains")

	putPolicy = storage.PutPolicy{
		Scope: bucket,
	}
	mac = qbox.NewMac(accessKey, secretKey)
	upToken = putPolicy.UploadToken(mac)

	if zone, err = storage.GetZone(accessKey, bucket); err != nil { // 空间对应的机房
		return
	}

	config = &storage.Config{
		Zone:          zone,          // bucket所在地区
		UseHTTPS:      useHTTPS,      // 是否使用https域名
		UseCdnDomains: useCdnDomains, // 上传是否使用CDN上传加速
	}
	formUploader = storage.NewFormUploader(config)
	return
}

// GetFile 通过给定key获取bucket中的资源
func GetFile(key string) string {
	var (
		accessKey     string
		secretKey     string
		useHTTPS      bool
		domain        string
		privateBucket bool
		mac           *auth.Credentials
		deadline      int64
	)

	accessKey = viper.GetString("qiniu.accessKey")
	secretKey = viper.GetString("qiniu.secretKey")
	privateBucket = viper.GetBool("qiniu.privateBucket")
	domain = viper.GetString("qiniu.domain")
	useHTTPS = viper.GetBool("qiniu.useHTTPS")

	if useHTTPS {
		domain = fmt.Sprintf("%s://%s", "https", domain)
	} else {
		domain = fmt.Sprintf("%s://%s", "http", domain)
	}

	if privateBucket { // 私有空间访问
		mac = auth.New(accessKey, secretKey)
		deadline = time.Now().Add(time.Second * 3600).Unix() // 1小时有效期
		return storage.MakePrivateURL(mac, domain, key, deadline)
	} else { // 公开空间访问
		return storage.MakePublicURL(domain, key)
	}
}
