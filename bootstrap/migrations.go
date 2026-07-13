package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ui_greenmetric/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
	}
}
