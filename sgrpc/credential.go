package sgrpc

import (
	context "context"

	"google.golang.org/grpc/credentials"
)

const (
	_authHeader = "authorization"
	_tokenType  = "Bearer "
)

type tokenCredential struct {
	Token      string
	RequireTLS bool
}

// GetRequestMetadata 获取请求Meta
func (x tokenCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if x.Token != "" {
		return map[string]string{
			_authHeader: _tokenType + x.Token,
		}, nil
	} else {
		return map[string]string{}, nil
	}
}

// RequireTransportSecurity 是否开启TLS
func (x tokenCredential) RequireTransportSecurity() bool {
	return x.RequireTLS
}

func NewTokenCredential(token string, requireTLS bool) credentials.PerRPCCredentials {
	return &tokenCredential{
		Token:      token,
		RequireTLS: requireTLS,
	}
}
