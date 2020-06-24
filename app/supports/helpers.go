package supports

import (
	"net/http"
	"path"
)

// 获取重定向后的URL地址
func GetRedirectURL(originURL string) (redirectURL string, err error) {
	var (
		client       *http.Client
		httpResponse *http.Response
	)

	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	if httpResponse, err = client.Get(originURL); err != nil {
		redirectURL = originURL
		return
	}

	if httpResponse.StatusCode == 301 || httpResponse.StatusCode == 302 {
		redirectURL = httpResponse.Header.Get("Location")
		return
	}

	redirectURL = originURL

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
