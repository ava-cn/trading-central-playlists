package supports

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"path"
)

// 获取重定向后的URL地址
func GetRedirectURL(originURL string) (redirectURL string, err error) {
	var (
		client       *http.Client
		httpRequest  *http.Request
		httpResponse *http.Response
		responseURL  *url.URL
	)
	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			//自用，将url根据需求进行组合
			if len(via) >= 1 {
				return errors.New("stopped after 1 redirects")
			}
			return nil
		},
	}

	if httpRequest, err = http.NewRequest("GET", originURL, nil); err != nil {
		return
	}

	if httpResponse, err = client.Do(httpRequest); err != nil {
		//log.Println(httpResponse.StatusCode)
	}

	if responseURL, err = httpResponse.Location(); err != nil {
		log.Println(responseURL.String())
		return
	}

	redirectURL = responseURL.String()

	return
}

func GetFileNameFromURL(originURL string) (fileName string, err error) {
	var (
		httpRequest *http.Request
	)

	if httpRequest, err = http.NewRequest("GET", originURL, nil); err != nil {
		return
	}

	fileName = path.Base(httpRequest.URL.Path)

	return
}
