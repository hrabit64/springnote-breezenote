package utils

import "strings"

// GetFileExt 파일 확장자 추출
func GetFileExt(filename string) string {
	split := strings.Split(filename, ".")
	return split[len(split)-1]
}

// GetFileName 파일 이름 추출
func GetFileName(filename string) string {
	split := strings.Split(filename, ".")
	return split[0]
}
