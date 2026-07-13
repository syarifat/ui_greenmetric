package bootstrap

import (
	"github.com/goravel/framework/contracts/database/seeder"

	"ui_greenmetric/database/seeders"
)

func Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.CategorySeeder{},
		&seeders.IndicatorSeeder{},
		&seeders.IndicatorScoringTierSeeder{},
		&seeders.UserSeeder{},
	}
}
