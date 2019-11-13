package shttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/syncfuture/go/surl"

	u "github.com/syncfuture/go/util"
)

const (
	CNT_TYPE      = "Content-Type"
	CNT_TYPE_JSON = "application/json; charset=utf-8"
)

var (
	BytesPool = &sync.Pool{
		New: func() interface{} {
			b := make([]byte, 32)
			return &b
		},
	}
)

type APIClient struct {
	urlProvider surl.IURLProvider
}

// NewAPIClient create new api client
// urlProvider is optional, if doesn't provide, it will not render url
func NewAPIClient(urlProvider surl.IURLProvider) *APIClient {
	return &APIClient{
		urlProvider: urlProvider,
	}
}

// CallAPI call api
// bodyObj is optional
func (x *APIClient) CallAPI(client *http.Client, method, url string, bodyObj interface{}) *[]byte {
	return x.SendRequest(client, method, url, nil, bodyObj)
}

func (x *APIClient) SendRequest(client *http.Client, method, url string, configRequest func(*http.Request), bodyObj interface{}) *[]byte {
	buffer := BytesPool.Get().(*[]byte)

	var err error
	var request *http.Request

	if x.urlProvider != nil {
		// 渲染Url
		url = x.urlProvider.RenderURLCache(url)
	}

	// 创建Request
	if bodyObj != nil {
		var body []byte

		switch v := bodyObj.(type) {
		case []byte:
			body = v
			break
		case string:
			body = []byte(v)
			break
		default:
			body, err = json.Marshal(v)
			if u.LogError(err) {
				*buffer = []byte(err.Error())
				return buffer
			}
		}

		request, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		request, err = http.NewRequest(method, url, nil)
	}
	if u.LogError(err) {
		*buffer = []byte(err.Error())
		return buffer
	}

	// 配置Request
	request.Header.Set(CNT_TYPE, CNT_TYPE_JSON)
	if configRequest != nil {
		configRequest(request)
	}

	// 发送请求
	resp, err := client.Do(request)
	if u.LogError(err) {
		*buffer = []byte(err.Error())
		return buffer
	}
	defer resp.Body.Close()

	// 读取Response Body
	*buffer, err = ioutil.ReadAll(resp.Body)
	if u.LogError(err) {
		*buffer = []byte(err.Error())
		return buffer
	}

	return buffer
}
