package u

import (
	"fmt"
	"net/url"
	"path"
)

// JointURL join urls
func JointURL(basePath string, paths ...string) (*url.URL, error) {
	r, err := url.Parse(basePath)

	if err != nil {
		return nil, fmt.Errorf("invalid url")
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
