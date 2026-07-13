package seeders
	
import (
	"fmt"
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"
)

type IndicatorScoringTierSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *IndicatorScoringTierSeeder) Signature() string {
	return "IndicatorScoringTierSeeder"
}

func floatPtr(v float64) *float64 {
	return &v
}

// Run executes the seeder logic.
func (s *IndicatorScoringTierSeeder) Run() error {
	// Define tiers for representative indicators
	data := map[string][]struct {
		OptionLabel     string
		MinValue        *float64
		MaxValue        *float64
		Operator        string
		PointMultiplier float64
	}{
		// SI1: Rasio Luas Ruang Terbuka terhadap Luas Total
		"SI1": {
			{"<= 1%", nil, floatPtr(1.0), "<=", 0.05},
			{"> 1% - 80%", floatPtr(1.0), floatPtr(80.0), "BETWEEN", 0.25},
			{"> 80% - 90%", floatPtr(80.0), floatPtr(90.0), "BETWEEN", 0.50},
			{"> 90% - 95%", floatPtr(90.0), floatPtr(95.0), "BETWEEN", 0.75},
			{"> 95%", floatPtr(95.0), nil, ">", 1.00},
		},
		// SI4: Total Luas Ruang Terbuka Dibagi Total Populasi Kampus
		"SI4": {
			{"<= 10 m2/orang", nil, floatPtr(10.0), "<=", 0.05},
			{"> 10 - 20 m2/orang", floatPtr(10.0), floatPtr(20.0), "BETWEEN", 0.25},
			{"> 20 - 40 m2/orang", floatPtr(20.0), floatPtr(40.0), "BETWEEN", 0.50},
			{"> 40 - 70 m2/orang", floatPtr(40.0), floatPtr(70.0), "BETWEEN", 0.75},
			{"> 70 m2/orang", floatPtr(70.0), nil, ">", 1.00},
		},
		// SI5: Fasilitas Kampus untuk Disabilitas
		"SI5": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"Kebijakan sudah tersedia", nil, nil, "CHOICE", 0.25},
			{"Fasilitas masih pada tahap perencanaan", nil, nil, "CHOICE", 0.50},
			{"Fasilitas tersedia sebagian dan sudah beroperasi", nil, nil, "CHOICE", 0.75},
			{"Fasilitas tersedia di semua bangunan dan beroperasi sepenuhnya", nil, nil, "CHOICE", 1.00},
		},
		// EC4: Total Penggunaan Listrik Dibagi Total Populasi Kampus (kWh per orang)
		"EC4": {
			{">= 2.400 kWh", floatPtr(2400.0), nil, ">=", 0.05},
			{"> 1.500 - 2.400 kWh", floatPtr(1500.0), floatPtr(2400.0), "BETWEEN", 0.25},
			{"> 600 - 1.500 kWh", floatPtr(600.0), floatPtr(1500.0), "BETWEEN", 0.50},
			{">= 250 - 600 kWh", floatPtr(250.0), floatPtr(600.0), "BETWEEN", 0.75},
			{"< 250 kWh", nil, floatPtr(250.0), "<", 1.00},
		},
		// WS1: Program 3R
		"WS1": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"Program 3R dalam tahap persiapan", nil, nil, "CHOICE", 0.25},
			{"Program 3R diimplementasikan 1%-50%", nil, nil, "CHOICE", 0.50},
			{"Program 3R diimplementasikan > 50-75%", nil, nil, "CHOICE", 0.75},
			{"Program 3R diimplementasikan > 75%", nil, nil, "CHOICE", 1.00},
		},
		// WR1: Total area resapan air
		"WR1": {
			{"<= 2%", nil, floatPtr(2.0), "<=", 0.05},
			{"> 2% - 10%", floatPtr(2.0), floatPtr(10.0), "BETWEEN", 0.25},
			{"> 10% - 20%", floatPtr(10.0), floatPtr(20.0), "BETWEEN", 0.50},
			{"> 20% - 40%", floatPtr(20.0), floatPtr(40.0), "BETWEEN", 0.75},
			{"> 40%", floatPtr(40.0), nil, ">", 1.00},
		},
		// TR1: Total jumlah kendaraan beremisi dibagi total populasi
		"TR1": {
			{">= 1", floatPtr(1.0), nil, ">=", 0.05},
			{"> 0.5 - 1", floatPtr(0.5), floatPtr(1.0), "BETWEEN", 0.25},
			{"> 0.125 - 0.5", floatPtr(0.125), floatPtr(0.5), "BETWEEN", 0.50},
			{"> 0.045 - 0.125", floatPtr(0.045), floatPtr(0.125), "BETWEEN", 0.75},
			{"< 0.045", nil, floatPtr(0.045), "<", 1.00},
		},
	}

	for indCode, tiers := range data {
		var indicator models.Indicator
		err := facades.Orm().Query().Where("code = ?", indCode).First(&indicator)
		if err != nil {
			return fmt.Errorf("indicator %s not found: %v", indCode, err)
		}

		for _, t := range tiers {
			count, err := facades.Orm().Query().Model(&models.IndicatorScoringTier{}).
				Where("indicator_id = ? AND option_label = ?", indicator.ID, t.OptionLabel).
				Count()
			if err != nil {
				return err
			}

			if count == 0 {
				newTier := models.IndicatorScoringTier{
					IndicatorID:     indicator.ID,
					OptionLabel:     t.OptionLabel,
					MinValue:        t.MinValue,
					MaxValue:        t.MaxValue,
					Operator:        t.Operator,
					PointMultiplier: t.PointMultiplier,
				}
				if err := facades.Orm().Query().Create(&newTier); err != nil {
					return fmt.Errorf("failed to create scoring tier for indicator %s: %v", indCode, err)
				}
			}
		}
	}

	return nil
}
