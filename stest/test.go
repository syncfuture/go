package stest

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"

	r "github.com/syncfuture/go/srand"
)

func AutoFixture(objPtr interface{}) {
	t := reflect.TypeOf(objPtr)
	if t.Kind() != reflect.Ptr {
		panic("objPtr must be pointer")
	}

	v := reflect.ValueOf(objPtr).Elem()

	fieldCount := v.NumField()
	for i := 0; i < fieldCount; i++ {
		field := v.Field(i)
		if field.CanAddr() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(r.String(10))
			case reflect.Int:
				field.SetInt(rand.Int63())
			case reflect.Int32:
				field.SetInt(rand.Int63())
			case reflect.Int64:
				field.SetInt(rand.Int63())
			case reflect.Int16:
				field.SetInt(rand.Int63())
			case reflect.Int8:
				field.SetInt(rand.Int63())
			case reflect.Float32:
				field.SetFloat(rand.Float64())
			case reflect.Float64:
				field.SetFloat(rand.Float64())
			}
		}
	}
}

func CreateTestHttpClient(proxyURL string) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 跳过证书验证
	}

	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			panic(err)
		}
		tr.Proxy = http.ProxyURL(proxy)
	}

	return &http.Client{Transport: tr}
}
