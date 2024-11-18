package img

import (
	"github.com/corona10/goimagehash"
	"github.com/franela/goblin"
	"github.com/hrabit64/springnote-breezenote/config"
	"github.com/hrabit64/springnote-breezenote/pkg/utils"
	dirTestUtil "github.com/hrabit64/springnote-breezenote/test/utils/dir"
	imageTestUtil "github.com/hrabit64/springnote-breezenote/test/utils/image"
	"image/jpeg"
	"os"
	"path"
	"testing"
)

func TestConvertToJpeg(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("올바른 이미지가 주어지면, Jpeg로 리사이징 및 변환에 성공한다.", func() {
		testImgs := []struct {
			imgName        string
			expectedWidth  int
			expectedHeight int
			description    string
		}{
			{"427×640.jpg", 133, 200, "세로 형태의 이미지가 주어지면, 세로 길이가 주어진 길이로 변환되고, 가로 길이는 비율에 맞춰 변환되어 저장 된다."},
			{"500×333.jpg", 200, 133, "가로 형태의 이미지가 주어지면, 가로 길이가 주어진 길이로 변환되고, 세로 길이는 비율에 맞춰 변환되어 저장 된다."},
			{"511×511.jpg", 200, 200, "정사각형 형태의 이미지가 주어지면, 가로 세로 길이가 주어진 길이로 변환되어 저장된다."},
			{"378×504.webp", 150, 200, "webp 확장자의 이미지가 주어지면, 정상적으로 변환되어 저장된다."},
			{"378×504.png", 150, 200, "png 확장자의 이미지가 주어지면, 정상적으로 변환되어 저장된다."},
		}

		for _, item := range testImgs {
			//fmt.Sprintf("%s 이미지가 주어지면, %d x %d 로 변환된다.")
			g.It(item.description, func() {
				runTestConvertToJpeg(g, item.imgName, item.expectedWidth, item.expectedHeight)
			})
		}
	})

}

func runTestConvertToJpeg(g *goblin.G, imageName string, exceptedWidth int, exceptedHeight int) {
	err := config.SetupConfig()
	if err != nil {
		g.Fatalf("config.SetupConfig() error : %v", err)
		g.FailNow()
	}

	tmpDir, err := dirTestUtil.CreateTmpDir()
	if err != nil {
		g.Fatalf("dirTestUtil.CreateTmpDir() error : %v", err)
		g.FailNow()
	}
	defer func() {
		err := dirTestUtil.RemoveTmpDir(tmpDir)
		if err != nil {
			g.Fatalf("dirTestUtil.RemoveTmpDir() error : %v", err)
			g.FailNow()
		}
	}()

	config.RootConfig.SavePath = path.Join("tmp", tmpDir)

	targetDir := path.Join(utils.GetRootPath(), "test", "data", imageName)

	imgReader, err := os.Open(targetDir)
	if err != nil {
		g.Fatalf("Open Test Image error : %v", err)
		g.FailNow()
	}
	defer imgReader.Close()

	err = ConvertToJpeg(imgReader, 200, "test")

	resFile, err := os.Open(path.Join(config.GetSavePath(), "test.jpeg"))
	g.Assert(err).IsNil(err)
	defer resFile.Close()

	resImage, err := jpeg.Decode(resFile)
	g.Assert(err).IsNil(err)
	g.Assert(resImage.Bounds().Dx()).Equal(exceptedWidth)
	g.Assert(resImage.Bounds().Dy()).Equal(exceptedHeight)

	originImgHash, err := imageTestUtil.GetTestImageHash(imageName)
	if err != nil {
		g.Fatalf("imageTestUtil.GetTestImageHash() error : %v", err)
		g.FailNow()
	}

	resHash, err := goimagehash.AverageHash(resImage)
	g.Assert(err).IsNil(err)

	isSame, err := imageTestUtil.IsSame(originImgHash, resHash, 5)
	g.Assert(err).IsNil(err)
	g.Assert(isSame).IsTrue()

}
