package img

import (
	"encoding/base64"
	"github.com/hrabit64/springnote-breezenote/config"
	"image/gif"
	"os"
	"path"
)

// SaveGif base64로 인코딩된 gif 이미지를 저장합니다.
func SaveGif(base64gif, name string) error {
	decodedGif, err := base64.StdEncoding.DecodeString(base64gif)
	if err != nil {
		return err
	}

	outputFileName := name + ".gif"
	outputPath := path.Join(config.GetSavePath(), outputFileName)

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(decodedGif)
	if err != nil {
		return err
	}

	return nil

}

func LoadGif(fileName string) (*gif.GIF, error) {
	targetDir := path.Join(config.GetSavePath(), fileName)

	file, err := os.Open(targetDir)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	gifImg, err := gif.DecodeAll(file)
	if err != nil {
		return nil, err
	}

	return gifImg, nil
}

//
////func ConvertGif(imgReader io.Reader, size int, name string) error {
////	gifImg, err := gif.DecodeAll(imgReader)
////	if err != nil {
////		return err
////	}
////
////	width, height := calImageSize(gifImg.Image[0].Bounds().Dx(), gifImg.Image[0].Bounds().Dy(), size)
////
////	gifImg.Config.Width = width
////	gifImg.Config.Height = height
////
////	RGBAImg := image.NewRGBA(image.Rect(0, 0, gifImg.Image[0].Bounds().Dx(), gifImg.Image[0].Bounds().Dy()))
////
////	newImages := make([]*image.Paletted, len(gifImg.Image))
////	for i, frame := range gifImg.Image {
////		bound := frame.Bounds()
////		draw.Draw(RGBAImg, bound, frame, bound.Min, draw.Over)
////		resizedFrame := resize.Resize(uint(width), uint(height), frame, resize.Lanczos3)
////
////		bound = resizedFrame.Bounds()
////		palettedImage := image.NewPaletted(bound, palette.Plan9)
////		draw.FloydSteinberg.Draw(palettedImage, bound, resizedFrame, image.ZP)
////
////		newImages[i] = palettedImage
////	}
////	newGif := &gif.GIF{
////		Image:    newImages,
////		Delay:    gifImg.Delay,
////		Disposal: gifImg.Disposal,
////		Config: image.Config{
////			ColorModel: gifImg.Config.ColorModel,
////			Width:      width,
////			Height:     height,
////		},
////		LoopCount: gifImg.LoopCount,
////	}
////	err = convertGifAndSave(newGif, name)
////	if err != nil {
////		return err
////	}
////
////	return nil
////}
////
////
////// ref : https://github.com/elsonwu/resizegif.go
////func resizeGif(img *gif.GIF, width int, height int) *gif.GIF {
////	img.Config.Width = width
////	img.Config.Height = height
////
////	firstFrame := img.Image[0].Bounds()
////	rgbaImg := image.NewRGBA(image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy()))
////
////	for index, frame := range img.Image {
////		b := frame.Bounds()
////		draw.Draw(rgbaImg, b, frame, b.Min, draw.Over)
////		img.Image[index] = imageToPaletted(resize.Resize(uint(width), uint(height), rgbaImg, resize.Lanczos3))
////	}
////
////	return img
////}
