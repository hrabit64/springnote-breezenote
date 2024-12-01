package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetRootPath 프로젝트 루트 경로를 반환합니다.
func GetRootPath() string {

	envPath := os.Getenv("BREEZENOTE_ROOT_PATH")

	if envPath != "" {
		return envPath
	}

	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "../..")
}
