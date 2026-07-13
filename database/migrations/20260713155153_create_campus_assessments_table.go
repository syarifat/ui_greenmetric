package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ui_greenmetric/app/facades"
)

type M20260713155153CreateCampusAssessmentsTable struct{}

// Signature The unique signature for the migration.
func (r *M20260713155153CreateCampusAssessmentsTable) Signature() string {
	return "20260713155153_create_campus_assessments_table"
}

// Up Run the migrations.
func (r *M20260713155153CreateCampusAssessmentsTable) Up() error {
	if !facades.Schema().HasTable("campus_assessments") {
		return facades.Schema().Create("campus_assessments", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("campus_id")
			table.Integer("assessment_year")
			table.Decimal("overall_score").Total(10).Places(2).Default(0)
			table.Enum("status", []any{"DRAFT", "SUBMITTED", "VERIFIED"}).Default("DRAFT")
			table.Timestamps()

			table.Foreign("campus_id").References("id").On("campuses").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260713155153CreateCampusAssessmentsTable) Down() error {
	return facades.Schema().DropIfExists("campus_assessments")
}
