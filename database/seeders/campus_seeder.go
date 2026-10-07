package seeders

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"
	"ui_greenmetric/app/services"
)

type CampusSeeder struct{}

// Signature The name and signature of the seeder.
func (s *CampusSeeder) Signature() string {
	return "CampusSeeder"
}

// campusProfile menyimpan nilai dasar kampus yang dipakai untuk membuat
// raw_input_data yang realistis dan saling konsisten antar indikator.
type campusProfile struct {
	LuasTotal          float64 // m2
	Populasi           float64 // mahasiswa + staf
	LuasTotalBangunan  float64 // m2
	TotalAirDikonsumsi float64 // m3/tahun
	TotalEnergi        float64 // kWh/tahun
	TotalMK            float64 // jumlah mata kuliah
	TotalDanaRiset     float64 // miliar rupiah
	TotalPublikasi     float64 // publikasi/tahun
	TotalLulusan       float64 // lulusan/tahun
	TotalAnggaran      float64 // miliar rupiah
	TotalPimpinan      float64 // jumlah pimpinan
}

// campusDef mendefinisikan 10 kampus nyam di Indonesia beserta level kualitas
// kampusnya (Quality 0-4: 0 = sangat kurang, 4 = sangat baik).
type campusDef struct {
	Code            string
	Name            string
	InstitutionType string
	Climate         string
	Setting         string
	Quality         int
	AdminName       string
}

var campusDefs = []campusDef{
	{"UI", "Universitas Indonesia", "University", "Tropical", "Urban", 4, "Andi Pratama"},
	{"ITB", "Institut Teknologi Bandung", "University", "Tropical", "Urban", 4, "Sari Nurhaliza"},
	{"UGM", "Universitas Gadjah Mada", "University", "Tropical", "Suburban", 3, "Budi Santoso"},
	{"UB", "Universitas Brawijaya", "University", "Tropical", "Urban", 3, "Rina Kartika"},
	{"UNDIP", "Universitas Diponegoro", "University", "Tropical", "Urban", 2, "Agus Salim"},
	{"UM", "Universitas Negeri Malang", "University", "Tropical", "Suburban", 2, "Dewi Lestari"},
	{"POLINEMA", "Politeknik Negeri Malang", "Vocational", "Tropical", "Urban", 1, "Joko Susanto"},
	{"UMM", "Universitas Muhammadiyah Malang", "University", "Tropical", "Urban", 2, "Fitri Handayani"},
	{"UNRAM", "Universitas Mataram", "University", "Tropical", "Suburban", 1, "Made Sukarta"},
	{"UNCEN", "Universitas Cenderawasih", "University", "Tropical", "Rural", 0, "Petrus Wamafma"},
}

var firstNames = []string{
	"Rizki", "Siti", "Bagas", "Ayu", "Wawan", "Nabila", "Yoga", "Intan",
	"Fajar", "Lina", "Dimas", "Salsa", "Bayu", "Tria", "Galih", "Nanda",
	"Aditya", "Putri", "Reza", "Dinda", "Hendra", "Kartika", "Yusuf", "Sri",
	"Agung", "Ratna", "Bima", "Citra", "Dian", "Eko", "Farah", "Gilang",
	"Irfan", "Maya", "Oscar", "Prilia", "Qori", "Rama", "Teguh", "Umar",
	"Vino", "Winda", "Xavi", "Yasmine", "Zaki", "Alan", "Budi", "Cahya",
	"Dewa", "Endah", "Fanny", "Gerry", "Hafiz", "Ika", "Joko", "Kirana",
	"Lukman", "Mira", "Nia", "Olie", "Prita", "Rian", "Surya", "Tono",
}

var lastNames = []string{
	"Hidayat", "Aminah", "Pratama", "Wulandari", "Setiawan", "Zahra", "Nugraha", "Permata",
	"Saputra", "Maharani", "Kusuma", "Rahmawati", "Wibowo", "Anggraini", "Santoso", "Lestari",
	"Pranata", "Utami", "Maulana", "Handayani", "Wijaya", "Susanti", "Permana", "Febrianti",
	"Kurnia", "Marlina", "Nurdin", "Oktaviani", "Pangestu", "Qurbani", "Ramadhan", "Suryani",
	"Tirta", "Ulfa", "Verdiansyah", "Wardani", "Xenovia", "Yulianti", "Zainuddin", "Adipura",
	"Brantas", "Cahyadi", "Daniswara", "Erlangga", "Firmansyah", "Gumilar", "Hartono", "Iskandar",
	"Jailani", "Kuswara", "Laksono", "Manggala", "Nugroho", "Oktriani", "Prajita", "Rusmana",
	"Setiabudi", "Tarigan", "Utomo", "Vitalis", "Wibisono", "Yulianto", "Zakaria", "Atmadja",
}

var categoryCodes = []string{"SI", "EC", "WS", "WR", "TR", "ED", "GD"}

// numericGenerators membangun raw_input_data per indikator numerik agar hasil
// rumusnya jatuh pada nilai target (tier yang diinginkan).
var numericGenerators = map[string]func(p *campusProfile, v float64) map[string]any{
	"SI1": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"luas_total": p.LuasTotal, "luas_dasar": round2(p.LuasTotal * (100 - v) / 100)}
	},
	"SI2": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"luas_total": p.LuasTotal, "luas_hutan": round2(p.LuasTotal * v / 100)}
	},
	"SI3": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"luas_total": p.LuasTotal, "luas_vegetasi": round2(p.LuasTotal * v / 100)}
	},
	"SI4": func(p *campusProfile, v float64) map[string]any {
		dasar := p.LuasTotal - v*p.Populasi
		if dasar < p.LuasTotal*0.02 {
			dasar = p.LuasTotal * 0.02
		}
		return map[string]any{"luas_total": p.LuasTotal, "luas_dasar": round2(dasar), "populasi": p.Populasi}
	},
	"EC1": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"persentase_alat_hemat_energi": v}
	},
	"EC2": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"luas_smart_building": round2(p.LuasTotalBangunan * v / 100), "luas_total_bangunan": p.LuasTotalBangunan}
	},
	"EC4": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"total_listrik": round2(v * p.Populasi), "populasi": p.Populasi}
	},
	"EC5": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"produksi_energi_terbarukan": round2(p.TotalEnergi * v / 100), "total_penggunaan_energi": p.TotalEnergi}
	},
	"EC8": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"jejak_karbon": round2(v * p.Populasi), "populasi": p.Populasi}
	},
	"WR1": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"area_resapan": round2(p.LuasTotal * v / 100), "luas_total": p.LuasTotal}
	},
	"WR5": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"air_olahan_dikonsumsi": round2(p.TotalAirDikonsumsi * v / 100), "total_air_dikonsumsi": p.TotalAirDikonsumsi}
	},
	"TR1": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"total_kendaraan": round2(v * p.Populasi), "populasi": p.Populasi}
	},
	"TR4": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"total_zev": round2(v * p.Populasi), "populasi": p.Populasi}
	},
	"TR5": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"luas_parkir": round2(p.LuasTotal * v / 100), "luas_total": p.LuasTotal}
	},
	"ED1": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"mk_keberlanjutan": round2(p.TotalMK * v / 100), "total_mk": p.TotalMK}
	},
	"ED2": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"dana_riset_keberlanjutan": round2(p.TotalDanaRiset * v / 100), "total_dana_riset": p.TotalDanaRiset}
	},
	"ED3": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"publikasi_keberlanjutan": round2(p.TotalPublikasi * v / 100), "total_publikasi": p.TotalPublikasi}
	},
	"ED10": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"lulusan_green_jobs": round2(p.TotalLulusan * v / 100), "total_lulusan": p.TotalLulusan}
	},
	"GD1": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"anggaran_keberlanjutan": round2(p.TotalAnggaran * v / 100), "total_anggaran": p.TotalAnggaran}
	},
	"GD8": func(p *campusProfile, v float64) map[string]any {
		return map[string]any{"pimpinan_perempuan": round2(p.TotalPimpinan * v / 100), "total_pimpinan": p.TotalPimpinan}
	},
}

// Run executes the seeder logic.
func (s *CampusSeeder) Run() error {
	hashedPassword, err := facades.Hash().Make("password")
	if err != nil {
		return err
	}

	// Muat master data sekali di awal
	var categories []models.Category
	if err := facades.Orm().Query().OrderBy("id").Get(&categories); err != nil {
		return err
	}
	catOrder := make(map[string]int)
	for i, c := range categories {
		catOrder[c.Code] = i
	}

	var indicators []models.Indicator
	if err := facades.Orm().Query().OrderBy("category_id, id").Get(&indicators); err != nil {
		return err
	}
	catIDToCode := make(map[uint]string)
	for _, c := range categories {
		catIDToCode[c.ID] = c.Code
	}

	var allTiers []models.IndicatorScoringTier
	if err := facades.Orm().Query().Get(&allTiers); err != nil {
		return err
	}
	tiersByIndicator := make(map[uint][]models.IndicatorScoringTier)
	for _, t := range allTiers {
		tiersByIndicator[t.IndicatorID] = append(tiersByIndicator[t.IndicatorID], t)
	}
	for k, ts := range tiersByIndicator {
		sort.Slice(ts, func(i, j int) bool {
			return ts[i].PointMultiplier < ts[j].PointMultiplier
		})
		tiersByIndicator[k] = ts
	}

	scoring := services.NewScoringService()
	currentYear := time.Now().Year()

	for ci, def := range campusDefs {
		// 1. Buat kampus (idempoten berdasarkan code)
		var campus models.Campus
		cnt, err := facades.Orm().Query().Model(&models.Campus{}).Where("code = ?", def.Code).Count()
		if err != nil {
			return err
		}
		if cnt == 0 {
			campus = models.Campus{
				Code:            def.Code,
				Name:            def.Name,
				InstitutionType: def.InstitutionType,
				Climate:         def.Climate,
				Setting:         def.Setting,
			}
			if err := facades.Orm().Query().Create(&campus); err != nil {
				return err
			}
		} else {
			if err := facades.Orm().Query().Where("code = ?", def.Code).First(&campus); err != nil {
				return err
			}
		}

		// 2. Buat user per kampus (1 admin kampus + 7 operator, satu per kategori)
		lowerCode := strings.ToLower(def.Code)
		users := []models.User{
			{
				CampusID: campus.ID,
				Name:     def.AdminName,
				Email:    fmt.Sprintf("admin.%s@gmail.com", lowerCode),
				Password: hashedPassword,
				Role:     "ADMIN_KAMPUS",
			},
		}
		for oi, cat := range categoryCodes {
			// Kombinasikan pool nama depan & belakang agar 70 operator semua unik
			first := firstNames[(ci*3+oi*7)%len(firstNames)]
			last := lastNames[(ci*5+oi*11)%len(lastNames)]
			name := fmt.Sprintf("%s %s", first, last)
			users = append(users, models.User{
				CampusID: campus.ID,
				Name:     name,
				Email:    fmt.Sprintf("operator.%s.%s@gmail.com", strings.ToLower(cat), lowerCode),
				Password: hashedPassword,
				Role:     "OPERATOR_" + cat,
			})
		}
		for _, u := range users {
			exists, err := facades.Orm().Query().Model(&models.User{}).Where("email = ?", u.Email).Count()
			if err != nil {
				return err
			}
			if exists == 0 {
				if err := facades.Orm().Query().Create(&u); err != nil {
					return err
				}
			}
		}

		// 3. Assessment tahun berjalan + jawaban seluruh indikator
		var assessment models.CampusAssessment
		_ = facades.Orm().Query().Where("campus_id = ? AND assessment_year = ?", campus.ID, currentYear).First(&assessment)
		if assessment.ID == 0 {
			assessment = models.CampusAssessment{
				CampusID:       campus.ID,
				AssessmentYear: currentYear,
				OverallScore:   0,
				Status:         "DRAFT",
			}
			if err := facades.Orm().Query().Create(&assessment); err != nil {
				return err
			}
		}

		answerCnt, err := facades.Orm().Query().Model(&models.AssessmentAnswer{}).
			Where("campus_assessment_id = ?", assessment.ID).Count()
		if err != nil {
			return err
		}
		if answerCnt == 0 {
			profile := buildProfile(ci, def.Quality)
			for ii, ind := range indicators {
				tiers := tiersByIndicator[ind.ID]
				if len(tiers) == 0 {
					continue
				}
				catIdx := catOrder[catIDToCode[ind.CategoryID]]
				idx := tierIndexFor(ci, def.Quality, catIdx, ii, len(tiers))
				raw := buildRawInput(&ind, tiers[idx], &profile)
				if _, err := scoring.Calculate(assessment.ID, ind.Code, raw); err != nil {
					return fmt.Errorf("gagal menghitung indikator %s kampus %s: %w", ind.Code, def.Code, err)
				}
			}
			// Ambil overall_score terbaru setelah seluruh jawaban dihitung
			_ = facades.Orm().Query().Where("id = ?", assessment.ID).First(&assessment)
		}

		// 4. Assessment 2 tahun sebelumnya untuk grafik tren (tanpa jawaban detail)
		for yi, year := range []int{currentYear - 2, currentYear - 1} {
			var factor float64
			if def.Quality <= 1 {
				// Kampus yang mundur: skor masa lalu lebih tinggi dari sekarang
				factor = 1.05 + 0.20*hash01(ci, yi+1)
				if yi == 0 {
					factor *= 1.05 + 0.15*hash01(ci, yi+50)
				}
			} else {
				// Kampus yang membaik: skor masa lalu lebih rendah
				factor = 0.72 + 0.23*hash01(ci, yi+1)
				if yi == 0 {
					factor *= 0.72 + 0.20*hash01(ci, yi+50)
				}
			}
			score := round2(assessment.OverallScore * factor)
			if score > 10000 {
				score = 10000
			}
			status := "SUBMITTED"
			if def.Quality >= 3 {
				status = "VERIFIED"
			}
			exists, err := facades.Orm().Query().Model(&models.CampusAssessment{}).
				Where("campus_id = ? AND assessment_year = ?", campus.ID, year).Count()
			if err != nil {
				return err
			}
			if exists == 0 {
				if err := facades.Orm().Query().Create(&models.CampusAssessment{
					CampusID:       campus.ID,
					AssessmentYear: year,
					OverallScore:   score,
					Status:         status,
				}); err != nil {
					return err
				}
			}
		}

		// 5. Kampus yang rapi sudah submit assessment tahun berjalan
		if def.Quality >= 3 && assessment.Status == "DRAFT" {
			assessment.Status = "SUBMITTED"
			if err := facades.Orm().Query().Where("id = ?", assessment.ID).Save(&assessment); err != nil {
				return err
			}
		}
	}

	return nil
}

// buildProfile menyusun profil kampus yang konsisten dan bervariasi. Kampus
// berkualitas baik cenderung berlahan luas dengan populasi lebih kecil.
func buildProfile(ci int, quality int) campusProfile {
	r := hash01(ci, 99)
	luas := 120000 + float64(quality)*55000 + r*150000
	pop := 48000 - float64(quality)*4000 - r*15000
	if pop < 3500 {
		pop = 3500
	}
	return campusProfile{
		LuasTotal:          round2(luas),
		Populasi:           round2(pop),
		LuasTotalBangunan:  round2(luas * (0.10 + 0.05*r)),
		TotalAirDikonsumsi: round2(pop * (80 + 40*r)),
		TotalEnergi:        round2(pop * (900 + 300*r)),
		TotalMK:            round2(400 + 200*r + float64(quality)*40),
		TotalDanaRiset:     round2(50 + 20*r + float64(quality)*8),
		TotalPublikasi:     round2(300 + 150*r + float64(quality)*40),
		TotalLulusan:       round2(pop * 0.18),
		TotalAnggaran:      round2(800 + 300*r + float64(quality)*60),
		TotalPimpinan:      round2(60 + 40*r),
	}
}

// buildRawInput menyusun raw_input_data sesuai tipe indikator dan tier target.
func buildRawInput(ind *models.Indicator, tier models.IndicatorScoringTier, p *campusProfile) map[string]any {
	if ind.InputType == "SINGLE_CHOICE" {
		return map[string]any{"option_label": tier.OptionLabel}
	}
	gen, ok := numericGenerators[ind.Code]
	if !ok {
		return map[string]any{"value": tierValue(tier)}
	}
	return gen(p, tierValue(tier))
}

// tierValue mengambil nilai tengah tier yang aman (selalu di dalam rentang tier).
func tierValue(t models.IndicatorScoringTier) float64 {
	switch {
	case t.MinValue != nil && t.MaxValue != nil:
		return (*t.MinValue + *t.MaxValue) / 2
	case t.MinValue != nil:
		return *t.MinValue * 1.1
	default:
		return *t.MaxValue * 0.5
	}
}

// tierIndexFor memilih index tier berdasarkan kualitas kampus dengan variasi
// per kategori & indikator agar tidak semua indikator dapat nilai sama.
func tierIndexFor(ci, quality, catIdx, indIdx, nTiers int) int {
	jitter := ((ci*5 + catIdx*3 + indIdx*7) % 3) - 1
	idx := quality + jitter
	if idx < 0 {
		idx = 0
	}
	if idx >= nTiers {
		idx = nTiers - 1
	}
	return idx
}

// hash01 menghasilkan nilai pseudo-acak deterministik 0-1 dari dua integer,
// sehingga seeder menghasilkan data yang sama setiap dijalankan.
func hash01(a, b int) float64 {
	x := uint32((a+1)*2654435761 ^ (b+1)*40503)
	x ^= x >> 13
	x *= 0x5bd1e995
	x ^= x >> 15
	return float64(x%10000) / 10000.0
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
