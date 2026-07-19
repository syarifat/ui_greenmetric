package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ui_greenmetric/app/facades"
)

type M20260719091940CreateIndicatorFieldsTable struct{}

// Signature The unique signature for the migration.
func (r *M20260719091940CreateIndicatorFieldsTable) Signature() string {
	return "20260719091940_create_indicator_fields_table"
}

// Up Run the migrations.
func (r *M20260719091940CreateIndicatorFieldsTable) Up() error {
		return facades.Schema().Create("indicator_fields", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("indicator_id")
			table.String("key")
			table.String("label")
			table.String("type") // int, float, date, varchar, choice
			table.Text("options").Nullable() // comma-separated options
			table.Boolean("required").Default(true)
			table.Timestamps()

			table.Foreign("indicator_id").References("id").On("indicators").CascadeOnDelete()
		})

	return nil
}

// Down Reverse the migrations.
func (r *M20260719091940CreateIndicatorFieldsTable) Down() error {
	return facades.Schema().DropIfExists("indicator_fields")
}
