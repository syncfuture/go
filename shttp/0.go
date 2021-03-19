package shttp

import (
	"bytes"
	"net/http"

	"github.com/syncfuture/go/spool"
)

const (
	HEADER_CTYPE = "Content-Type"
	HEADER_AUTH  = "Authorization"
	CTYPE_TEXT   = "text/plain"
	CTYPE_HTML   = "text/html"
	CTYPE_XML    = "text/xml"
	CTYPE_CSS    = "text/css"
	CTYPE_JS     = "text/javascript"
	CTYPE_JSON   = "application/json"
	CTYPE_FORM   = "application/x-www-form-urlencoded"
	CTYPE_MFORM  = "multipart/form-data"
	CHARSET_UTF8 = "charset=utf-8"
)

var (
	_bufferPool = spool.NewSyncBufferPool(1024)
)

func GetRespBuffer(resp *http.Response, err error) (*bytes.Buffer, error) {
	if err != nil {
		return nil, err
	}

	bf := _bufferPool.GetBuffer()
	defer func() {
		_bufferPool.PutBuffer(bf)
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	_, err = bf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return bf, nil
}

func RecycleBuffer(buffer *bytes.Buffer) {
	_bufferPool.PutBuffer(buffer)
}
