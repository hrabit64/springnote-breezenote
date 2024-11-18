package imageTestUtil

import (
	b64 "encoding/base64"
	"github.com/corona10/goimagehash"
	"github.com/hrabit64/springnote-breezenote/pkg/utils"
	"image"
	"io"
	"os"
	"path"
)

// IsSame 은 두 이미지 해시가 같은지 비교합니다.
// distance  < 5 : same image(or resize), < 10 : similar image(like crop)
func IsSame(h1 *goimagehash.ImageHash, h2 *goimagehash.ImageHash, distance int) (bool, error) {
	d, err := h1.Distance(h2)
	if err != nil {
		return false, err
	}

	return d < distance, nil
}

// GetTestImageHash 테스트 이미지의 해시값을 가져옵니다.
func GetTestImageHash(name string) (*goimagehash.ImageHash, error) {

	file, err := os.Open(path.Join(utils.GetRootPath(), "test", "data", name))
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return GetImageHashFromReader(file)
}

func GetImageHashFromReader(reader io.Reader) (*goimagehash.ImageHash, error) {
	image, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	hash, err := goimagehash.AverageHash(image)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

// GetTestImageBase64 테스트 이미지의 base64 값을 가져옵니다.
// 테스트 이미지는 test/data 디렉토리에 있어야 합니다.
func GetTestImageBase64(name string) (string, error) {

	targetDir := path.Join(utils.GetRootPath(), "test", "data", name)

	file, err := os.Open(targetDir)
	if err != nil {
		return "", err
	}

	defer file.Close()

	testImgByte, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	// encode to base64
	imgBase64 := b64.StdEncoding.EncodeToString(testImgByte)
	return imgBase64, nil
}

// GetTestImageHashFromOtherDir 다른 디렉토리에 있는 테스트 이미지의 해시값을 가져옵니다.
func GetTestImageHashFromOtherDir(name, dir string) (*goimagehash.ImageHash, error) {
	targetDir := path.Join(utils.GetRootPath(), dir, name)

	file, err := os.Open(targetDir)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return GetImageHashFromReader(file)

}
