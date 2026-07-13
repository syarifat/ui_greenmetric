package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ui_greenmetric/app/facades"
)

type M20260713155204CreateAssessmentAnswersTable struct{}

// Signature The unique signature for the migration.
func (r *M20260713155204CreateAssessmentAnswersTable) Signature() string {
	return "20260713155204_create_assessment_answers_table"
}

// Up Run the migrations.
func (r *M20260713155204CreateAssessmentAnswersTable) Up() error {
	if !facades.Schema().HasTable("assessment_answers") {
		return facades.Schema().Create("assessment_answers", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("campus_assessment_id")
			table.UnsignedBigInteger("indicator_id")
			table.Json("raw_input_data").Nullable()
			table.Decimal("calculated_value").Total(15).Places(4).Nullable()
			table.UnsignedBigInteger("selected_tier_id").Nullable()
			table.Decimal("earned_points").Total(10).Places(2).Default(0)
			table.Timestamps()

			table.Foreign("campus_assessment_id").References("id").On("campus_assessments").CascadeOnDelete()
			table.Foreign("indicator_id").References("id").On("indicators").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260713155204CreateAssessmentAnswersTable) Down() error {
	return facades.Schema().DropIfExists("assessment_answers")
}
