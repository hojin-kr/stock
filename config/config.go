package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	AllowedOrigins []string `json:"allowed_origins"`
}

var (
	config *Config
	once   sync.Once
)

// GetConfig는 싱글톤 패턴으로 설정을 반환합니다.
func GetConfig() *Config {
	once.Do(func() {
		config = &Config{
			AllowedOrigins: []string{
				"http://localhost:3000",       // React 개발 서버
				"http://localhost:3001",       // React 개발 서버
				"http://localhost:8080",       // 백엔드 서버
				"http://127.0.0.1:3000",       // React 개발 서버 (IP)
				"http://127.0.0.1:3001",       // React 개발 서버 (IP)
				"http://127.0.0.1:8080",       // 백엔드 서버 (IP)
				"http://192.168.200.150:3001", // React 개발 서버 (IP)
			},
		}
		loadConfig()
	})
	return config
}

// loadConfig는 config.json 파일에서 설정을 로드합니다.
func loadConfig() {
	file, err := os.Open("config/config.json")
	if err != nil {
		// config.json 파일이 없으면 기본 설정을 사용
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		// JSON 파싱 에러가 발생하면 기본 설정을 사용
		return
	}
}

// IsOriginAllowed는 주어진 origin이 허용된 도메인인지 확인합니다.
func (c *Config) IsOriginAllowed(origin string) bool {
	for _, allowed := range c.AllowedOrigins {
		if allowed == origin {
			return true
		}
	}
	return false
}
