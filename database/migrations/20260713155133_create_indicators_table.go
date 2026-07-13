package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ui_greenmetric/app/facades"
)

type M20260713155133CreateIndicatorsTable struct{}

// Signature The unique signature for the migration.
func (r *M20260713155133CreateIndicatorsTable) Signature() string {
	return "20260713155133_create_indicators_table"
}

// Up Run the migrations.
func (r *M20260713155133CreateIndicatorsTable) Up() error {
	if !facades.Schema().HasTable("indicators") {
		return facades.Schema().Create("indicators", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("category_id")
			table.String("code")
			table.String("title", 500)
			table.Enum("input_type", []any{"NUMERIC_FORMULA", "SINGLE_CHOICE"})
			table.Integer("max_points")
			table.Timestamps()

			table.Unique("code")
			table.Foreign("category_id").References("id").On("categories").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260713155133CreateIndicatorsTable) Down() error {
	return facades.Schema().DropIfExists("indicators")
}
