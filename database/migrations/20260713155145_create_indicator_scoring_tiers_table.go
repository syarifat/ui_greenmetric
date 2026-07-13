package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ui_greenmetric/app/facades"
)

type M20260713155145CreateIndicatorScoringTiersTable struct{}

// Signature The unique signature for the migration.
func (r *M20260713155145CreateIndicatorScoringTiersTable) Signature() string {
	return "20260713155145_create_indicator_scoring_tiers_table"
}

// Up Run the migrations.
func (r *M20260713155145CreateIndicatorScoringTiersTable) Up() error {
	if !facades.Schema().HasTable("indicator_scoring_tiers") {
		return facades.Schema().Create("indicator_scoring_tiers", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("indicator_id")
			table.String("option_label")
			table.Decimal("min_value").Total(15).Places(4).Nullable()
			table.Decimal("max_value").Total(15).Places(4).Nullable()
			table.Enum("operator", []any{"<=", ">", ">=", "BETWEEN", "CHOICE"})
			table.Decimal("point_multiplier").Total(5).Places(2)
			table.Timestamps()

			table.Foreign("indicator_id").References("id").On("indicators").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260713155145CreateIndicatorScoringTiersTable) Down() error {
	return facades.Schema().DropIfExists("indicator_scoring_tiers")
}
