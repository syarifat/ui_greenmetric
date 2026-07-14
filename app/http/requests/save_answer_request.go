package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type SaveAnswerRequest struct {
	IndicatorCode  string         `form:"indicator_code" json:"indicator_code"`
	AssessmentYear int            `form:"assessment_year" json:"assessment_year"`
	RawInputData   map[string]any `form:"raw_input_data" json:"raw_input_data"`
}

func (r *SaveAnswerRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *SaveAnswerRequest) Filters(ctx http.Context) map[string]any {
	return map[string]any{}
}

func (r *SaveAnswerRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"indicator_code":  "required",
		"assessment_year": "required",
		"raw_input_data":  "required",
	}
}

func (r *SaveAnswerRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"indicator_code.required":  "Indicator code is required",
		"assessment_year.required": "Assessment year is required",
		"raw_input_data.required":  "Raw input data is required",
	}
}

func (r *SaveAnswerRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}
