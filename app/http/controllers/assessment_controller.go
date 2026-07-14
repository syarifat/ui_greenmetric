package controllers

import (
	"strings"
	"time"
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/http/requests"
	"ui_greenmetric/app/models"
	"ui_greenmetric/app/services"

	"github.com/goravel/framework/contracts/http"
)

type AssessmentController struct {
	scoringService *services.ScoringService
}

func NewAssessmentController() *AssessmentController {
	return &AssessmentController{
		scoringService: services.NewScoringService(),
	}
}

// GetIndicatorsByCategory fetches the questions and previous answers for a specific category
func (r *AssessmentController) GetIndicatorsByCategory(ctx http.Context) http.Response {
	categoryCode := strings.ToUpper(ctx.Request().Route("category_code"))

	var currentUser models.User
	if err := facades.Auth(ctx).User(&currentUser); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	// 1. Fetch Category
	var category models.Category
	err := facades.Orm().Query().Where("code = ?", categoryCode).First(&category)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "Category not found",
		})
	}

	// 2. Fetch Indicators in Category
	var indicators []models.Indicator
	err = facades.Orm().Query().Where("category_id = ?", category.ID).Get(&indicators)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to fetch indicators",
		})
	}

	// 3. Fetch Tiers for all these indicators
	var indicatorIDs []uint
	for _, ind := range indicators {
		indicatorIDs = append(indicatorIDs, ind.ID)
	}

	var tiers []models.IndicatorScoringTier
	if len(indicatorIDs) > 0 {
		facades.Orm().Query().Where("indicator_id IN ?", indicatorIDs).Get(&tiers)
	}

	// Map tiers to indicator IDs
	tiersMap := make(map[uint][]models.IndicatorScoringTier)
	for _, t := range tiers {
		tiersMap[t.IndicatorID] = append(tiersMap[t.IndicatorID], t)
	}

	// 4. Fetch Current Assessment and Answers
	currentYear := time.Now().Year()
	var assessment models.CampusAssessment
	facades.Orm().Query().Where("campus_id = ? AND assessment_year = ?", currentUser.CampusID, currentYear).First(&assessment)

	answersMap := make(map[uint]models.AssessmentAnswer)
	if assessment.ID != 0 && len(indicatorIDs) > 0 {
		var answers []models.AssessmentAnswer
		facades.Orm().Query().With("Evidences").Where("campus_assessment_id = ? AND indicator_id IN ?", assessment.ID, indicatorIDs).Get(&answers)
		for _, ans := range answers {
			answersMap[ans.IndicatorID] = ans
		}
	}

	var resultList []http.Json
	for _, ind := range indicators {
		ans, hasAnswer := answersMap[ind.ID]
		var answerData any = nil
		if hasAnswer {
			answerData = http.Json{
				"id":               ans.ID,
				"raw_input_data":   ans.RawInputData,
				"calculated_value": ans.CalculatedValue,
				"selected_tier_id": ans.SelectedTierID,
				"earned_points":    ans.EarnedPoints,
				"evidences":        ans.Evidences,
			}
		}

		indTiers := tiersMap[ind.ID]
		if indTiers == nil {
			indTiers = []models.IndicatorScoringTier{}
		}

		resultList = append(resultList, http.Json{
			"id":         ind.ID,
			"code":       ind.Code,
			"title":      ind.Title,
			"input_type": ind.InputType,
			"max_points": ind.MaxPoints,
			"tiers":      indTiers,
			"answer":     answerData,
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Indicators and answers loaded successfully",
		"data": http.Json{
			"category_code": category.Code,
			"category_name": category.Name,
			"indicators":    resultList,
		},
	})
}

// SaveAnswer stores/calculates an answer for an indicator
func (r *AssessmentController) SaveAnswer(ctx http.Context) http.Response {
	var currentUser models.User
	if err := facades.Auth(ctx).User(&currentUser); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	var saveRequest requests.SaveAnswerRequest
	errors, err := ctx.Request().ValidateRequest(&saveRequest)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "Validation failed",
			"errors":  errors.All(),
		})
	}

	// 1. Get or Create CampusAssessment for the year
	var assessment models.CampusAssessment
	err = facades.Orm().Query().Where("campus_id = ? AND assessment_year = ?", currentUser.CampusID, saveRequest.AssessmentYear).First(&assessment)
	if err != nil {
		// Not found, create one
		assessment = models.CampusAssessment{
			CampusID:       currentUser.CampusID,
			AssessmentYear: saveRequest.AssessmentYear,
			OverallScore:   0.0,
			Status:         "DRAFT",
		}
		if err := facades.Orm().Query().Create(&assessment); err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to create campus assessment",
			})
		}
	}

	// Check if assessment is locked (already SUBMITTED or VERIFIED)
	if assessment.Status != "DRAFT" {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"status":  "error",
			"code":    http.StatusForbidden,
			"message": "Forbidden: Assessment is locked and cannot be modified",
		})
	}

	// 2. Perform calculation and save answer
	ans, err := r.scoringService.Calculate(assessment.ID, saveRequest.IndicatorCode, saveRequest.RawInputData)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": "Calculation error: " + err.Error(),
		})
	}

	// Fetch updated assessment to get the new overall score
	facades.Orm().Query().Where("id = ?", assessment.ID).First(&assessment)

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Answer saved and calculated successfully",
		"data": http.Json{
			"earned_points":    ans.EarnedPoints,
			"calculated_value": ans.CalculatedValue,
			"overall_score":    assessment.OverallScore,
		},
	})
}

// SubmitAssessment finalizes and locks the assessment for a specific year
func (r *AssessmentController) SubmitAssessment(ctx http.Context) http.Response {
	var currentUser models.User
	if err := facades.Auth(ctx).User(&currentUser); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	// Validate role: must be ADMIN_KAMPUS or SUPER_ADMIN
	if currentUser.Role != "ADMIN_KAMPUS" && currentUser.Role != "SUPER_ADMIN" {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"status":  "error",
			"code":    http.StatusForbidden,
			"message": "Forbidden: Only campus admins can submit and finalize assessments",
		})
	}

	year := ctx.Request().InputInt("assessment_year")
	if year == 0 {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "assessment_year is required",
		})
	}

	// Get Campus ID
	var campusID uint
	if currentUser.Role == "SUPER_ADMIN" {
		campusID = uint(ctx.Request().InputInt("campus_id"))
		if campusID == 0 {
			return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
				"status":  "error",
				"code":    http.StatusUnprocessableEntity,
				"message": "campus_id is required for SUPER_ADMIN",
			})
		}
	} else {
		campusID = currentUser.CampusID
	}

	var assessment models.CampusAssessment
	err := facades.Orm().Query().Where("campus_id = ? AND assessment_year = ?", campusID, year).First(&assessment)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "Campus assessment not found for the specified year",
		})
	}

	// Check if already submitted
	if assessment.Status == "SUBMITTED" || assessment.Status == "VERIFIED" {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "Assessment has already been submitted and locked",
		})
	}

	assessment.Status = "SUBMITTED"
	if err := facades.Orm().Query().Save(&assessment); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to finalize assessment",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Assessment finalized and locked successfully",
		"data": http.Json{
			"campus_assessment_id": assessment.ID,
			"status":               assessment.Status,
			"overall_score":        assessment.OverallScore,
		},
	})
}	
