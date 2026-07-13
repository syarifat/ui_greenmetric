package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/support"

	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/http/controllers"
	"ui_greenmetric/app/http/middleware"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version": support.Version,
		})
	})

	facades.Route().Static("public", "./public")

	authController := controllers.NewAuthController()

	// API V1 Group
	facades.Route().Prefix("/api/v1").Group(func(router route.Router) {
		// Public Auth Routes
		router.Post("/auth/login", authController.Login)

		// Protected Routes (JWT & RBAC Middleware)
		router.Middleware(&middleware.JwtMiddleware{}, &middleware.RbacMiddleware{}).Group(func(protected route.Router) {
			protected.Post("/auth/logout", authController.Logout)
		})
	})
}
