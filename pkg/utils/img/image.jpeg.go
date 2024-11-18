package img

import (
	"github.com/hrabit64/springnote-breezenote/config"
	_ "golang.org/x/image/webp"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path"
)

// ConvertToJpeg 이미지를 jpeg로 변환하여 저장합니다.
// 가로 이미지는 가로 길이를, 세로 이미지는 세로 길이를 맞춰 변환합니다. (비율 유지)
// 정사각형 이미지는 가로 세로 길이를 맞춰 변환합니다.
func ConvertToJpeg(imgReader io.Reader, size int, name string) error {
	img, _, err := decodeImageWithResize(imgReader, size)
	if err != nil {
		return err
	}

	err = convertJpegAndSave(img, name)
	if err != nil {
		return err
	}

	return nil
}

func convertJpegAndSave(img image.Image, name string) error {
	outputFileName := name + ".jpeg"
	outputPath := path.Join(config.GetSavePath(), outputFileName)

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return err
	}

	return nil
}
