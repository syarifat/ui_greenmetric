package controllers

import (
	"time"
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"

	"github.com/goravel/framework/contracts/http"
)

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

// Index returns dashboard statistics for the authenticated operator's campus
func (r *DashboardController) Index(ctx http.Context) http.Response {
	var currentUser models.User
	if err := facades.Auth(ctx).User(&currentUser); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	// 1. Fetch Campus Info
	var campus models.Campus
	err := facades.Orm().Query().Where("id = ?", currentUser.CampusID).First(&campus)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Campus info not found",
		})
	}

	// 2. Fetch Current Assessment
	currentYear := time.Now().Year()
	var assessment models.CampusAssessment
	err = facades.Orm().Query().Where("campus_id = ? AND assessment_year = ?", currentUser.CampusID, currentYear).First(&assessment)
	if err != nil {
		// Default to a draft assessment if not created yet
		assessment = models.CampusAssessment{
			CampusID:       currentUser.CampusID,
			AssessmentYear: currentYear,
			OverallScore:   0.0,
			Status:         "DRAFT",
		}
	}

	// 3. Fetch all categories, indicators, and answers for in-memory sum
	var categories []models.Category
	facades.Orm().Query().Get(&categories)

	var indicators []models.Indicator
	facades.Orm().Query().Get(&indicators)

	var answers []models.AssessmentAnswer
	if assessment.ID != 0 {
		facades.Orm().Query().Where("campus_assessment_id = ?", assessment.ID).Get(&answers)
	}

	// Create helper maps for calculation
	indicatorToCategory := make(map[uint]uint)
	for _, ind := range indicators {
		indicatorToCategory[ind.ID] = ind.CategoryID
	}

	categoryEarned := make(map[uint]float64)
	for _, ans := range answers {
		catID := indicatorToCategory[ans.IndicatorID]
		categoryEarned[catID] += ans.EarnedPoints
	}

	// Build category breakdown list
	var categoryBreakdown []http.Json
	for _, cat := range categories {
		earned := categoryEarned[cat.ID]
		categoryBreakdown = append(categoryBreakdown, http.Json{
			"category_code":     cat.Code,
			"category_name":     cat.Name,
			"earned_points":     earned,
			"max_points":        cat.MaxPoints,
			"weight_percentage": cat.WeightPercentage,
		})
	}

	// 4. Calculate estimated rank (among all campuses for current year)
	var allAssessments []models.CampusAssessment
	facades.Orm().Query().Where("assessment_year = ?", currentYear).Order("overall_score desc").Get(&allAssessments)

	rank := 1
	hasAssessmentRanked := false
	for i, assess := range allAssessments {
		if assess.CampusID == currentUser.CampusID {
			rank = i + 1
			hasAssessmentRanked = true
			break
		}
	}
	// If the current campus doesn't have an assessment in DB yet, put them at the end of the list
	if !hasAssessmentRanked {
		rank = len(allAssessments) + 1
	}

	// 5. Fetch Trend History
	var historicalAssessments []models.CampusAssessment
	facades.Orm().Query().Where("campus_id = ?", currentUser.CampusID).Order("assessment_year asc").Get(&historicalAssessments)

	var trendHistory []http.Json
	for _, hist := range historicalAssessments {
		trendHistory = append(trendHistory, http.Json{
			"year":  hist.AssessmentYear,
			"score": hist.OverallScore,
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Dashboard statistics loaded successfully",
		"data": http.Json{
			"campus_name":       campus.Name,
			"current_year":      currentYear,
			"assessment_status": assessment.Status,
			"overall_score":     assessment.OverallScore,
			"max_overall_score": 10000,
			"estimated_rank":    rank,
			"category_breakdown": categoryBreakdown,
			"trend_history":     trendHistory,
		},
	})
}	
