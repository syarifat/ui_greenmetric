package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"
)

type ScoringService struct{}

func NewScoringService() *ScoringService {
	return &ScoringService{}
}

// Calculate computes scores and points for an assessment answer
func (s *ScoringService) Calculate(assessmentID uint, indicatorCode string, rawInput map[string]any) (models.AssessmentAnswer, error) {
	var answer models.AssessmentAnswer

	// 1. Fetch Indicator
	var indicator models.Indicator
	err := facades.Orm().Query().Where("code = ?", indicatorCode).First(&indicator)
	if err != nil {
		return answer, fmt.Errorf("indicator not found: %v", err)
	}

	// Fetch existing answer if any, or initialize a new one
	err = facades.Orm().Query().Where("campus_assessment_id = ? AND indicator_id = ?", assessmentID, indicator.ID).First(&answer)
	if err != nil {
		answer = models.AssessmentAnswer{
			CampusAssessmentID: assessmentID,
			IndicatorID:        indicator.ID,
		}
	}

	// Save raw input data as JSON string in model
	jsonBytes, err := json.Marshal(rawInput)
	if err != nil {
		return answer, fmt.Errorf("failed to marshal raw input: %v", err)
	}
	answer.RawInputData = string(jsonBytes)

	// Fetch Tiers for this indicator
	var tiers []models.IndicatorScoringTier
	err = facades.Orm().Query().Where("indicator_id = ?", indicator.ID).Get(&tiers)
	if err != nil {
		return answer, fmt.Errorf("failed to fetch scoring tiers: %v", err)
	}

	var matchedTier *models.IndicatorScoringTier

	if indicator.InputType == "SINGLE_CHOICE" {
		// Single Choice uses selected option label
		optionLabelVal, ok := rawInput["option_label"]
		if !ok {
			return answer, errors.New("option_label is required for SINGLE_CHOICE indicator")
		}
		optionLabel, ok := optionLabelVal.(string)
		if !ok {
			return answer, errors.New("option_label must be a string")
		}

		for i := range tiers {
			if tiers[i].OptionLabel == optionLabel {
				matchedTier = &tiers[i]
				break
			}
		}

		if matchedTier == nil {
			return answer, fmt.Errorf("matching tier not found for option: %s", optionLabel)
		}

		answer.CalculatedValue = nil
		answer.SelectedTierID = &matchedTier.ID
		answer.EarnedPoints = matchedTier.PointMultiplier * float64(indicator.MaxPoints)

	} else {
		// Numeric Formula requires mathematical calculation
		calculatedVal, err := s.EvaluateFormula(indicator.Code, rawInput)
		if err != nil {
			return answer, err
		}
		answer.CalculatedValue = &calculatedVal

		// Find matching tier based on value and operators
		for i := range tiers {
			t := &tiers[i]
			if t.Operator == "<=" && t.MaxValue != nil && calculatedVal <= *t.MaxValue {
				matchedTier = t
				break
			} else if t.Operator == "<" && t.MaxValue != nil && calculatedVal < *t.MaxValue {
				matchedTier = t
				break
			} else if t.Operator == ">=" && t.MinValue != nil && calculatedVal >= *t.MinValue {
				matchedTier = t
				break
			} else if t.Operator == ">" && t.MinValue != nil && calculatedVal > *t.MinValue {
				matchedTier = t
				break
			} else if t.Operator == "BETWEEN" && t.MinValue != nil && t.MaxValue != nil {
				if calculatedVal > *t.MinValue && calculatedVal <= *t.MaxValue {
					matchedTier = t
					break
				}
			}
		}

		// Fallback to default/lowest tier multiplier if none match (or use 0)
		multiplier := 0.0
		if matchedTier != nil {
			multiplier = matchedTier.PointMultiplier
			answer.SelectedTierID = &matchedTier.ID
		} else {
			answer.SelectedTierID = nil
		}
		answer.EarnedPoints = multiplier * float64(indicator.MaxPoints)
	}

	// 2. Save the assessment answer
	if answer.ID == 0 {
		err = facades.Orm().Query().Create(&answer)
	} else {
		err = facades.Orm().Query().Save(&answer)
	}
	if err != nil {
		return answer, fmt.Errorf("failed to save answer: %v", err)
	}

	// 3. Recalculate and update the CampusAssessment overall score
	err = s.UpdateOverallScore(assessmentID)
	if err != nil {
		return answer, fmt.Errorf("failed to update overall score: %v", err)
	}

	return answer, nil
}

// UpdateOverallScore sums all earned points and updates the campus assessment
func (s *ScoringService) UpdateOverallScore(assessmentID uint) error {
	var assessment models.CampusAssessment
	err := facades.Orm().Query().Where("id = ?", assessmentID).First(&assessment)
	if err != nil {
		return err
	}

	var answers []models.AssessmentAnswer
	err = facades.Orm().Query().Where("campus_assessment_id = ?", assessmentID).Get(&answers)
	if err != nil {
		return err
	}

	totalScore := 0.0
	for _, ans := range answers {
		totalScore += ans.EarnedPoints
	}

	assessment.OverallScore = totalScore
	return facades.Orm().Query().Save(&assessment)
}

// EvaluateFormula implements mathematics calculations based on indicator code
func (s *ScoringService) EvaluateFormula(indicatorCode string, rawData map[string]any) (float64, error) {
	getFloat := func(key string) float64 {
		val, ok := rawData[key]
		if !ok {
			return 0.0
		}
		switch v := val.(type) {
		case float64:
			return v
		case float32:
			return float64(v)
		case int:
			return float64(v)
		case int64:
			return float64(v)
		default:
			return 0.0
		}
	}

	switch indicatorCode {
	case "SI1":
		total := getFloat("luas_total")
		dasar := getFloat("luas_dasar")
		if total == 0 {
			return 0, errors.New("luas_total cannot be zero")
		}
		return ((total - dasar) / total) * 100, nil

	case "SI2":
		hutan := getFloat("luas_hutan")
		total := getFloat("luas_total")
		if total == 0 {
			return 0, errors.New("luas_total cannot be zero")
		}
		return (hutan / total) * 100, nil

	case "SI3":
		tanam := getFloat("luas_vegetasi")
		total := getFloat("luas_total")
		if total == 0 {
			return 0, errors.New("luas_total cannot be zero")
		}
		return (tanam / total) * 100, nil

	case "SI4":
		total := getFloat("luas_total")
		dasar := getFloat("luas_dasar")
		populasi := getFloat("populasi")
		if populasi == 0 {
			return 0, errors.New("populasi cannot be zero")
		}
		return (total - dasar) / populasi, nil

	case "EC1":
		return getFloat("persentase_alat_hemat_energi"), nil

	case "EC2":
		smart := getFloat("luas_smart_building")
		totalBangunan := getFloat("luas_total_bangunan")
		if totalBangunan == 0 {
			return 0, errors.New("luas_total_bangunan cannot be zero")
		}
		return (smart / totalBangunan) * 100, nil

	case "EC4":
		listrik := getFloat("total_listrik")
		populasi := getFloat("populasi")
		if populasi == 0 {
			return 0, errors.New("populasi cannot be zero")
		}
		return listrik / populasi, nil

	case "EC5":
		produksi := getFloat("produksi_energi_terbarukan")
		total := getFloat("total_penggunaan_energi")
		if total == 0 {
			return 0, errors.New("total_penggunaan_energi cannot be zero")
		}
		return (produksi / total) * 100, nil

	case "EC8":
		jejak := getFloat("jejak_karbon")
		populasi := getFloat("populasi")
		if populasi == 0 {
			return 0, errors.New("populasi cannot be zero")
		}
		return jejak / populasi, nil

	case "WR1":
		resapan := getFloat("area_resapan")
		total := getFloat("luas_total")
		if total == 0 {
			return 0, errors.New("luas_total cannot be zero")
		}
		return (resapan / total) * 100, nil

	case "WR5":
		olahan := getFloat("air_olahan_dikonsumsi")
		total := getFloat("total_air_dikonsumsi")
		if total == 0 {
			return 0, errors.New("total_air_dikonsumsi cannot be zero")
		}
		return (olahan / total) * 100, nil

	case "TR1":
		kendaraan := getFloat("total_kendaraan")
		populasi := getFloat("populasi")
		if populasi == 0 {
			return 0, errors.New("populasi cannot be zero")
		}
		return kendaraan / populasi, nil

	case "TR4":
		zev := getFloat("total_zev")
		populasi := getFloat("populasi")
		if populasi == 0 {
			return 0, errors.New("populasi cannot be zero")
		}
		return zev / populasi, nil

	case "TR5":
		parkir := getFloat("luas_parkir")
		total := getFloat("luas_total")
		if total == 0 {
			return 0, errors.New("luas_total cannot be zero")
		}
		return (parkir / total) * 100, nil

	case "ED1":
		mk := getFloat("mk_keberlanjutan")
		total := getFloat("total_mk")
		if total == 0 {
			return 0, errors.New("total_mk cannot be zero")
		}
		return (mk / total) * 100, nil

	case "ED2":
		dana := getFloat("dana_riset_keberlanjutan")
		total := getFloat("total_dana_riset")
		if total == 0 {
			return 0, errors.New("total_dana_riset cannot be zero")
		}
		return (dana / total) * 100, nil

	case "ED3":
		pub := getFloat("publikasi_keberlanjutan")
		total := getFloat("total_publikasi")
		if total == 0 {
			return 0, errors.New("total_publikasi cannot be zero")
		}
		return (pub / total) * 100, nil

	case "ED10":
		green := getFloat("lulusan_green_jobs")
		total := getFloat("total_lulusan")
		if total == 0 {
			return 0, errors.New("total_lulusan cannot be zero")
		}
		return (green / total) * 100, nil

	case "GD1":
		anggaran := getFloat("anggaran_keberlanjutan")
		total := getFloat("total_anggaran")
		if total == 0 {
			return 0, errors.New("total_anggaran cannot be zero")
		}
		return (anggaran / total) * 100, nil

	case "GD8":
		perempuan := getFloat("pimpinan_perempuan")
		total := getFloat("total_pimpinan")
		if total == 0 {
			return 0, errors.New("total_pimpinan cannot be zero")
		}
		return (perempuan / total) * 100, nil

	default:
		return getFloat("value"), nil
	}
}
