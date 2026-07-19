package controllers

import (
	"fmt"
	"strings"
	"time"
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"

	"github.com/goravel/framework/contracts/http"
)

type EvidenceController struct{}

func NewEvidenceController() *EvidenceController {
	return &EvidenceController{}
}

// Upload handles uploading a new document evidence
func (r *EvidenceController) Upload(ctx http.Context) http.Response {
	var currentUser models.User
	if err := facades.Auth(ctx).User(&currentUser); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	// 1. Get and Validate File
	file, err := ctx.Request().File("file")
	if err != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "File is required",
		})
	}

	// Validate file size (max 2MB)
	size, err := file.Size()
	if err != nil || size > 2*1024*1024 {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "File size must not exceed 2MB",
		})
	}

	// Validate file extension
	ext := strings.ToLower(file.GetClientOriginalExtension())
	if ext != "pdf" && ext != "jpg" && ext != "jpeg" && ext != "png" {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "Only PDF, JPG, JPEG, and PNG files are allowed",
		})
	}

	// 2. Validate input and check ownership
	answerID := ctx.Request().InputInt("assessment_answer_id")
	description := ctx.Request().Input("description")

	if answerID == 0 {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "assessment_answer_id is required",
		})
	}

	var answer models.AssessmentAnswer
	err = facades.Orm().Query().Where("id = ?", answerID).FirstOrFail(&answer)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "Assessment answer not found",
		})
	}

	var assessment models.CampusAssessment
	err = facades.Orm().Query().Where("id = ?", answer.CampusAssessmentID).FirstOrFail(&assessment)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Associated campus assessment not found",
		})
	}

	// Verify multi-tenant ownership
	if currentUser.Role != "SUPER_ADMIN" && assessment.CampusID != currentUser.CampusID {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"status":  "error",
			"code":    http.StatusForbidden,
			"message": "Forbidden: You do not own this assessment",
		})
	}

	// Check if assessment is locked
	if assessment.Status != "DRAFT" {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"status":  "error",
			"code":    http.StatusForbidden,
			"message": "Forbidden: Assessment is locked and cannot be modified",
		})
	}

	// 3. Store file to public/evidences
	hashedName := file.HashName()
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), hashedName)

	storedPath, err := file.Disk("public_dir").StoreAs("evidences", filename)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to store file: " + err.Error(),
		})
	}

	// 4. Save evidence record to DB
	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	evidence := models.AssessmentEvidence{
		AssessmentAnswerID: uint(answerID),
		DocumentName:       file.GetClientOriginalName(),
		Description:        descPtr,
		FileUrl:            "/public/" + storedPath,
	}

	if err := facades.Orm().Query().Create(&evidence); err != nil {
		// Clean up stored file if DB insert fails
		facades.Storage().Disk("public_dir").Delete(storedPath)

		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to save evidence to database",
		})
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{
		"status":  "success",
		"message": "Evidence document uploaded successfully",
		"data":    evidence,
	})
}

// Destroy deletes an uploaded evidence document
func (r *EvidenceController) Destroy(ctx http.Context) http.Response {
	var currentUser models.User
	if err := facades.Auth(ctx).User(&currentUser); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	evidenceID := ctx.Request().RouteInt("id")
	var evidence models.AssessmentEvidence
	err := facades.Orm().Query().Where("id = ?", evidenceID).FirstOrFail(&evidence)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "Evidence not found",
		})
	}

	var answer models.AssessmentAnswer
	err = facades.Orm().Query().Where("id = ?", evidence.AssessmentAnswerID).FirstOrFail(&answer)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Assessment answer not found",
		})
	}

	var assessment models.CampusAssessment
	err = facades.Orm().Query().Where("id = ?", answer.CampusAssessmentID).FirstOrFail(&assessment)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Campus assessment not found",
		})
	}

	// Verify multi-tenant ownership
	if currentUser.Role != "SUPER_ADMIN" && assessment.CampusID != currentUser.CampusID {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"status":  "error",
			"code":    http.StatusForbidden,
			"message": "Forbidden: You do not own this assessment",
		})
	}

	// Check if assessment is locked
	if assessment.Status != "DRAFT" {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"status":  "error",
			"code":    http.StatusForbidden,
			"message": "Forbidden: Assessment is locked and cannot be modified",
		})
	}

	// Delete physical file
	// FileUrl format: "/public/evidences/filename.ext" -> storage path is "evidences/filename.ext"
	storagePath := strings.TrimPrefix(evidence.FileUrl, "/public/")
	if storagePath != "" {
		facades.Storage().Disk("public_dir").Delete(storagePath)
	}

	// Delete DB record
	if _, err := facades.Orm().Query().Delete(&evidence); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to delete evidence from database",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Evidence document deleted successfully",
	})
}	
