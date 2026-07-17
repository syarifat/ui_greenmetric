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
	token := ctx.Request().Header("Authorization")
	if token == "" {
		ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Token tidak valid atau kedaluwarsa",
		}).Abort()
		return
	}

	// Parse will verify token signature, expiration, and store it in context
	_, err := facades.Auth(ctx).Parse(token)
	if err != nil {
		ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Token tidak valid atau kedaluwarsa",
		}).Abort()
		return
	}

	var user models.User
	err = facades.Auth(ctx).User(&user)
	if err != nil {
		ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Token tidak valid atau kedaluwarsa",
		}).Abort()
		return
	}

	ctx.Request().Next()
}
