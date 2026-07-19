package controllers

import (
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"

	"github.com/goravel/framework/contracts/http"
)

type AdminIndicatorController struct{}

func NewAdminIndicatorController() *AdminIndicatorController {
	return &AdminIndicatorController{}
}

// Index lists all indicators with fields
func (r *AdminIndicatorController) Index(ctx http.Context) http.Response {
	var indicators []models.Indicator
	err := facades.Orm().Query().With("Fields").With("Category").Get(&indicators)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to fetch indicators: " + err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Indicators retrieved successfully",
		"data":    indicators,
	})
}

// Store creates a new indicator with dynamic fields
func (r *AdminIndicatorController) Store(ctx http.Context) http.Response {
	var req struct {
		CategoryID uint `json:"category_id"`
		Code       string `json:"code"`
		Title      string `json:"title"`
		InputType  string `json:"input_type"`
		MaxPoints  int    `json:"max_points"`
		Fields []struct {
			Key      string `json:"key"`
			Label    string `json:"label"`
			Type     string `json:"type"`
			Options  string `json:"options"`
			Required bool   `json:"required"`
		} `json:"fields"`
	}

	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": "Invalid request payload: " + err.Error(),
		})
	}

	if req.CategoryID == 0 || req.Code == "" || req.Title == "" || req.InputType == "" {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "category_id, code, title, and input_type are required",
		})
	}

	// Verify unique code
	count, err := facades.Orm().Query().Model(&models.Indicator{}).Where("code = ?", req.Code).Count()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Database error: " + err.Error(),
		})
	}
	if count > 0 {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "Indicator code already exists",
		})
	}

	// Create indicator
	indicator := models.Indicator{
		CategoryID: req.CategoryID,
		Code:       req.Code,
		Title:      req.Title,
		InputType:  req.InputType,
		MaxPoints:  req.MaxPoints,
	}
	if err := facades.Orm().Query().Create(&indicator); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to create indicator: " + err.Error(),
		})
	}

	// If SINGLE_CHOICE, always make sure we have the default choice field
	actualFields := req.Fields
	if req.InputType == "SINGLE_CHOICE" && len(actualFields) == 0 {
		actualFields = append(actualFields, struct {
			Key      string `json:"key"`
			Label    string `json:"label"`
			Type     string `json:"type"`
			Options  string `json:"options"`
			Required bool   `json:"required"`
		}{
			Key:      "option_label",
			Label:    "Pilih Opsi Kriteria",
			Type:     "choice",
			Required: true,
		})
	}

	// Save fields
	for _, f := range actualFields {
		newField := models.IndicatorField{
			IndicatorID: indicator.ID,
			Key:         f.Key,
			Label:       f.Label,
			Type:        f.Type,
			Required:    f.Required,
		}
		if f.Options != "" {
			opts := f.Options
			newField.Options = &opts
		}
		_ = facades.Orm().Query().Create(&newField)
	}

	// Load fields for final response
	_ = facades.Orm().Query().Load(&indicator, "Fields")

	return ctx.Response().Json(http.StatusCreated, http.Json{
		"status":  "success",
		"message": "Indicator and fields created successfully",
		"data":    indicator,
	})
}

// Update modifies an existing indicator and its fields
func (r *AdminIndicatorController) Update(ctx http.Context) http.Response {
	id := ctx.Request().RouteInt("id")
	var indicator models.Indicator
	err := facades.Orm().Query().Where("id = ?", id).FirstOrFail(&indicator)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "Indicator not found",
		})
	}

	title := ctx.Request().Input("title")
	inputType := ctx.Request().Input("input_type")
	maxPoints := ctx.Request().InputInt("max_points")

	if title != "" {
		indicator.Title = title
	}
	if inputType != "" {
		indicator.InputType = inputType
	}
	if maxPoints > 0 {
		indicator.MaxPoints = int(maxPoints)
	}

	if err := facades.Orm().Query().Save(&indicator); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to update indicator: " + err.Error(),
		})
	}

	// Handle sync fields
	var wrapper struct {
		Fields []struct {
			Key      string `json:"key"`
			Label    string `json:"label"`
			Type     string `json:"type"`
			Options  string `json:"options"`
			Required bool   `json:"required"`
		} `json:"fields"`
	}
	if err := ctx.Request().Bind(&wrapper); err == nil && len(wrapper.Fields) > 0 {
		// Drop existing fields first to rewrite
		_, _ = facades.Orm().Query().Model(&models.IndicatorField{}).Where("indicator_id = ?", indicator.ID).Delete()

		for _, f := range wrapper.Fields {
			newField := models.IndicatorField{
				IndicatorID: indicator.ID,
				Key:         f.Key,
				Label:       f.Label,
				Type:        f.Type,
				Required:    f.Required,
			}
			if f.Options != "" {
				opts := f.Options
				newField.Options = &opts
			}
			_ = facades.Orm().Query().Create(&newField)
		}
	}

	_ = facades.Orm().Query().Load(&indicator, "Fields")

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Indicator updated successfully",
		"data":    indicator,
	})
}

// Destroy deletes an indicator
func (r *AdminIndicatorController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().RouteInt("id")
	var indicator models.Indicator
	err := facades.Orm().Query().Where("id = ?", id).FirstOrFail(&indicator)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "Indicator not found",
		})
	}

	if _, err := facades.Orm().Query().Delete(&indicator); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to delete indicator: " + err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Indicator deleted successfully",
	})
}
