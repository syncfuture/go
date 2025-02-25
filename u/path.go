package u

import (
	"net/url"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/syncfuture/go/serr"
)

// JointURL join urls
func JointURL(basePath string, paths ...string) (*url.URL, error) {
	r, err := url.Parse(basePath)

	if err != nil {
		return nil, errors.Errorf("invalid url")
	}

	p2 := append([]string{r.Path}, paths...)

	result := path.Join(p2...)

	r.Path = result

	return r, nil
}

// JointURLString join urls as string
func JointURLString(basePath string, paths ...string) (string, error) {
	r, err := JointURL(basePath, paths...)
	if LogError(err) {
		return "", err
	}

	return r.String(), nil
}

// AbsURL Get a absolute url
func AbsURL(baseURL, inputURL string) (string, error) {
	parsedInput, err := url.Parse(inputURL)
	if err != nil {
		return "", serr.WithStack(err)
	}

	// 如果 inputURL 已经是绝对URL，则直接返回
	if parsedInput.IsAbs() {
		return parsedInput.String(), nil
	}

	// 解析 baseURL
	parsedBase, err := url.Parse(baseURL)
	if err != nil {
		return "", serr.WithStack(err)
	}

	// 处理 / 开头的路径（基于站点根目录）
	if strings.HasPrefix(inputURL, "/") {
		rootBase := &url.URL{
			Scheme: parsedBase.Scheme,
			Host:   parsedBase.Host,
		}
		return rootBase.ResolveReference(parsedInput).String(), nil
	}

	// 否则，使用 baseURL 作为基础路径
	return parsedBase.ResolveReference(parsedInput).String(), nil
}
