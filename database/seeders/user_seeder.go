package seeders
	
import (
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"
)

type UserSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *UserSeeder) Signature() string {
	return "UserSeeder"
}

// Run executes the seeder logic.
func (s *UserSeeder) Run() error {
	// Create default Campus
	campus := models.Campus{
		Code:            "POLINEMA-KDR",
		Name:            "Politeknik Negeri Malang - PSDKU Kediri",
		InstitutionType: "Vocational",
		Climate:         "Tropical",
		Setting:         "Suburban",
	}

	campusCount, err := facades.Orm().Query().Model(&models.Campus{}).Where("code = ?", campus.Code).Count()
	if err != nil {
		return err
	}

	var activeCampus models.Campus
	if campusCount == 0 {
		if err := facades.Orm().Query().Create(&campus); err != nil {
			return err
		}
		activeCampus = campus
	} else {
		if err := facades.Orm().Query().Where("code = ?", campus.Code).First(&activeCampus); err != nil {
			return err
		}
	}

	// Create default User
	hashedPassword, err := facades.Hash().Make("secretpassword")
	if err != nil {
		return err
	}

	users := []models.User{
		{
			CampusID: activeCampus.ID,
			Name:     "Syarif Super Admin",
			Email:    "super@greencampus.org",
			Password: hashedPassword,
			Role:     "SUPER_ADMIN",
		},
		{
			CampusID: activeCampus.ID,
			Name:     "Syarif Admin Green Campus",
			Email:    "admin@polinema.ac.id",
			Password: hashedPassword,
			Role:     "ADMIN_KAMPUS",
		},
		{
			CampusID: activeCampus.ID,
			Name:     "Syarif Operator SI",
			Email:    "operator.si@polinema.ac.id",
			Password: hashedPassword,
			Role:     "OPERATOR_SI",
		},
	}

	for _, user := range users {
		userCount, err := facades.Orm().Query().Model(&models.User{}).Where("email = ?", user.Email).Count()
		if err != nil {
			return err
		}

		if userCount == 0 {
			if err := facades.Orm().Query().Create(&user); err != nil {
				return err
			}
		}
	}

	return nil
}
