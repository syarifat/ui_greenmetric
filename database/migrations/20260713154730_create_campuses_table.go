package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ui_greenmetric/app/facades"
)

type M20260713154730CreateCampusesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260713154730CreateCampusesTable) Signature() string {
	return "20260713154730_create_campuses_table"
}

// Up Run the migrations.
func (r *M20260713154730CreateCampusesTable) Up() error {
	if !facades.Schema().HasTable("campuses") {
		return facades.Schema().Create("campuses", func(table schema.Blueprint) {
			table.ID()
			table.String("code")
			table.String("name")
			table.String("institution_type").Nullable()
			table.String("climate").Nullable()
			table.String("setting").Nullable()
			table.Timestamps()

			table.Unique("code")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260713154730CreateCampusesTable) Down() error {
	return facades.Schema().DropIfExists("campuses")
}
