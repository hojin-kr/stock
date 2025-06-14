package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"stock-backend/config"
)

type Report struct {
	Date     string `json:"date"`
	Name     string `json:"name"`
	Decision string `json:"decision"`
	Locale   string `json:"locale"`
	Content  string `json:"content"`
	Type     string `json:"type"`
}

// CORS 미들웨어 함수
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		cfg := config.GetConfig()

		// Origin이 허용된 도메인인 경우에만 CORS 헤더 설정
		if cfg.IsOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// OPTIONS 요청 처리
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 다음 핸들러 실행
		next(w, r)
	}
}

func getLocaleFromPath(path string) string {
	if strings.Contains(path, "/EN/") {
		return "EN"
	} else if strings.Contains(path, "/KO/") {
		return "KO"
	}
	return "EN" // 기본값
}

func getDecisionFromPath(path string) string {
	parts := strings.Split(path, string(os.PathSeparator))
	for _, part := range parts {
		if part == "BUY" || part == "HOLD" || part == "SELL" {
			return part
		}
	}
	return ""
}

func getReportType(filename string) string {
	switch filename {
	case "complete.md":
		return "complete"
	case "final_trade_decision.md":
		return "final_decision"
	case "investment_plan.md":
		return "investment_plan"
	case "market_report.md":
		return "market_report"
	case "trader_investment_plan.md":
		return "trader_plan"
	default:
		return "unknown"
	}
}

func getAllReports(w http.ResponseWriter, r *http.Request) {
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "EN"
	}

	reportType := r.URL.Query().Get("type")
	if reportType == "" {
		reportType = "complete" // 기본값
	}

	dir := "./reports"
	var reports []Report

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		fileLocale := getLocaleFromPath(path)
		if fileLocale != locale {
			return nil
		}

		parts := strings.Split(path, string(os.PathSeparator))
		if len(parts) < 4 {
			return nil
		}

		filename := filepath.Base(path)
		if getReportType(filename) != reportType {
			return nil
		}

		date := parts[1]
		decision := getDecisionFromPath(path)
		name := parts[len(parts)-2]
		if name == "EN" || name == "KO" {
			name = parts[len(parts)-3]
		}

		content, _ := ioutil.ReadFile(path)
		reports = append(reports, Report{
			Date:     date,
			Name:     name,
			Decision: decision,
			Locale:   fileLocale,
			Content:  string(content),
			Type:     getReportType(filename),
		})
		return nil
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func getReportsByDate(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	locale := r.URL.Query().Get("locale")
	reportType := r.URL.Query().Get("type")

	if locale == "" {
		locale = "EN"
	}
	if reportType == "" {
		reportType = "complete"
	}

	if date == "" {
		http.Error(w, "Date parameter is required", http.StatusBadRequest)
		return
	}

	dir := fmt.Sprintf("./reports/%s", date)
	var reports []Report

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		fileLocale := getLocaleFromPath(path)
		if fileLocale != locale {
			return nil
		}

		filename := filepath.Base(path)
		if getReportType(filename) != reportType {
			return nil
		}

		parts := strings.Split(path, string(os.PathSeparator))
		if len(parts) < 4 {
			return nil
		}

		decision := getDecisionFromPath(path)
		name := parts[len(parts)-2]
		if name == "EN" || name == "KO" {
			name = parts[len(parts)-3]
		}

		content, _ := ioutil.ReadFile(path)
		reports = append(reports, Report{
			Date:     date,
			Name:     name,
			Decision: decision,
			Locale:   fileLocale,
			Content:  string(content),
			Type:     getReportType(filename),
		})
		return nil
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func getReportsByDecision(w http.ResponseWriter, r *http.Request) {
	decision := r.URL.Query().Get("decision")
	locale := r.URL.Query().Get("locale")
	reportType := r.URL.Query().Get("type")

	if locale == "" {
		locale = "EN"
	}
	if reportType == "" {
		reportType = "complete"
	}

	if decision == "" {
		http.Error(w, "Decision parameter is required", http.StatusBadRequest)
		return
	}

	dir := "./reports"
	var reports []Report

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		fileLocale := getLocaleFromPath(path)
		if fileLocale != locale {
			return nil
		}

		filename := filepath.Base(path)
		if getReportType(filename) != reportType {
			return nil
		}

		parts := strings.Split(path, string(os.PathSeparator))
		if len(parts) < 4 {
			return nil
		}

		fileDecision := getDecisionFromPath(path)
		if fileDecision != decision {
			return nil
		}

		date := parts[1]
		name := parts[len(parts)-2]
		if name == "EN" || name == "KO" {
			name = parts[len(parts)-3]
		}

		content, _ := ioutil.ReadFile(path)
		reports = append(reports, Report{
			Date:     date,
			Name:     name,
			Decision: decision,
			Locale:   fileLocale,
			Content:  string(content),
			Type:     getReportType(filename),
		})
		return nil
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func main() {
	fmt.Println("Server is starting...")

	// 정적 파일 서빙
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// API 엔드포인트
	http.HandleFunc("/api/reports", enableCORS(getAllReports))
	http.HandleFunc("/api/reports/date", enableCORS(getReportsByDate))
	http.HandleFunc("/api/reports/decision", enableCORS(getReportsByDecision))

	// 서버 시작
	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
