package img

import (
	"errors"
	"github.com/hrabit64/springnote-breezenote/config"
	"github.com/nfnt/resize"
	"golang.org/x/image/webp"
	_ "golang.org/x/image/webp"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path"
	"strings"
)

// LoadImageWithResize 이미지를 리사이즈하여 반환합니다.
// width, height가 0이면 원본 이미지를 반환합니다.
// 주어진 값이 원본 이미지의 비율과 같으면 리사이즈만 수행합니다.
// 주어진 값이 원본 이미지의 비율과 다르면 리사이즈 후 중앙을 기준으로 자릅니다.
// 만약 둘중 하나가 0이면 주어진 값에 맞춰 비율을 유지하며 리사이즈합니다.
func LoadImageWithResize(fileName string, width int, height int) (image.Image, error) {
	targetDir := path.Join(config.GetSavePath(), fileName)

	file, err := os.Open(targetDir)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	img, err := decodeImage(file, strings.Split(fileName, ".")[1])
	if err != nil {
		return nil, err
	}
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	// case 1. given origin size
	if width == 0 && height == 0 {
		return img, nil
	}

	// case 2. ratio is not same
	if !isSameRatio(imgWidth, imgHeight, width, height) {
		var resizedImg image.Image
		if getRatio(imgWidth, imgHeight) < getRatio(width, height) {
			resizedImg = resizeImage(img, width, 0)
		} else {
			resizedImg = resizeImage(img, 0, height)
		}

		cropImg := CropCenterImage(resizedImg, width, height)
		return cropImg, nil
	}

	// case 3. ratio is same
	resizedImg := resizeImage(img, width, height)
	return resizedImg, nil
}

// decodeImage 이미지를 디코딩합니다.
func decodeImage(file io.Reader, format string) (image.Image, error) {
	switch format {
	case "webp":
		img, err := webp.Decode(file)
		if err != nil {
			return nil, err
		}
		return img, nil
	case "jpeg":
		img, err := jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
		return img, nil
	default:
		return nil, errors.New("not supported format")
	}
}

// resizeImage 이미지를 리사이즈합니다.
func resizeImage(img image.Image, width int, height int) image.Image {
	return resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
}

func thumbnailImage(img image.Image, width int, height int) image.Image {
	return resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)
}

func decodeImageWithResize(imgReader io.Reader, size int) (image.Image, string, error) {
	img, format, err := image.Decode(imgReader)
	if err != nil {
		return nil, "", err
	}

	// raise error if format is gif
	if format == "gif" {
		return nil, "", err
	}

	img = thumbnailImage(img, size, size)
	return img, format, nil
}

func isSameRatio(imgWidth, imgHeight, width, height int) bool {
	return width == 0 || height == 0 || imgWidth*height == imgHeight*width
}

// getRatio 이미지의 비율을 반환합니다.
func getRatio(imgWidth, imgHeight int) float64 {
	return float64(imgWidth) / float64(imgHeight)
}

// CropCenterImage 이미지를 중앙을 기준으로 자릅니다.
func CropCenterImage(img image.Image, width, height int) image.Image {
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	// Return the original image if cropping is unnecessary
	if imgWidth == width && imgHeight == height {
		return img
	}

	// Ensure the target dimensions are smaller than the original dimensions
	if imgWidth < width || imgHeight < height {
		return img
	}

	x := (imgWidth - width) / 2
	y := (imgHeight - height) / 2

	// Define the cropping rectangle
	cropRect := image.Rect(0, 0, width, height)

	// Create a new RGBA image with the desired dimensions
	rgbaImg := image.NewRGBA(cropRect)

	// Offset the source image to crop from the center area
	draw.Draw(rgbaImg, cropRect, img, image.Point{X: x, Y: y}, draw.Src)

	return rgbaImg
}
