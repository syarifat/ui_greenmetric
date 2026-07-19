package middleware

import (
	"strings"
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"

	"github.com/goravel/framework/contracts/http"
)

type RbacMiddleware struct{}

func (m *RbacMiddleware) Signature() string {
	return "rbac"
}

func (m *RbacMiddleware) Handle(ctx http.Context) {
	var user models.User
	err := facades.Auth(ctx).User(&user)
	if err != nil {
		ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Token tidak valid atau kedaluwarsa",
		}).Abort()
		return
	}

	// 1. SUPER_ADMIN has access to everything
	if user.Role == "SUPER_ADMIN" {
		ctx.Request().Next()
		return
	}

	// 2. Check category-specific restrictions (SI, EC, WS, WR, TR, ED, GD)
	categoryCode := ctx.Request().Route("category_code")
	if categoryCode != "" {
		categoryCode = strings.ToUpper(categoryCode)
		// ADMIN_KAMPUS has access to all category forms
		if user.Role == "ADMIN_KAMPUS" {
			ctx.Request().Next()
			return
		}

		// Operator role pattern: OPERATOR_<CATEGORY_CODE>
		expectedOperatorRole := "OPERATOR_" + categoryCode
		if user.Role != expectedOperatorRole {
			ctx.Response().Json(http.StatusForbidden, http.Json{
				"status":  "error",
				"code":    http.StatusForbidden,
				"message": "Anda tidak memiliki hak akses",
			}).Abort()
			return
		}

		ctx.Request().Next()
		return
	}

	// 3. Check management actions (Campuses: SUPER_ADMIN only, Users: ADMIN_KAMPUS/SUPER_ADMIN)
	path := ctx.Request().Path()
	if strings.Contains(path, "/campuses") || strings.Contains(path, "/admin/indicators") {
		// Only SUPER_ADMIN can manage campuses and indicators
		ctx.Response().Json(http.StatusForbidden, http.Json{
			"status":  "error",
			"code":    http.StatusForbidden,
			"message": "Anda tidak memiliki hak akses",
		}).Abort()
		return
	}

	if strings.Contains(path, "/users") {
		// ADMIN_KAMPUS can manage users of their own campus
		if user.Role == "ADMIN_KAMPUS" {
			ctx.Request().Next()
			return
		}

		ctx.Response().Json(http.StatusForbidden, http.Json{
			"status":  "error",
			"code":    http.StatusForbidden,
			"message": "Anda tidak memiliki hak akses",
		}).Abort()
		return
	}

	ctx.Request().Next()
}
