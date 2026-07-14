package controllers

import (
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"

	"github.com/goravel/framework/contracts/http"
)

type CampusController struct{}

func NewCampusController() *CampusController {
	return &CampusController{}
}

// Index lists all campuses
func (r *CampusController) Index(ctx http.Context) http.Response {
	var campuses []models.Campus
	err := facades.Orm().Query().Get(&campuses)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to fetch campuses",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Campuses retrieved successfully",
		"data":    campuses,
	})
}

// Store creates a new campus
func (r *CampusController) Store(ctx http.Context) http.Response {
	code := ctx.Request().Input("code")
	name := ctx.Request().Input("name")
	instType := ctx.Request().Input("institution_type")
	climate := ctx.Request().Input("climate")
	setting := ctx.Request().Input("setting")

	if code == "" || name == "" {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "Code and Name are required",
		})
	}

	// Check unique code
	count, err := facades.Orm().Query().Model(&models.Campus{}).Where("code = ?", code).Count()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Database error",
		})
	}
	if count > 0 {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "Campus code already exists",
		})
	}

	campus := models.Campus{
		Code:            code,
		Name:            name,
		InstitutionType: instType,
		Climate:         climate,
		Setting:         setting,
	}

	if err := facades.Orm().Query().Create(&campus); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to create campus",
		})
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{
		"status":  "success",
		"message": "Campus created successfully",
		"data":    campus,
	})
}

// Update modifies an existing campus
func (r *CampusController) Update(ctx http.Context) http.Response {
	id := ctx.Request().RouteInt("id")
	var campus models.Campus
	err := facades.Orm().Query().Where("id = ?", id).First(&campus)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "Campus not found",
		})
	}

	code := ctx.Request().Input("code")
	name := ctx.Request().Input("name")
	instType := ctx.Request().Input("institution_type")
	climate := ctx.Request().Input("climate")
	setting := ctx.Request().Input("setting")

	if code != "" {
		// Check unique code (excluding current campus id)
		count, err := facades.Orm().Query().Model(&models.Campus{}).Where("code = ? AND id != ?", code, id).Count()
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Database error",
			})
		}
		if count > 0 {
			return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
				"status":  "error",
				"code":    http.StatusUnprocessableEntity,
				"message": "Campus code already exists",
			})
		}
		campus.Code = code
	}

	if name != "" {
		campus.Name = name
	}
	if instType != "" {
		campus.InstitutionType = instType
	}
	if climate != "" {
		campus.Climate = climate
	}
	if setting != "" {
		campus.Setting = setting
	}

	if err := facades.Orm().Query().Save(&campus); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to update campus",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Campus updated successfully",
		"data":    campus,
	})
}

// Destroy deletes an existing campus
func (r *CampusController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().RouteInt("id")
	var campus models.Campus
	err := facades.Orm().Query().Where("id = ?", id).First(&campus)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "Campus not found",
		})
	}

	if _, err := facades.Orm().Query().Delete(&campus); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to delete campus",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Campus deleted successfully",
	})
}	
