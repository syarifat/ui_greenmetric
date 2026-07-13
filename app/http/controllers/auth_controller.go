package controllers

import (
	"ui_greenmetric/app/services"

	"github.com/goravel/framework/contracts/http"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

func (r *AuthController) Login(ctx http.Context) http.Response {
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")

	if email == "" || password == "" {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "Email and password are required",
		})
	}

	token, user, err := r.authService.Login(ctx, email, password)
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Login successful",
		"data": http.Json{
			"token": token,
			"user": http.Json{
				"id":        user.ID,
				"name":      user.Name,
				"email":     user.Email,
				"role":      user.Role,
				"campus_id": user.CampusID,
				"campus":    user.Campus,
			},
		},
	})
}

func (r *AuthController) Logout(ctx http.Context) http.Response {
	err := r.authService.Logout(ctx)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Logout failed",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Logout successful",
	})
}	
