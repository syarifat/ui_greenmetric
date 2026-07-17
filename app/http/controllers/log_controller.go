package controllers

import (
	"bufio"
	"os"
	"ui_greenmetric/app/facades"

	"github.com/goravel/framework/contracts/http"
)

type LogController struct{}

func NewLogController() *LogController {
	return &LogController{}
}

// ViewLogs returns the last 500 lines of goravel.log if the secret key matches
func (r *LogController) ViewLogs(ctx http.Context) http.Response {
	secret := ctx.Request().Query("secret")
	expectedSecret := facades.Config().GetString("app.logs_secret")
	if expectedSecret == "" {
		expectedSecret = "super-secret-logs-key"
	}

	if secret == "" || secret != expectedSecret {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"status":  "error",
			"code":    http.StatusForbidden,
			"message": "Akses ditolak. Secret key tidak valid.",
		})
	}

	logFilePath := "storage/logs/goravel.log"
	file, err := os.Open(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ctx.Response().String(http.StatusOK, "File log kosong atau belum dibuat.")
		}
		return ctx.Response().String(http.StatusInternalServerError, "Gagal membuka file log: "+err.Error())
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return ctx.Response().String(http.StatusInternalServerError, "Gagal membaca file log: "+err.Error())
	}

	maxLines := 500
	start := 0
	if len(lines) > maxLines {
		start = len(lines) - maxLines
	}

	var result string
	for i := start; i < len(lines); i++ {
		result += lines[i] + "\n"
	}

	return ctx.Response().String(http.StatusOK, result)
}
