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
			var dbIndicator models.Indicator
			_ = facades.Orm().Query().Where("code = ?", ind.Code).First(&dbIndicator)
			if dbIndicator.ID == 0 {
				// Create new indicator
				dbIndicator = models.Indicator{
					CategoryID: category.ID,
					Code:       ind.Code,
					Title:      ind.Title,
					InputType:  ind.InputType,
					MaxPoints:  ind.MaxPoints,
				}
				if err := facades.Orm().Query().Create(&dbIndicator); err != nil {
					return fmt.Errorf("failed to create indicator %s: %v", ind.Code, err)
				}
			}

			// Seed fields for this indicator
			fields := getFieldsForIndicator(ind.Code, ind.InputType)
			for _, f := range fields {
				count, err := facades.Orm().Query().Model(&models.IndicatorField{}).
					Where("indicator_id = ? AND `key` = ?", dbIndicator.ID, f.Key).Count()
				if err != nil {
					return err
				}
				if count == 0 {
					f.IndicatorID = dbIndicator.ID
					if err := facades.Orm().Query().Create(&f); err != nil {
						return fmt.Errorf("failed to create field %s for indicator %s: %v", f.Key, ind.Code, err)
					}
				}
			}
		}
	}

	return nil
}

func getFieldsForIndicator(code string, inputType string) []models.IndicatorField {
	if inputType == "SINGLE_CHOICE" {
		return []models.IndicatorField{
			{Key: "option_label", Label: "Pilih Opsi Kriteria", Type: "choice", Required: true},
		}
	}

	switch code {
	case "SI1":
		return []models.IndicatorField{
			{Key: "luas_total", Label: "Total Luas Kampus (m2)", Type: "float", Required: true},
			{Key: "luas_dasar", Label: "Total Luas Lantai Dasar Bangunan (m2)", Type: "float", Required: true},
		}
	case "SI2":
		return []models.IndicatorField{
			{Key: "luas_hutan", Label: "Total Luas Kampus Ditutupi Hutan (m2)", Type: "float", Required: true},
			{Key: "luas_total", Label: "Total Luas Kampus (m2)", Type: "float", Required: true},
		}
	case "SI3":
		return []models.IndicatorField{
			{Key: "luas_vegetasi", Label: "Total Luas Kampus Ditutupi Vegetasi Tanam (m2)", Type: "float", Required: true},
			{Key: "luas_total", Label: "Total Luas Kampus (m2)", Type: "float", Required: true},
		}
	case "SI4":
		return []models.IndicatorField{
			{Key: "luas_total", Label: "Total Luas Kampus (m2)", Type: "float", Required: true},
			{Key: "luas_dasar", Label: "Total Luas Lantai Dasar Bangunan (m2)", Type: "float", Required: true},
			{Key: "populasi", Label: "Total Populasi Kampus (Orang)", Type: "float", Required: true},
		}
	case "EC1":
		return []models.IndicatorField{
			{Key: "persentase_alat_hemat_energi", Label: "Persentase Peralatan Hemat Energi (%)", Type: "float", Required: true},
		}
	case "EC2":
		return []models.IndicatorField{
			{Key: "luas_smart_building", Label: "Total Luas Smart Building (m2)", Type: "float", Required: true},
			{Key: "luas_total_bangunan", Label: "Total Luas Bangunan Kampus (m2)", Type: "float", Required: true},
		}
	case "EC4":
		return []models.IndicatorField{
			{Key: "total_listrik", Label: "Total Penggunaan Listrik Tahunan (kWh)", Type: "float", Required: true},
			{Key: "populasi", Label: "Total Populasi Kampus (Orang)", Type: "float", Required: true},
		}
	case "EC5":
		return []models.IndicatorField{
			{Key: "produksi_energi_terbarukan", Label: "Total Produksi Energi Terbarukan (kWh)", Type: "float", Required: true},
			{Key: "total_penggunaan_energi", Label: "Total Penggunaan Energi Tahunan (kWh)", Type: "float", Required: true},
		}
	case "EC8":
		return []models.IndicatorField{
			{Key: "jejak_karbon", Label: "Total Jejak Karbon Tahunan (Metrik Ton CO2)", Type: "float", Required: true},
			{Key: "populasi", Label: "Total Populasi Kampus (Orang)", Type: "float", Required: true},
		}
	case "WR1":
		return []models.IndicatorField{
			{Key: "area_resapan", Label: "Total Area Resapan Air (m2)", Type: "float", Required: true},
			{Key: "luas_total", Label: "Total Luas Kampus (m2)", Type: "float", Required: true},
		}
	case "WR5":
		return []models.IndicatorField{
			{Key: "air_olahan_dikonsumsi", Label: "Total Konsumsi Air Olahan (m3)", Type: "float", Required: true},
			{Key: "total_air_dikonsumsi", Label: "Total Konsumsi Air Kampus (m3)", Type: "float", Required: true},
		}
	case "TR1":
		return []models.IndicatorField{
			{Key: "total_kendaraan", Label: "Total Jumlah Kendaraan Mobil & Motor Bermesin Pembakaran", Type: "float", Required: true},
			{Key: "populasi", Label: "Total Populasi Kampus (Orang)", Type: "float", Required: true},
		}
	case "TR4":
		return []models.IndicatorField{
			{Key: "total_zev", Label: "Total Jumlah Zero Emission Vehicles (ZEV)", Type: "float", Required: true},
			{Key: "populasi", Label: "Total Populasi Kampus (Orang)", Type: "float", Required: true},
		}
	case "TR5":
		return []models.IndicatorField{
			{Key: "luas_parkir", Label: "Total Luas Area Parkir Permukaan (m2)", Type: "float", Required: true},
			{Key: "luas_total", Label: "Total Luas Kampus (m2)", Type: "float", Required: true},
		}
	case "ED1":
		return []models.IndicatorField{
			{Key: "mk_keberlanjutan", Label: "Jumlah Mata Kuliah Terkait Keberlanjutan", Type: "float", Required: true},
			{Key: "total_mk", Label: "Total Jumlah Mata Kuliah", Type: "float", Required: true},
		}
	case "ED2":
		return []models.IndicatorField{
			{Key: "dana_riset_keberlanjutan", Label: "Total Dana Riset Terkait Keberlanjutan (Rp)", Type: "float", Required: true},
			{Key: "total_dana_riset", Label: "Total Dana Riset Universitas (Rp)", Type: "float", Required: true},
		}
	case "ED3":
		return []models.IndicatorField{
			{Key: "publikasi_keberlanjutan", Label: "Jumlah Publikasi Keberlanjutan", Type: "float", Required: true},
			{Key: "total_publikasi", Label: "Total Jumlah Publikasi", Type: "float", Required: true},
		}
	case "ED10":
		return []models.IndicatorField{
			{Key: "lulusan_green_jobs", Label: "Jumlah Lulusan dengan Green Jobs", Type: "float", Required: true},
			{Key: "total_lulusan", Label: "Total Jumlah Lulusan", Type: "float", Required: true},
		}
	case "GD1":
		return []models.IndicatorField{
			{Key: "anggaran_keberlanjutan", Label: "Total Anggaran Universitas Untuk Keberlanjutan (Rp)", Type: "float", Required: true},
			{Key: "total_anggaran", Label: "Total Anggaran Universitas (Rp)", Type: "float", Required: true},
		}
	case "GD8":
		return []models.IndicatorField{
			{Key: "pimpinan_perempuan", Label: "Jumlah Pimpinan Perempuan", Type: "float", Required: true},
			{Key: "total_pimpinan", Label: "Total Jumlah Pimpinan", Type: "float", Required: true},
		}
	default:
		return []models.IndicatorField{
			{Key: "value", Label: "Nilai Masukan", Type: "float", Required: true},
		}
	}
}
