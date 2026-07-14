package middleware

import (
	"github.com/goravel/framework/contracts/http"
)

type NotFoundMiddleware struct{}

func (m *NotFoundMiddleware) Signature() string {
	return "not_found"
}

func (m *NotFoundMiddleware) Handle(ctx http.Context) {
	ctx.Request().Next()

	if ctx.Request().OriginPath() == "" {
		ctx.Response().Json(404, http.Json{
			"status":  "error",
			"code":    404,
			"message": "Resource tidak ditemukan",
		})
	}
}
