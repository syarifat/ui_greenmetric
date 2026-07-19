package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ui_greenmetric/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20260713154730CreateCampusesTable{},
		&migrations.M20260713154824CreateUsersTable{},
		&migrations.M20260713155115CreateCategoriesTable{},
		&migrations.M20260713155133CreateIndicatorsTable{},
		&migrations.M20260713155145CreateIndicatorScoringTiersTable{},
		&migrations.M20260713155153CreateCampusAssessmentsTable{},
		&migrations.M20260713155204CreateAssessmentAnswersTable{},
		&migrations.M20260713155211CreateAssessmentEvidencesTable{},
		&migrations.M20260719091940CreateIndicatorFieldsTable{},
	}
}
