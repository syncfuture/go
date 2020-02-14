package sgrpc

import (
	context "context"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/dgrijalva/jwt-go"
	u "github.com/syncfuture/go/util"
	"google.golang.org/grpc"
)

func AttachJWTToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	token := extractToken(ctx)
	if token == "" {
		return handler(ctx, req)
	}

	jwtToken, err := new(jwt.Parser).Parse(token, nil)
	vErr, ok := err.(*jwt.ValidationError)
	if err != nil && ok && vErr.Errors != jwt.ValidationErrorUnverifiable {
		u.LogError(err)
		return handler(ctx, req)
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok {
		return handler(context.WithValue(ctx, "user", claims), req)
	} else {
		return handler(ctx, req)
	}
}

// func extractToken(ctx context.Context) (string, error) {
// 	tokenHeader, err := extractHeader(ctx, "authorization")
// 	if u.LogError(err) {
// 		return "", err
// 	}
// 	array := strings.Split(tokenHeader, " ")
// 	if len(array) != 2 {
// 		return "", status.Error(codes.Unauthenticated, "invalid token in request")
// 	}
// 	return array[1], nil
// }

// func extractHeader(ctx context.Context, header string) (string, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return "", status.Error(codes.Unauthenticated, "no headers in request")
// 	}

// 	authHeaders, ok := md[header]
// 	if !ok {
// 		return "", status.Error(codes.Unauthenticated, "no header in request")
// 	}

// 	if len(authHeaders) != 1 {
// 		return "", status.Error(codes.Unauthenticated, "more than 1 header in request")
// 	}

// 	return authHeaders[0], nil
// }

func extractToken(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	if authStrs, ok := md["authorization"]; ok {
		if len(authStrs) == 1 {
			array := strings.Split(authStrs[0], " ")
			if len(array) == 2 {
				return array[1]
			}
		}
	}
	return ""
}
