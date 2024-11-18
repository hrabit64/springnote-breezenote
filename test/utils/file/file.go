package fileTestUtil

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/hrabit64/springnote-breezenote/pkg/utils"
	"io"
	"os"
	"path"
)

func getFileHash(filePath string) (string, error) {
	file, err := os.Open(path.Join(utils.GetRootPath(), filePath))
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// CompareFiles 두 파일의 해시값을 비교하여 같은지 여부를 반환합니다.
func CompareFiles(filePath1, filePath2 string) (bool, error) {
	hash1, err := getFileHash(filePath1)
	if err != nil {
		return false, err
	}

	hash2, err := getFileHash(filePath2)
	if err != nil {
		return false, err
	}

	return hash1 == hash2, nil
}
