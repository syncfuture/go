package core

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/syncfuture/go/json"
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

func CallAPI(client *http.Client, method, apiUrl string, datas ...interface{}) *[]byte {
	return CallAPIRequest(client, method, apiUrl, nil, datas)
}

func CallAPIRequest(client *http.Client, method, apiUrl string, configRequest func(*http.Request), body ...interface{}) *[]byte {
	buffer := BytesPool.Get().(*[]byte)

	var err error
	var request *http.Request

	if len(body) == 1 {
		var dataJson string
		dataJson, err = json.Serialize(body[0])
		if u.LogError(err) {
			*buffer = []byte(err.Error())
			return buffer
		}
		request, err = http.NewRequest(method, apiUrl, bytes.NewBufferString(dataJson))

	} else {
		request, err = http.NewRequest(method, apiUrl, nil)
	}
	if u.LogError(err) {
		*buffer = []byte(err.Error())
		return buffer
	}
	// 设置读取
	request.Header.Set(CNT_TYPE, CNT_TYPE_JSON)
	if configRequest != nil {
		configRequest(request)
	}

	resp, err := client.Do(request)
	if u.LogError(err) {
		*buffer = []byte(err.Error())
		return buffer
	}
	defer resp.Body.Close()

	*buffer, err = ioutil.ReadAll(resp.Body)
	if u.LogError(err) {
		*buffer = []byte(err.Error())
		return buffer
	}

	return buffer
}
