package claims

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12/context"
	"strconv"
)

func GetClaimInt64(claimName string, ctx context.Context) int64 {
	str := GetClaimString(claimName, ctx)
	r, _ := strconv.ParseInt(str, 10, 64)
	return r
}

func GetClaimString(claimName string, ctx context.Context) string {
	j := ctx.Values().Get("jwt")
	if j != nil {
		if token, ok := j.(*jwt.Token); ok {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if str, ok := claims[claimName].(string); ok && str != "" {
					return str
				}
			}
		}
	}

	return ""
}
