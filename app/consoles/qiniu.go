package consoles

import (
	"bytes"
	"context"
	"fmt"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

// 通过给定远程资源地址上传到七牛
func UploadToQiniu(remoteURL string, distName string) (key string, err error) {
	var (
		upToken      string
		putRet       storage.PutRet
		formUploader *storage.FormUploader
		putExtra     storage.PutExtra
		u            bytes.Buffer
		readByte     []byte
	)

	if formUploader, upToken, err = InitQiniu(); err != nil {
		return
	}

	putRet = storage.PutRet{}
	putExtra = storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	u = bytes.Buffer{}
	u.WriteString(remoteURL)
	response, _ := http.Get(u.String())
	defer response.Body.Close()
	readByte, _ = ioutil.ReadAll(response.Body)

	if err = formUploader.Put(context.Background(), &putRet, upToken, distName, bytes.NewReader(readByte), int64(len(readByte)), &putExtra); err != nil {
		fmt.Println(err)
		return
	}

	key = putRet.Key

	return
}

// 初始化七牛云存储
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
