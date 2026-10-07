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
	// Gunakan kampus POLINEMA yang dibuat CampusSeeder (fallback: buat jika belum ada)
	campus := models.Campus{
		Code:            "POLINEMA",
		Name:            "Politeknik Negeri Malang",
		InstitutionType: "Vocational",
		Climate:         "Tropical",
		Setting:         "Urban",
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
	hashedPassword, err := facades.Hash().Make("password")
	if err != nil {
		return err
	}

	users := []models.User{
		{
			CampusID: activeCampus.ID,
			Name:     "Syarif Super Admin",
			Email:    "superadmin@gmail.com",
			Password: hashedPassword,
			Role:     "SUPER_ADMIN",
		},
		{
			CampusID: activeCampus.ID,
			Name:     "Syarif Admin Green Campus",
			Email:    "adminkampus@gmail.com",
			Password: hashedPassword,
			Role:     "ADMIN_KAMPUS",
		},
		{
			CampusID: activeCampus.ID,
			Name:     "Syarif Operator SI",
			Email:    "operatorsi@gmail.com",
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
