package shttp

import (
	"bytes"
	"net/http"

	"github.com/syncfuture/go/serr"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/spool"
	"github.com/syncfuture/go/u"
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
	_, err = bf.ReadFrom(resp.Body)
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return nil, serr.WithStack(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Warnf("%s %s [%d] -> %s", resp.Request.Method, resp.Request.URL.String(), resp.StatusCode, u.BytesToStr(bf.Bytes()))
	}

	return bf, nil
}

func RecycleBuffer(buffer *bytes.Buffer) {
	_bufferPool.PutBuffer(buffer)
}
