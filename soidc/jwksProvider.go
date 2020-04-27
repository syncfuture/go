package soidc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	lru "github.com/hashicorp/golang-lru"
	log "github.com/syncfuture/go/slog"
	"gopkg.in/square/go-jose.v2"
)

var lruCache *lru.Cache

type cacheValue struct {
	Algorithm string
	Data      interface{}
	Expiry    time.Time
}

func init() {
	var err error
	lruCache, err = lru.New(128)
	if err != nil {
		log.Fatal("cannot initialize cache")
	}
}

type jwksProvider struct {
	url      string
	issuer   string
	audience string
}

func NewPublicKeyProvider(issuer, jwksPath, audience string) IPublicKeyProvider {
	url := issuer + jwksPath

	return &jwksProvider{
		url:      url,
		issuer:   issuer,
		audience: audience,
	}
}

// GetPublicKey verifies the desired iss and aud against the token's claims, and then
// tries to fetch a public key from the iss. It returns the public key as byte slice
// on success and error on failure.
func (x *jwksProvider) GetKey(token *jwt.Token) (interface{}, error) {
	claims := token.Claims.(jwt.MapClaims)

	// Get iss from JWT and validate against desired iss
	if claims["iss"].(string) != x.issuer {
		return nil, fmt.Errorf("cannot validate iss claim")
	}

	// Get audience from JWT and validate against desired audience
	if claims["aud"].(string) != x.audience {
		return nil, fmt.Errorf("cannot validate audience claim")
	}

	cacheKey := token.Header["kid"].(string) + "|" + claims["iss"].(string)

	// Try to get and return existing entry from cache. If cache is expired,
	// it will try proceed with rest of the function call
	cached, ok := lruCache.Get(cacheKey)
	if ok {
		val := cached.(*cacheValue)

		// Check for alg
		if val.Algorithm != token.Header["alg"] {
			return nil, fmt.Errorf("mismatch in token and JWKS algorithms")
		}

		// Check for expiry
		if time.Now().Before(cached.(*cacheValue).Expiry) {
			cert := cached.(*cacheValue)
			return cert.Data, nil
		}
	}

	resp, err := http.DefaultClient.Get(x.url)
	if err != nil {
		return nil, fmt.Errorf("json validation error: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("json validation error: %s", err)
	}
	defer resp.Body.Close()

	jwks := &jose.JSONWebKeySet{}
	err = json.Unmarshal(body, jwks)
	if err != nil {
		return nil, fmt.Errorf("json validation error c: %s", err)
	}

	// Get desired key from JWKS
	kid := token.Header["kid"].(string)
	key := jwks.Key(kid)[0]
	if !key.Valid() {
		return nil, fmt.Errorf("invalid JWKS key")
	}

	// Check for alg
	alg := key.Algorithm
	if alg != token.Header["alg"] {
		return nil, fmt.Errorf("mismatch in token and JWKS algorithms")
	}

	// pk := key.Certificates[0].PublicKey
	pk := key.Key

	// Store value in cache
	exp := claims["exp"].(float64)
	lruCache.Add(cacheKey, &cacheValue{
		Algorithm: alg,
		Data:      pk,
		Expiry:    time.Unix(int64(exp), 0),
	})

	return pk, nil
}
