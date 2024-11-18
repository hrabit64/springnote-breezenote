package utils

import (
	b64 "encoding/base64"
	"io"
	"strings"
)

// EncodeBase64 base64 인코딩
func DecodeBase64(encoded string) io.Reader {
	return b64.NewDecoder(b64.StdEncoding, strings.NewReader(encoded))
}

// IsBase64 base64 인코딩 여부 확인
func IsBase64(s string) bool {
	_, err := b64.StdEncoding.DecodeString(s)
	return err == nil
}
