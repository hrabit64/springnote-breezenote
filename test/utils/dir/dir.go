package dirTestUtil

import (
	"github.com/google/uuid"
	"github.com/hrabit64/springnote-breezenote/pkg/utils"
	"os"
	"path"
)

// CreateTmpDir 임시 디렉토리를 생성하여 디렉토리 이름을 반환합니다.
// 경로는 root/tmp/uuid 형식입니다.
func CreateTmpDir() (string, error) {
	dirName, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	tmpDir := path.Join(utils.GetRootPath(), "tmp", dirName.String())

	err = os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dirName.String(), nil
}

// RemoveTmpDir 임시 디렉토리를 삭제합니다.
// 이때 삭제할 디렉토리 이름은 CreateTmpDir 함수에서 반환한 값을 제공해야 합니다.
func RemoveTmpDir(dirName string) error {
	tmpDir := path.Join(utils.GetRootPath(), "tmp", dirName)
	err := os.RemoveAll(tmpDir)
	if err != nil {
		return err
	}
	return nil
}
