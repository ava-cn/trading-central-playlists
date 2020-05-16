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

// 通过给定地址上传文件
func UploadToQiniu(remoteURL string, distName string) (key string, err error) {
	var (
		upToken      string
		putRet       storage.PutRet
		formUploader *storage.FormUploader
		putExtra     storage.PutExtra
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

	url := bytes.Buffer{}
	url.WriteString(remoteURL)
	fmt.Println(url.String())
	response, _ := http.Get(url.String())
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	dataLen := int64(len(data))
	err = formUploader.Put(context.Background(), &putRet, upToken, distName, bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}

	key = putRet.Key

	fmt.Println(key, putRet.Hash, putRet.PersistentID)

	return
}

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
	accessKey = viper.GetString("qiniu.bucket")
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
