package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ui_greenmetric/app/facades"
)

type M20260713155115CreateCategoriesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260713155115CreateCategoriesTable) Signature() string {
	return "20260713155115_create_categories_table"
}

// Up Run the migrations.
func (r *M20260713155115CreateCategoriesTable) Up() error {
	if !facades.Schema().HasTable("categories") {
		return facades.Schema().Create("categories", func(table schema.Blueprint) {
			table.ID()
			table.String("code")
			table.String("name")
			table.Integer("max_points")
			table.Decimal("weight_percentage").Total(5).Places(2)
			table.Timestamps()

			table.Unique("code")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260713155115CreateCategoriesTable) Down() error {
	return facades.Schema().DropIfExists("categories")
}
