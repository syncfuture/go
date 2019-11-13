package surl

type IURLProvider interface {
	GetURL(urlKey string) string
	GetURLCache(urlKey string) string
}
