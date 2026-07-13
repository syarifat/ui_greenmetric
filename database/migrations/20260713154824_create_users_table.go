package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ui_greenmetric/app/facades"
)

type M20260713154824CreateUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20260713154824CreateUsersTable) Signature() string {
	return "20260713154824_create_users_table"
}

// Up Run the migrations.
func (r *M20260713154824CreateUsersTable) Up() error {
	if !facades.Schema().HasTable("users") {
		return facades.Schema().Create("users", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("campus_id")
			table.String("name")
			table.String("email")
			table.String("password")
			table.String("role")
			table.Timestamps()

			table.Unique("email")
			table.Foreign("campus_id").References("id").On("campuses").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260713154824CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
