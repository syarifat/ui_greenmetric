package middleware

import (
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"

	"github.com/goravel/framework/contracts/http"
)

type JwtMiddleware struct{}

func (m *JwtMiddleware) Signature() string {
	return "jwt"
}

func (m *JwtMiddleware) Handle(ctx http.Context) {
	var user models.User
	err := facades.Auth(ctx).User(&user)
	if err != nil {
		ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized: Invalid or expired token",
		})
		return
	}

	ctx.Request().Next()
}
