package shttp

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/syncfuture/go/spool"
	"github.com/syncfuture/go/surl"
)

const (
	CNT_TYPE      = "Content-Type"
	CNT_TYPE_JSON = "application/json; charset=utf-8"
)

var (
	_bufferPool = spool.NewSyncBufferPool(1024)
)

type APIClient struct {
	URLProvider surl.IURLProvider
}

func (x *APIClient) DoBuffer(client *http.Client, method, url string, configRequest func(*http.Request), bodyObj interface{}) (buffer *bytes.Buffer, err error) {
	buffer = _bufferPool.GetBuffer()

	var resp *http.Response
	resp, err = x.Do(client, method, url, configRequest, bodyObj)
	if err != nil {
		return nil, err
	}
	// 读取Response Body
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return buffer, err
}
func (x *APIClient) RecycleBuffer(buffer *bytes.Buffer) {
	_bufferPool.PutBuffer(buffer)
}

func (x *APIClient) Do(client *http.Client, method, url string, configRequest func(*http.Request), bodyObj interface{}) (resp *http.Response, err error) {

	var request *http.Request

	if x.URLProvider != nil {
		// 渲染Url
		url = x.URLProvider.RenderURLCache(url)
	}

	// 创建Request
	bodyBuffer := _bufferPool.GetBuffer()
	if bodyObj != nil {
		switch v := bodyObj.(type) {
		case []byte:
			bodyBuffer.Write(v)
			break
		case string:
			bodyBuffer.WriteString(v)
			break
		default:
			var body []byte
			body, err = json.Marshal(v)
			if err != nil {
				return nil, err
			}
			bodyBuffer.Write(body)
		}

		request, err = http.NewRequest(method, url, bodyBuffer)
	} else {
		request, err = http.NewRequest(method, url, nil)
	}
	defer _bufferPool.PutBuffer(bodyBuffer)

	if err != nil {
		return nil, err
	}

	// 配置Request
	request.Header.Set(CNT_TYPE, CNT_TYPE_JSON)
	if configRequest != nil {
		configRequest(request)
	}

	// 发送请求
	resp, err = client.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, err
}
