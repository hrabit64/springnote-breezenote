package img

import (
	"github.com/chai2010/webp"
	"github.com/hrabit64/springnote-breezenote/config"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path"
)

// ConvertToWebp 이미지를 webp로 변환하여 저장합니다.
// 가로 이미지는 가로 길이를, 세로 이미지는 세로 길이를 맞춰 변환합니다. (비율 유지)
// 정사각형 이미지는 가로 세로 길이를 맞춰 변환합니다.
func ConvertToWebp(imgReader io.Reader, size int, name string) error {

	img, _, err := decodeImageWithResize(imgReader, size)
	if err != nil {
		return err
	}

	err = convertWebpAndSave(img, name)
	if err != nil {
		return err
	}

	return nil
}

func convertWebpAndSave(img image.Image, name string) error {
	outputFileName := name + ".webp"
	outputPath := path.Join(config.GetSavePath(), outputFileName)
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	err = webp.Encode(file, img, &webp.Options{Lossless: true})
	if err != nil {
		return err
	}

	return nil
}
