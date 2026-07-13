package seeders
	
import (
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"
)

type CategorySeeder struct {
}

// Signature The name and signature of the seeder.
func (s *CategorySeeder) Signature() string {
	return "CategorySeeder"
}

// Run executes the seeder logic.
func (s *CategorySeeder) Run() error {
	categories := []models.Category{
		{Code: "SI", Name: "Setting and Infrastructure", MaxPoints: 1100, WeightPercentage: 11.0},
		{Code: "EC", Name: "Energy and Climate Change", MaxPoints: 2000, WeightPercentage: 20.0},
		{Code: "WS", Name: "Waste", MaxPoints: 1700, WeightPercentage: 17.0},
		{Code: "WR", Name: "Water", MaxPoints: 1100, WeightPercentage: 11.0},
		{Code: "TR", Name: "Transportation", MaxPoints: 1700, WeightPercentage: 17.0},
		{Code: "ED", Name: "Education and Research", MaxPoints: 1300, WeightPercentage: 13.0},
		{Code: "GD", Name: "Governance and Digitalization", MaxPoints: 1100, WeightPercentage: 11.0},
	}

	for _, category := range categories {
		// Check if already exists to prevent duplicate key error when running multiple times
		count, err := facades.Orm().Query().Model(&models.Category{}).Where("code = ?", category.Code).Count()
		if err != nil {
			return err
		}
		if count == 0 {
			if err := facades.Orm().Query().Create(&category); err != nil {
				return err
			}
		}
	}

	return nil
}
