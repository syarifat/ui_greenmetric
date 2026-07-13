package seeders
	
import (
	"fmt"
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"
)

type IndicatorSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *IndicatorSeeder) Signature() string {
	return "IndicatorSeeder"
}

// Run executes the seeder logic.
func (s *IndicatorSeeder) Run() error {
	// Define list of indicators grouped by category code
	data := map[string][]struct {
		Code      string
		Title     string
		InputType string
		MaxPoints int
	}{
		"SI": {
			{"SI1", "Rasio luas ruang terbuka terhadap total luas", "NUMERIC_FORMULA", 200},
			{"SI2", "Total luas kampus yang ditutupi vegetasi hutan yang digunakan untuk riset, pengajaran, dan/atau pelibatan masyarakat", "NUMERIC_FORMULA", 100},
			{"SI3", "Total luas kampus yang ditutupi vegetasi tanam", "NUMERIC_FORMULA", 200},
			{"SI4", "Total luas ruang terbuka dibagi total populasi kampus", "NUMERIC_FORMULA", 200},
			{"SI5", "Fasilitas kampus untuk penyandang disabilitas, kebutuhan khusus, dan/atau layanan maternitas", "SINGLE_CHOICE", 100},
			{"SI6", "Fasilitas keamanan dan keselamatan", "SINGLE_CHOICE", 100},
			{"SI7", "Infrastruktur kesehatan untuk mendukung kesejahteraan mahasiswa, staf akademik, dan staf administrasi", "SINGLE_CHOICE", 100},
			{"SI8", "Konservasi flora, fauna, satwa liar, dan/atau sumber daya genetik yang diamankan dalam fasilitas konservasi jangka menengah atau panjang", "SINGLE_CHOICE", 100},
		},
		"EC": {
			{"EC1", "Penggunaan peralatan hemat energi", "NUMERIC_FORMULA", 200},
			{"EC2", "Implementasi smart building", "NUMERIC_FORMULA", 300},
			{"EC3", "Jumlah sumber energi terbarukan di kampus", "SINGLE_CHOICE", 300},
			{"EC4", "Total penggunaan listrik dibagi total populasi kampus (kWh per orang)", "NUMERIC_FORMULA", 200},
			{"EC5", "Rasio produksi energi terbarukan terhadap total penggunaan energi tahunan", "NUMERIC_FORMULA", 200},
			{"EC6", "Elemen green building yang diterapkan di seluruh bangunan", "SINGLE_CHOICE", 200},
			{"EC7", "Program pengurangan emisi gas rumah kaca (GHG)", "SINGLE_CHOICE", 200},
			{"EC8", "Total jejak karbon dibagi total populasi kampus (metrik ton per orang)", "NUMERIC_FORMULA", 200},
			{"EC9", "Jumlah program inovatif di bidang energi dan perubahan iklim", "SINGLE_CHOICE", 100},
			{"EC10", "Program universitas yang berdampak pada perubahan iklim", "SINGLE_CHOICE", 100},
		},
		"WS": {
			{"WS1", "Program 3R (Reduce, Reuse, Recycle)", "SINGLE_CHOICE", 200},
			{"WS2", "Program mengurangi penggunaan kertas dan plastik", "SINGLE_CHOICE", 300},
			{"WS3", "Pengolahan sampah organik", "SINGLE_CHOICE", 300},
			{"WS4", "Pengolahan sampah anorganik", "SINGLE_CHOICE", 300},
			{"WS5", "Pengolahan sampah beracun", "SINGLE_CHOICE", 300},
			{"WS6", "Pengolahan air limbah", "SINGLE_CHOICE", 300},
		},
		"WR": {
			{"WR1", "Total area untuk resapan air (di luar area hutan dan vegetasi tanam)", "NUMERIC_FORMULA", 100},
			{"WR2", "Program konservasi air dan implementasinya", "SINGLE_CHOICE", 200},
			{"WR3", "Implementasi program daur ulang air", "SINGLE_CHOICE", 200},
			{"WR4", "Penggunaan peralatan hemat air", "SINGLE_CHOICE", 200},
			{"WR5", "Konsumsi air olahan", "NUMERIC_FORMULA", 200},
			{"WR6", "Pengendalian pencemaran air di area kampus", "SINGLE_CHOICE", 200},
		},
		"TR": {
			{"TR1", "Total jumlah kendaraan (mobil dan motor bermesin pembakaran) dibagi total populasi kampus", "NUMERIC_FORMULA", 200},
			{"TR2", "Layanan shuttle", "SINGLE_CHOICE", 250},
			{"TR3", "Ketersediaan Zero Emission Vehicles (ZEV) di kampus", "SINGLE_CHOICE", 200},
			{"TR4", "Total jumlah ZEV dibagi total populasi kampus", "NUMERIC_FORMULA", 200},
			{"TR5", "Rasio area parkir permukaan terhadap total area kampus", "NUMERIC_FORMULA", 200},
			{"TR6", "Program untuk membatasi atau mengurangi area parkir di kampus selama 3 tahun terakhir", "SINGLE_CHOICE", 200},
			{"TR7", "Jumlah inisiatif untuk mengurangi kendaraan pribadi di kampus", "SINGLE_CHOICE", 200},
			{"TR8", "Jalur pejalan kaki di kampus", "SINGLE_CHOICE", 250},
		},
		"ED": {
			{"ED1", "Rasio mata kuliah terkait keberlanjutan", "NUMERIC_FORMULA", 200},
			{"ED2", "Rasio pendanaan riset keberlanjutan", "NUMERIC_FORMULA", 200},
			{"ED3", "Rasio publikasi keberlanjutan", "NUMERIC_FORMULA", 200},
			{"ED4", "Jumlah kegiatan keberlanjutan", "SINGLE_CHOICE", 100},
			{"ED5", "Kegiatan organisasi mahasiswa terkait keberlanjutan", "SINGLE_CHOICE", 150},
			{"ED6", "Jumlah kegiatan budaya di kampus", "SINGLE_CHOICE", 100},
			{"ED7", "Program keberlanjutan dengan kolaborasi internasional", "SINGLE_CHOICE", 100},
			{"ED8", "Program pengabdian terkait keberlanjutan melibatkan mahasiswa", "SINGLE_CHOICE", 100},
			{"ED9", "Jumlah start-up terkait keberlanjutan", "SINGLE_CHOICE", 100},
			{"ED10", "Persentase lulusan dengan green jobs", "NUMERIC_FORMULA", 50},
		},
		"GD": {
			{"GD1", "Persentase anggaran universitas untuk keberlanjutan", "NUMERIC_FORMULA", 200},
			{"GD2", "Situs web berkelanjutan", "SINGLE_CHOICE", 200},
			{"GD3", "Laporan berkelanjutan", "SINGLE_CHOICE", 100},
			{"GD4", "Laporan keuangan", "SINGLE_CHOICE", 100},
			{"GD5", "Unit atau kantor yang mengoordinasikan keberlanjutan", "SINGLE_CHOICE", 100},
			{"GD6", "Penggunaan TIK untuk perencanaan, implementasi, pemantauan, dan evaluasi program keberlanjutan", "SINGLE_CHOICE", 50},
			{"GD7", "Kebijakan penggunaan teknologi digital lanjutan (AI/IoT dan sejenisnya) untuk mendukung pengambilan keputusan, efisiensi operasional, dan layanan", "SINGLE_CHOICE", 50},
			{"GD8", "Rasio pimpinan perempuan terhadap total pimpinan institusi", "NUMERIC_FORMULA", 100},
			{"GD9", "Sistem anti-korupsi dan integritas universitas", "SINGLE_CHOICE", 50},
			{"GD10", "Sistem whistleblowing dan pengaduan universitas", "SINGLE_CHOICE", 50},
			{"GD11", "Program literasi digital berbasis LMS untuk mahasiswa dan staf", "SINGLE_CHOICE", 50},
			{"GD12", "Kode etik tertulis yang berlaku bagi pimpinan universitas, staf akademik, staf administrasi, dan mahasiswa", "SINGLE_CHOICE", 50},
		},
	}

	for catCode, indicators := range data {
		var category models.Category
		err := facades.Orm().Query().Where("code = ?", catCode).First(&category)
		if err != nil {
			return fmt.Errorf("category %s not found: %v", catCode, err)
		}

		for _, ind := range indicators {
			count, err := facades.Orm().Query().Model(&models.Indicator{}).Where("code = ?", ind.Code).Count()
			if err != nil {
				return err
			}
			if count == 0 {
				newIndicator := models.Indicator{
					CategoryID: category.ID,
					Code:       ind.Code,
					Title:      ind.Title,
					InputType:  ind.InputType,
					MaxPoints:  ind.MaxPoints,
				}
				if err := facades.Orm().Query().Create(&newIndicator); err != nil {
					return fmt.Errorf("failed to create indicator %s: %v", ind.Code, err)
				}
			}
		}
	}

	return nil
}
