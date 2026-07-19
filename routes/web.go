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
	// Register Global NotFound Middleware to handle 404 cleanly
	facades.Route().GlobalMiddleware(&middleware.NotFoundMiddleware{})

	// Global 500 Panic Recovery (must be registered at the beginning to avoid resetting router and wiping routes)
	facades.Route().Recover(func(ctx http.Context, err any) {
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Terjadi kesalahan internal pada server",
		})
	})

	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version": support.Version,
		})
	})

	facades.Route().Static("public", "./public")

	logController := controllers.NewLogController()
	facades.Route().Get("/logs", logController.ViewLogs)

	authController := controllers.NewAuthController()

	// API V1 Group
	facades.Route().Prefix("/api/v1").Group(func(router route.Router) {
		// Public Auth Routes
		router.Post("/auth/login", authController.Login)

		// Protected Routes (JWT & RBAC Middleware)
		router.Middleware(&middleware.JwtMiddleware{}, &middleware.RbacMiddleware{}).Group(func(protected route.Router) {
			protected.Post("/auth/logout", authController.Logout)

			// Campus Management (SUPER_ADMIN)
			campusController := controllers.NewCampusController()
			protected.Get("/campuses", campusController.Index)
			protected.Post("/campuses", campusController.Store)
			protected.Put("/campuses/{id}", campusController.Update)
			protected.Delete("/campuses/{id}", campusController.Destroy)

			// Indicator Management (SUPER_ADMIN)
			adminIndicatorController := controllers.NewAdminIndicatorController()
			protected.Get("/admin/indicators", adminIndicatorController.Index)
			protected.Post("/admin/indicators", adminIndicatorController.Store)
			protected.Put("/admin/indicators/{id}", adminIndicatorController.Update)
			protected.Delete("/admin/indicators/{id}", adminIndicatorController.Destroy)

			// User Management (ADMIN_KAMPUS / SUPER_ADMIN)
			userController := controllers.NewUserController()
			protected.Get("/users", userController.Index)
			protected.Post("/users", userController.Store)
			protected.Put("/users/{id}", userController.Update)
			protected.Delete("/users/{id}", userController.Destroy)

			// Dashboard API
			dashboardController := controllers.NewDashboardController()
			protected.Get("/assessments/dashboard", dashboardController.Index)

			// Assessment API (Core Scoring Engine)
			assessmentController := controllers.NewAssessmentController()
			protected.Get("/categories/{category_code}/indicators", assessmentController.GetIndicatorsByCategory)
			protected.Post("/assessments/answers", assessmentController.SaveAnswer)
			protected.Post("/assessments/submit", assessmentController.SubmitAssessment)

			// Evidence Upload API
			evidenceController := controllers.NewEvidenceController()
			protected.Post("/evidences/upload", evidenceController.Upload)
			protected.Delete("/evidences/{id}", evidenceController.Destroy)
		})
	})
}
