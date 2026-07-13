package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ui_greenmetric/app/facades"
)

type M20260713155211CreateAssessmentEvidencesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260713155211CreateAssessmentEvidencesTable) Signature() string {
	return "20260713155211_create_assessment_evidences_table"
}

// Up Run the migrations.
func (r *M20260713155211CreateAssessmentEvidencesTable) Up() error {
	if !facades.Schema().HasTable("assessment_evidences") {
		return facades.Schema().Create("assessment_evidences", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("assessment_answer_id")
			table.String("document_name")
			table.Text("description").Nullable()
			table.String("file_url", 500)
			table.Timestamps()

			table.Foreign("assessment_answer_id").References("id").On("assessment_answers").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260713155211CreateAssessmentEvidencesTable) Down() error {
	return facades.Schema().DropIfExists("assessment_evidences")
}
