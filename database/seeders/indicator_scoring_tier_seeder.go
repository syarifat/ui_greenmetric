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
	data := map[string][]struct {
		OptionLabel     string
		MinValue        *float64
		MaxValue        *float64
		Operator        string
		PointMultiplier float64
	}{
		// EC1
		"EC1": {
			{"<= 1%", nil, floatPtr(1.0), "<=", 0.05},
			{"> 1% - 25%", floatPtr(1.0), floatPtr(25.0), "BETWEEN", 0.25},
			{"> 25% - 50%", floatPtr(25.0), floatPtr(50.0), "BETWEEN", 0.50},
			{"> 50% - 75%", floatPtr(50.0), floatPtr(75.0), "BETWEEN", 0.75},
			{"> 75%", floatPtr(75.0), nil, ">", 1.00},
		},
		// EC10
		"EC10": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"Program dalam tahap persiapan", nil, nil, "CHOICE", 0.25},
			{"Pelatihan/ materi/ seminar/ kegiatan dilaksanakan bersama komunitas sekitar (tingkat lokal)", nil, nil, "CHOICE", 0.50},
			{"Pelatihan/ materi/ seminar/ kegiatan dilaksanakan pada tingkat nasional", nil, nil, "CHOICE", 0.75},
			{"Pelatihan/ materi/ seminar/ kegiatan dilaksanakan pada tingkat internasional", nil, nil, "CHOICE", 1.00},
		},
		// EC2
		"EC2": {
			{"<= 1%", nil, floatPtr(1.0), "<=", 0.05},
			{"> 1% - 25%", floatPtr(1.0), floatPtr(25.0), "BETWEEN", 0.25},
			{"> 25% - 50%", floatPtr(25.0), floatPtr(50.0), "BETWEEN", 0.50},
			{"> 50% - 75%", floatPtr(50.0), floatPtr(75.0), "BETWEEN", 0.75},
			{"> 75%", floatPtr(75.0), nil, ">", 1.00},
		},
		// EC3
		"EC3": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"1 sumber", nil, nil, "CHOICE", 0.25},
			{"2 sumber", nil, nil, "CHOICE", 0.50},
			{"3 sumber", nil, nil, "CHOICE", 0.75},
			{"Lebih dari 3 sumber", nil, nil, "CHOICE", 1.00},
		},
		// EC4
		"EC4": {
			{">= 2400 kWh", floatPtr(2400.0), nil, ">=", 0.05},
			{"> 1500 - 2400 kWh", floatPtr(1500.0), floatPtr(2400.0), "BETWEEN", 0.25},
			{"> 600 - 1500 kWh", floatPtr(600.0), floatPtr(1500.0), "BETWEEN", 0.50},
			{">= 250 - 600 kWh", floatPtr(250.0), floatPtr(600.0), "BETWEEN", 0.75},
			{"< 250 kWh", nil, floatPtr(250.0), "<", 1.00},
		},
		// EC5
		"EC5": {
			{"<= 0.5%", nil, floatPtr(0.5), "<=", 0.05},
			{"> 0.5% - 1%", floatPtr(0.5), floatPtr(1.0), "BETWEEN", 0.25},
			{"> 1% - 2%", floatPtr(1.0), floatPtr(2.0), "BETWEEN", 0.50},
			{"> 2% - 25%", floatPtr(2.0), floatPtr(25.0), "BETWEEN", 0.75},
			{"> 25%", floatPtr(25.0), nil, ">", 1.00},
		},
		// EC6
		"EC6": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"1 elemen", nil, nil, "CHOICE", 0.25},
			{"2 elemen", nil, nil, "CHOICE", 0.50},
			{"3 elemen", nil, nil, "CHOICE", 0.75},
			{"4 elemen", nil, nil, "CHOICE", 1.00},
		},
		// EC7
		"EC7": {
			{"Tidak ada (program diperlukan, tapi belum ada tindakan)", nil, nil, "CHOICE", 0.00},
			{"Program dalam tahap persiapan", nil, nil, "CHOICE", 0.25},
			{"Program menangani emisi dalam satu cakupan (Scope 1/ 2/ 3)", nil, nil, "CHOICE", 0.50},
			{"Program menangani emisi dalam dua cakupan (Scope 1 & 2 / 1 & 3 / 2 & 3)", nil, nil, "CHOICE", 0.75},
			{"Program menangani emisi dalam ketiga cakupan (Scope 1, 2, dan 3)", nil, nil, "CHOICE", 1.00},
		},
		// EC8
		"EC8": {
			{"> 2.05 metrik ton", floatPtr(2.05), nil, ">", 0.05},
			{"> 1.11 - 2.05 metrik ton", floatPtr(1.11), floatPtr(2.05), "BETWEEN", 0.25},
			{"> 0.42 - 1.11 metrik ton", floatPtr(0.42), floatPtr(1.11), "BETWEEN", 0.50},
			{"> 0.10 - 0.42 metrik ton", floatPtr(0.1), floatPtr(0.42), "BETWEEN", 0.75},
			{"< 0.10 metrik ton", nil, floatPtr(0.1), "<", 1.00},
		},
		// EC9
		"EC9": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"1 program", nil, nil, "CHOICE", 0.25},
			{"2 program", nil, nil, "CHOICE", 0.50},
			{"3 program", nil, nil, "CHOICE", 0.75},
			{"> 3 program", nil, nil, "CHOICE", 1.00},
		},
		// ED1
		"ED1": {
			{"<= 1%", nil, floatPtr(1.0), "<=", 0.05},
			{"> 1% - 5%", floatPtr(1.0), floatPtr(5.0), "BETWEEN", 0.25},
			{"> 5% - 10%", floatPtr(5.0), floatPtr(10.0), "BETWEEN", 0.50},
			{"> 10% - 20%", floatPtr(10.0), floatPtr(20.0), "BETWEEN", 0.75},
			{"> 20%", floatPtr(20.0), nil, ">", 1.00},
		},
		// ED10
		"ED10": {
			{"<= 1%", nil, floatPtr(1.0), "<=", 0.05},
			{"> 1% - 5%", floatPtr(1.0), floatPtr(5.0), "BETWEEN", 0.25},
			{"> 5% - 10%", floatPtr(5.0), floatPtr(10.0), "BETWEEN", 0.50},
			{"> 10% - 20%", floatPtr(10.0), floatPtr(20.0), "BETWEEN", 0.75},
			{"> 20%", floatPtr(20.0), nil, ">", 1.00},
		},
		// ED2
		"ED2": {
			{"<= 1%", nil, floatPtr(1.0), "<=", 0.05},
			{"> 1% - 10%", floatPtr(1.0), floatPtr(10.0), "BETWEEN", 0.25},
			{"> 10% - 20%", floatPtr(10.0), floatPtr(20.0), "BETWEEN", 0.50},
			{"> 20% - 40%", floatPtr(20.0), floatPtr(40.0), "BETWEEN", 0.75},
			{"> 40%", floatPtr(40.0), nil, ">", 1.00},
		},
		// ED3
		"ED3": {
			{"< 0.5", nil, floatPtr(0.5), "<", 0.00},
			{"> 0.5 - 1", floatPtr(0.5), floatPtr(1.0), "BETWEEN", 0.25},
			{"> 1 - 2", floatPtr(1.0), floatPtr(2.0), "BETWEEN", 0.50},
			{"> 2 - 3", floatPtr(2.0), floatPtr(3.0), "BETWEEN", 0.75},
			{"> 3", floatPtr(3.0), nil, ">", 1.00},
		},
		// ED4
		"ED4": {
			{"0", nil, nil, "CHOICE", 0.00},
			{"1-5", nil, nil, "CHOICE", 0.25},
			{"6-20", nil, nil, "CHOICE", 0.50},
			{"21-50", nil, nil, "CHOICE", 0.75},
			{"> 50", nil, nil, "CHOICE", 1.00},
		},
		// ED5
		"ED5": {
			{"0", nil, nil, "CHOICE", 0.00},
			{"1-5", nil, nil, "CHOICE", 0.25},
			{"6-10", nil, nil, "CHOICE", 0.50},
			{"11-20", nil, nil, "CHOICE", 0.75},
			{"> 20", nil, nil, "CHOICE", 1.00},
		},
		// ED6
		"ED6": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"1-3 kegiatan per tahun", nil, nil, "CHOICE", 0.25},
			{"4-6 kegiatan per tahun", nil, nil, "CHOICE", 0.50},
			{"7-10 kegiatan per tahun", nil, nil, "CHOICE", 0.75},
			{"Lebih dari 10 kegiatan per tahun", nil, nil, "CHOICE", 1.00},
		},
		// ED7
		"ED7": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"1-3 program per tahun", nil, nil, "CHOICE", 0.25},
			{"4-6 program per tahun", nil, nil, "CHOICE", 0.50},
			{"7-10 program per tahun", nil, nil, "CHOICE", 0.75},
			{"Lebih dari 10 program per tahun", nil, nil, "CHOICE", 1.00},
		},
		// ED8
		"ED8": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"1-3 program per tahun", nil, nil, "CHOICE", 0.25},
			{"4-6 program per tahun", nil, nil, "CHOICE", 0.50},
			{"7-10 program per tahun", nil, nil, "CHOICE", 0.75},
			{"Lebih dari 10 program per tahun", nil, nil, "CHOICE", 1.00},
		},
		// ED9
		"ED9": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"1-5 start-up", nil, nil, "CHOICE", 0.25},
			{"6-10 start-up", nil, nil, "CHOICE", 0.50},
			{"11-15 start-up", nil, nil, "CHOICE", 0.75},
			{"Lebih dari 15 start-up", nil, nil, "CHOICE", 1.00},
		},
		// GD1
		"GD1": {
			{"<= 1%", nil, floatPtr(1.0), "<=", 0.05},
			{"> 1% - 5%", floatPtr(1.0), floatPtr(5.0), "BETWEEN", 0.25},
			{"> 5% - 10%", floatPtr(5.0), floatPtr(10.0), "BETWEEN", 0.50},
			{"> 10% - 15%", floatPtr(10.0), floatPtr(15.0), "BETWEEN", 0.75},
			{"> 15%", floatPtr(15.0), nil, ">", 1.00},
		},
		// GD10
		"GD10": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"Sistem whistleblowing dan pengaduan masih pada tahap perencanaan", nil, nil, "CHOICE", 0.25},
			{"Sistem whistleblowing dan pengaduan sudah diimplementasikan", nil, nil, "CHOICE", 0.50},
			{"Sistem whistleblowing dan pengaduan sudah diimplementasikan dan dievaluasi", nil, nil, "CHOICE", 0.75},
			{"Sistem whistleblowing dan pengaduan sudah diimplementasikan, dievaluasi, dan saat ini sedang direvisi", nil, nil, "CHOICE", 1.00},
		},
		// GD11
		"GD11": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"Program masih pada tahap perencanaan", nil, nil, "CHOICE", 0.25},
			{"Program sudah diimplementasikan", nil, nil, "CHOICE", 0.50},
			{"Program sudah diimplementasikan dan dievaluasi", nil, nil, "CHOICE", 0.75},
			{"Program sudah diimplementasikan, dievaluasi, dan saat ini sedang direvisi", nil, nil, "CHOICE", 1.00},
		},
		// GD12
		"GD12": {
			{"Tidak ada kode etik tertulis", nil, nil, "CHOICE", 0.00},
			{"Kode etik sedang dipersiapkan atau masih berupa draf", nil, nil, "CHOICE", 0.25},
			{"Kode etik tertulis sudah ditetapkan secara formal, tetapi hanya berlaku untuk kelompok tertentu atau belum ditegakkan secara konsisten", nil, nil, "CHOICE", 0.50},
			{"Kode etik tertulis berlaku untuk semua kelompok dan sudah diimplementasikan serta dipantau", nil, nil, "CHOICE", 0.75},
			{"Kode etik tertulis berlaku untuk semua kelompok, diimplementasikan penuh, ditinjau secara berkala, dan ditegakkan secara aktif melalui mekanisme institusional", nil, nil, "CHOICE", 1.00},
		},
		// GD2
		"GD2": {
			{"Tidak tersedia", nil, nil, "CHOICE", 0.00},
			{"Situs web masih dalam proses atau sedang dibangun", nil, nil, "CHOICE", 0.25},
			{"Situs web tersedia dan dapat diakses", nil, nil, "CHOICE", 0.50},
			{"Situs web tersedia, dapat diakses, dan diperbarui sesekali", nil, nil, "CHOICE", 0.75},
			{"Situs web tersedia, dapat diakses, dan diperbarui secara rutin", nil, nil, "CHOICE", 1.00},
		},
		// GD3
		"GD3": {
			{"Tidak tersedia", nil, nil, "CHOICE", 0.00},
			{"Laporan keberlanjutan dalam tahap persiapan", nil, nil, "CHOICE", 0.25},
			{"Tersedia tetapi tidak dapat diakses publik", nil, nil, "CHOICE", 0.50},
			{"Laporan keberlanjutan dapat diakses dan diterbitkan sesekali", nil, nil, "CHOICE", 0.75},
			{"Laporan keberlanjutan dapat diakses dan diterbitkan setiap tahun", nil, nil, "CHOICE", 1.00},
		},
		// GD4
		"GD4": {
			{"Tidak tersedia", nil, nil, "CHOICE", 0.00},
			{"Laporan keuangan dalam tahap persiapan", nil, nil, "CHOICE", 0.25},
			{"Tersedia tetapi tidak dapat diakses publik", nil, nil, "CHOICE", 0.50},
			{"Laporan keuangan dapat diakses dan diterbitkan sesekali", nil, nil, "CHOICE", 0.75},
			{"Laporan keuangan dapat diakses dan diterbitkan setiap tahun", nil, nil, "CHOICE", 1.00},
		},
		// GD5
		"GD5": {
			{"Ad-hoc atau satuan tugas", nil, nil, "CHOICE", 0.00},
			{"Unit atau kantor masih dalam pengembangan", nil, nil, "CHOICE", 0.25},
			{"Unit atau kantor memiliki SK pimpinan universitas, struktur dan tugas pada tahap awal", nil, nil, "CHOICE", 0.50},
			{"Unit atau kantor memiliki SK pimpinan universitas, struktur dan tugas, dan sudah operasional", nil, nil, "CHOICE", 0.75},
			{"Unit atau kantor memiliki SK pimpinan universitas, struktur dan tugas, sudah operasional, dan memimpin implementasi keberlanjutan universitas", nil, nil, "CHOICE", 1.00},
		},
		// GD6
		"GD6": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"Tahap perencanaan", nil, nil, "CHOICE", 0.25},
			{"Sudah diimplementasikan", nil, nil, "CHOICE", 0.50},
			{"Sudah diimplementasikan dan dievaluasi", nil, nil, "CHOICE", 0.75},
			{"Sudah diimplementasikan, dievaluasi, dan saat ini sedang direvisi atau ditingkatkan", nil, nil, "CHOICE", 1.00},
		},
		// GD7
		"GD7": {
			{"Tidak ada kebijakan", nil, nil, "CHOICE", 0.00},
			{"Adopsi awal kebijakan, implementasi kebijakan terbatas pada unit tertentu", nil, nil, "CHOICE", 0.25},
			{"Implementasi sebagian, kebijakan digunakan pada beberapa proses administrasi atau akademik tetapi belum terintegrasi di tingkat institusi", nil, nil, "CHOICE", 0.50},
			{"Implementasi luas, kebijakan terintegrasi pada berbagai fungsi administrasi dan akademik serta mendukung pengambilan keputusan rutin dan layanan", nil, nil, "CHOICE", 0.75},
			{"Implementasi maju dan terintegrasi, kebijakan diterapkan di seluruh institusi, secara sistematis mendukung pengambilan keputusan strategis, optimasi operasional, dan layanan, serta dievaluasi dan ditingkatkan secara berkelanjutan", nil, nil, "CHOICE", 1.00},
		},
		// GD8
		"GD8": {
			{"<= 5%", nil, floatPtr(5.0), "<=", 0.00},
			{"> 5% - 20%", floatPtr(5.0), floatPtr(20.0), "BETWEEN", 0.25},
			{"> 20% - 35%", floatPtr(20.0), floatPtr(35.0), "BETWEEN", 0.50},
			{"> 35% - 50%", floatPtr(35.0), floatPtr(50.0), "BETWEEN", 0.75},
			{"> 50%", floatPtr(50.0), nil, ">", 1.00},
		},
		// GD9
		"GD9": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"Sistem anti-korupsi dan integritas masih pada tahap perencanaan", nil, nil, "CHOICE", 0.25},
			{"Sistem anti-korupsi dan integritas sudah diimplementasikan", nil, nil, "CHOICE", 0.50},
			{"Sistem anti-korupsi dan integritas sudah diimplementasikan dan dievaluasi", nil, nil, "CHOICE", 0.75},
			{"Sistem anti-korupsi dan integritas sudah diimplementasikan, dievaluasi, dan saat ini sedang direvisi", nil, nil, "CHOICE", 1.00},
		},
		// SI1
		"SI1": {
			{"<= 1%", nil, floatPtr(1.0), "<=", 0.05},
			{"> 1% - 80%", floatPtr(1.0), floatPtr(80.0), "BETWEEN", 0.25},
			{"> 80% - 90%", floatPtr(80.0), floatPtr(90.0), "BETWEEN", 0.50},
			{"> 90% - 95%", floatPtr(90.0), floatPtr(95.0), "BETWEEN", 0.75},
			{"> 95%", floatPtr(95.0), nil, ">", 1.00},
		},
		// SI2
		"SI2": {
			{"<= 2%", nil, floatPtr(2.0), "<=", 0.05},
			{"> 2% - 10%", floatPtr(2.0), floatPtr(10.0), "BETWEEN", 0.25},
			{"> 10% - 25%", floatPtr(10.0), floatPtr(25.0), "BETWEEN", 0.50},
			{"> 25% - 35%", floatPtr(25.0), floatPtr(35.0), "BETWEEN", 0.75},
			{"> 35%", floatPtr(35.0), nil, ">", 1.00},
		},
		// SI3
		"SI3": {
			{"<= 10%", nil, floatPtr(10.0), "<=", 0.05},
			{"> 10% - 20%", floatPtr(10.0), floatPtr(20.0), "BETWEEN", 0.25},
			{"> 20% - 30%", floatPtr(20.0), floatPtr(30.0), "BETWEEN", 0.50},
			{"> 30% - 50%", floatPtr(30.0), floatPtr(50.0), "BETWEEN", 0.75},
			{"> 50%", floatPtr(50.0), nil, ">", 1.00},
		},
		// SI4
		"SI4": {
			{"<= 10 m2/orang", nil, floatPtr(10.0), "<=", 0.05},
			{"> 10 - 20 m2/orang", floatPtr(10.0), floatPtr(20.0), "BETWEEN", 0.25},
			{"> 20 - 40 m2/orang", floatPtr(20.0), floatPtr(40.0), "BETWEEN", 0.50},
			{"> 40 - 70 m2/orang", floatPtr(40.0), floatPtr(70.0), "BETWEEN", 0.75},
			{"> 70 m2/orang", floatPtr(70.0), nil, ">", 1.00},
		},
		// SI5
		"SI5": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"Kebijakan sudah tersedia", nil, nil, "CHOICE", 0.25},
			{"Fasilitas masih pada tahap perencanaan", nil, nil, "CHOICE", 0.50},
			{"Fasilitas tersedia sebagian dan sudah beroperasi", nil, nil, "CHOICE", 0.75},
			{"Fasilitas tersedia di semua bangunan dan beroperasi sepenuhnya", nil, nil, "CHOICE", 1.00},
		},
		// SI6
		"SI6": {
			{"Sistem keamanan dan keselamatan pasif", nil, nil, "CHOICE", 0.00},
			{"Sistem keamanan dan keselamatan (CCTV, hotline/ tombol darurat) tersedia dan berfungsi penuh", nil, nil, "CHOICE", 0.25},
			{"Sistem keamanan dan keselamatan (CCTV, hotline/ tombol darurat, personel tersertifikasi, APAR, hidran) tersedia dan berfungsi penuh", nil, nil, "CHOICE", 0.50},
			{"Infrastruktur keamanan dan keselamatan tersedia dan berfungsi penuh, serta waktu respons untuk kecelakaan, kejahatan, kebakaran, dan bencana alam lebih dari 5 menit", nil, nil, "CHOICE", 0.75},
			{"Infrastruktur keamanan dan keselamatan tersedia dan berfungsi penuh, serta waktu respons untuk kecelakaan, kejahatan, kebakaran, dan bencana alam kurang dari 5 menit", nil, nil, "CHOICE", 1.00},
		},
		// SI7
		"SI7": {
			{"Infrastruktur kesehatan (pertolongan pertama) tidak tersedia", nil, nil, "CHOICE", 0.00},
			{"Infrastruktur kesehatan (pertolongan pertama, IGD, klinik, dan personel) tersedia", nil, nil, "CHOICE", 0.25},
			{"Infrastruktur kesehatan (pertolongan pertama, IGD, klinik, dan personel tersertifikasi) tersedia", nil, nil, "CHOICE", 0.50},
			{"Infrastruktur kesehatan (pertolongan pertama, IGD, klinik, rumah sakit, dan personel tersertifikasi) tersedia", nil, nil, "CHOICE", 0.75},
			{"Infrastruktur kesehatan (pertolongan pertama, IGD, klinik, rumah sakit, dan personel tersertifikasi) tersedia, tersistem, dan dapat diakses publik", nil, nil, "CHOICE", 1.00},
		},
		// SI8
		"SI8": {
			{"Program konservasi sedang dipersiapkan", nil, nil, "CHOICE", 0.05},
			{"Program konservasi sudah diimplementasikan 1% - 25%", nil, nil, "CHOICE", 0.25},
			{"Program konservasi sudah diimplementasikan 25% - 50%", nil, nil, "CHOICE", 0.50},
			{"Program konservasi sudah diimplementasikan 50% - 75%", nil, nil, "CHOICE", 0.75},
			{"Program konservasi sudah diimplementasikan lebih dari 75%", nil, nil, "CHOICE", 1.00},
		},
		// TR1
		"TR1": {
			{">= 1", floatPtr(1.0), nil, ">=", 0.05},
			{"> 0.5 - 1", floatPtr(0.5), floatPtr(1.0), "BETWEEN", 0.25},
			{"> 0.125 - 0.5", floatPtr(0.125), floatPtr(0.5), "BETWEEN", 0.50},
			{"> 0.045 - 0.125", floatPtr(0.045), floatPtr(0.125), "BETWEEN", 0.75},
			{"< 0.045", nil, floatPtr(0.045), "<", 1.00},
		},
		// TR2
		"TR2": {
			{"Memungkinkan, tetapi tidak disediakan oleh universitas 0", nil, nil, "CHOICE", 0.00},
			{"Disediakan (oleh universitas atau pihak lain), rutin, tetapi tidak gratis", nil, nil, "CHOICE", 0.25},
			{"Disediakan (oleh universitas atau pihak lain), dan universitas mensubsidi sebagian biaya", nil, nil, "CHOICE", 0.50},
			{"Disediakan oleh universitas, rutin, dan gratis", nil, nil, "CHOICE", 0.75},
			{"Disediakan oleh universitas, rutin, dan dioperasikan menggunakan kendaraan tanpa emisi", nil, nil, "CHOICE", 1.00},
		},
		// TR3
		"TR3": {
			{"Penggunaan ZEV tidak memungkinkan atau tidak praktis", nil, nil, "CHOICE", 0.00},
			{"ZEV tersedia, tetapi tidak disediakan oleh universitas", nil, nil, "CHOICE", 0.00},
			{"ZEV tersedia, disediakan oleh universitas, dan dikenakan biaya (berbayar)", nil, nil, "CHOICE", 0.00},
			{"ZEV tersedia dan disediakan oleh universitas secara gratis* *Digunakan secara rutin oleh komunitas kampus. Bukti wajib dilampirkan", nil, nil, "CHOICE", 0.00},
		},
		// TR4
		"TR4": {
			{"<= 0.002", nil, floatPtr(0.002), "<=", 0.05},
			{"> 0.002 - 0.004", floatPtr(0.002), floatPtr(0.004), "BETWEEN", 0.25},
			{"> 0.004 - 0.008", floatPtr(0.004), floatPtr(0.008), "BETWEEN", 0.50},
			{"> 0.008 - 0.02", floatPtr(0.008), floatPtr(0.02), "BETWEEN", 0.75},
			{"> 0.02", floatPtr(0.02), nil, ">", 1.00},
		},
		// TR5
		"TR5": {
			{"> 11%", floatPtr(11.0), nil, ">", 0.05},
			{"> 7% - 11%", floatPtr(7.0), floatPtr(11.0), "BETWEEN", 0.25},
			{"> 4% - 7%", floatPtr(4.0), floatPtr(7.0), "BETWEEN", 0.50},
			{"> 1% - 4%", floatPtr(1.0), floatPtr(4.0), "BETWEEN", 0.75},
			{"< 1%", nil, floatPtr(1.0), "<", 1.00},
		},
		// TR6
		"TR6": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"Dalam tahap persiapan", nil, nil, "CHOICE", 0.25},
			{"Penurunan area parkir kurang dari 10%", nil, nil, "CHOICE", 0.50},
			{"Penurunan area parkir 10%-30%", nil, nil, "CHOICE", 0.75},
			{"Penurunan area parkir lebih dari 30%, atau pengurangan area parkir telah mencapai batas praktisnya", nil, nil, "CHOICE", 1.00},
		},
		// TR7
		"TR7": {
			{"Tidak ada inisiatif", nil, nil, "CHOICE", 0.00},
			{"1 inisiatif", nil, nil, "CHOICE", 0.25},
			{"2 inisiatif", nil, nil, "CHOICE", 0.50},
			{"3 inisiatif", nil, nil, "CHOICE", 0.75},
			{"> 3 inisiatif, atau inisiatif sudah tidak diperlukan", nil, nil, "CHOICE", 1.00},
		},
		// TR8
		"TR8": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"Tersedia", nil, nil, "CHOICE", 0.25},
			{"Tersedia dan dirancang untuk keselamatan", nil, nil, "CHOICE", 0.50},
			{"Tersedia dan dirancang untuk keselamatan serta kenyamanan", nil, nil, "CHOICE", 0.75},
			{"Tersedia, dirancang untuk keselamatan dan kenyamanan, dan di beberapa bagian dilengkapi fitur ramah disabilitas", nil, nil, "CHOICE", 1.00},
		},
		// WR1
		"WR1": {
			{"<= 2%", nil, floatPtr(2.0), "<=", 0.05},
			{"> 2% - 10%", floatPtr(2.0), floatPtr(10.0), "BETWEEN", 0.25},
			{"> 10% - 20%", floatPtr(10.0), floatPtr(20.0), "BETWEEN", 0.50},
			{"> 20% - 40%", floatPtr(20.0), floatPtr(40.0), "BETWEEN", 0.75},
			{"> 40%", floatPtr(40.0), nil, ">", 1.00},
		},
		// WR2
		"WR2": {
			{"Tidak ada (program diperlukan, tetapi belum ada tindakan)", nil, nil, "CHOICE", 0.00},
			{"Program dalam tahap persiapan", nil, nil, "CHOICE", 0.25},
			{"1% - 25% air dikonservasi", nil, nil, "CHOICE", 0.50},
			{"> 25% - 50% air dikonservasi", nil, nil, "CHOICE", 0.75},
			{"> 50% air dikonservasi", nil, nil, "CHOICE", 1.00},
		},
		// WR3
		"WR3": {
			{"Tidak ada (program diperlukan, tetapi belum ada tindakan)", nil, nil, "CHOICE", 0.00},
			{"Program dalam tahap persiapan", nil, nil, "CHOICE", 0.25},
			{"1% - 25% air didaur ulang", nil, nil, "CHOICE", 0.50},
			{"> 25% - 50% air didaur ulang", nil, nil, "CHOICE", 0.75},
			{"> 50% air didaur ulang", nil, nil, "CHOICE", 1.00},
		},
		// WR4
		"WR4": {
			{"<= 20% perangkat hemat air dipasang", nil, nil, "CHOICE", 0.05},
			{"> 20% - 40% terpasang", nil, nil, "CHOICE", 0.25},
			{"> 40% - 60% terpasang", nil, nil, "CHOICE", 0.50},
			{"> 60% - 80% terpasang", nil, nil, "CHOICE", 0.75},
			{"> 80% terpasang", nil, nil, "CHOICE", 1.00},
		},
		// WR5
		"WR5": {
			{"<= 1%", nil, floatPtr(1.0), "<=", 0.05},
			{"> 1% - 25%", floatPtr(1.0), floatPtr(25.0), "BETWEEN", 0.25},
			{"> 25% - 50%", floatPtr(25.0), floatPtr(50.0), "BETWEEN", 0.50},
			{"> 50% - 75%", floatPtr(50.0), floatPtr(75.0), "BETWEEN", 0.75},
			{"> 75%", floatPtr(75.0), nil, ">", 1.00},
		},
		// WR6
		"WR6": {
			{"Kebijakan dan program pengendalian pencemaran air masih pada tahap perancangan", nil, nil, "CHOICE", 0.05},
			{"Kebijakan dan program pengendalian pencemaran air masih pada tahap konstruksi", nil, nil, "CHOICE", 0.25},
			{"Kebijakan dan program pengendalian pencemaran air masih pada tahap awal implementasi", nil, nil, "CHOICE", 0.50},
			{"Kebijakan dan program pengendalian pencemaran air sudah diimplementasikan sepenuhnya dan dipantau sesekali", nil, nil, "CHOICE", 0.75},
			{"Kebijakan dan program pengendalian pencemaran air sudah diimplementasikan sepenuhnya dan dipantau secara rutin", nil, nil, "CHOICE", 1.00},
		},
		// WS1
		"WS1": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"Program 3R dalam tahap persiapan", nil, nil, "CHOICE", 0.25},
			{"Program 3R diimplementasikan 1%-50%", nil, nil, "CHOICE", 0.50},
			{"Program 3R diimplementasikan > 50-75%", nil, nil, "CHOICE", 0.75},
			{"Program 3R diimplementasikan > 75%", nil, nil, "CHOICE", 1.00},
		},
		// WS2
		"WS2": {
			{"Tidak ada", nil, nil, "CHOICE", 0.00},
			{"1-3 program", nil, nil, "CHOICE", 0.25},
			{"4-6 program", nil, nil, "CHOICE", 0.50},
			{"7-10 program", nil, nil, "CHOICE", 0.75},
			{"Lebih dari 10 program", nil, nil, "CHOICE", 1.00},
		},
		// WS3
		"WS3": {
			{"Pembuangan terbuka", nil, nil, "CHOICE", 0.00},
			{"Sebagian (1%-35% diolah)", nil, nil, "CHOICE", 0.25},
			{"Sebagian (>35%-65% diolah)", nil, nil, "CHOICE", 0.50},
			{"Sebagian (>65%-85% diolah)", nil, nil, "CHOICE", 0.75},
			{"Menyeluruh (>85% diolah)", nil, nil, "CHOICE", 1.00},
		},
		// WS4
		"WS4": {
			{"Dibakar di area terbuka", nil, nil, "CHOICE", 0.00},
			{"Sebagian (1%-35% diolah)", nil, nil, "CHOICE", 0.25},
			{"Sebagian (> 35%-65% diolah)", nil, nil, "CHOICE", 0.50},
			{"Sebagian (> 65%-85% diolah)", nil, nil, "CHOICE", 0.75},
			{"Menyeluruh (> 85% diolah)", nil, nil, "CHOICE", 1.00},
		},
		// WS5
		"WS5": {
			{"Tidak diolah", nil, nil, "CHOICE", 0.00},
			{"Sebagian (1%-35% diolah)", nil, nil, "CHOICE", 0.25},
			{"Sebagian (> 35%-65% diolah)", nil, nil, "CHOICE", 0.50},
			{"Sebagian (> 65%-85% diolah)", nil, nil, "CHOICE", 0.75},
			{"Menyeluruh (> 85% diolah) atau kampus menghasilkan sampah beracun dalam jumlah minimal", nil, nil, "CHOICE", 1.00},
		},
		// WS6
		"WS6": {
			{"Dibuang tanpa diolah ke badan air", nil, nil, "CHOICE", 0.00},
			{"Diolah dengan pengolahan pendahuluan (penyaringan dan pemisahan)", nil, nil, "CHOICE", 0.25},
			{"Diolah dengan pengolahan primer (pengendapan dan koagulasi-flokulasi)", nil, nil, "CHOICE", 0.50},
			{"Diolah dengan pengolahan sekunder (pengolahan biologis)", nil, nil, "CHOICE", 0.75},
			{"Diolah dengan pengolahan tersier (pengolahan lanjutan untuk pemanfaatan ulang)", nil, nil, "CHOICE", 1.00},
		},
	}

	for indCode, tiers := range data {
		var indicator models.Indicator
		_ = facades.Orm().Query().Where("code = ?", indCode).First(&indicator)
		if indicator.ID == 0 {
			// Skip indicator if it doesn't exist in DB yet
			continue
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
