package img

import (
	"fmt"
	"github.com/corona10/goimagehash"
	goblin "github.com/franela/goblin"
	"github.com/hrabit64/springnote-breezenote/config"
	imageTestUtil "github.com/hrabit64/springnote-breezenote/test/utils/image"
	"path"
	"testing"
)

func TestLoadImageWithResize(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("올바른 Jpeg이미지가 주어지면, 변환에 성공한다.", func() {
		g.BeforeEach(func() {
			err := config.SetupConfig()
			if err != nil {
				g.Fatalf("config.SetupConfig() error : %v", err)
				g.FailNow()
			}
			config.RootConfig.SavePath = path.Join("test", "data")
		})

		testImgs := []struct {
			imgName        string
			reqWidth       int
			reqHeight      int
			expectedWidth  int
			expectedHeight int
			distance       int
			description    string
		}{
			{"427×640.jpeg", 200, 0, 200, 300, 5, "가로에 대한 길이만 주어지면, 세로 길이는 비율에 맞게 변환된다."},
			{"427×640.jpeg", 0, 200, 134, 200, 5, "세로에 대한 길이만 주어지면, 가로 길이는 비율에 맞게 변환된다."},
			{"427×640.jpeg", 214, 320, 214, 320, 5, "원본 이미지와 비율이 같은 가로 세로 길이가 주어지면, 해당 길이로 변환된다."},
			{"427×640.jpeg", 200, 200, 200, 200, 10, "원본 이미지와 비율이 다른 가로 세로 길이가 주어지면, 해당 길이 만큼 잘라내어 변환된다."},
		}

		for _, item := range testImgs {
			//fmt.Sprintf("%s 이미지가 %d x %d 사이즈로 리사이징 요청이 들어오면, %d x %d 로 변환된다.", item.imgName, item.reqWidth, item.reqHeight, item.expectedHeight, item.expectedWidth)
			g.It(item.description, func() {
				runLoadImageWithResize(g, item.imgName, item.reqWidth, item.reqHeight, item.expectedWidth, item.expectedHeight, item.distance)
			})
		}
	})

	g.Describe("올바른 Webp이미지가 주어지면, 변환에 성공한다.", func() {
		g.BeforeEach(func() {
			err := config.SetupConfig()
			if err != nil {
				g.Fatalf("config.SetupConfig() error : %v", err)
				g.FailNow()
			}
			config.RootConfig.SavePath = path.Join("test", "data")
		})

		testImgs := []struct {
			imgName        string
			reqWidth       int
			reqHeight      int
			expectedWidth  int
			expectedHeight int
			distance       int
			description    string
		}{
			{"427×640.jpeg", 200, 0, 200, 300, 5, "가로에 대한 길이만 주어지면, 세로 길이는 비율에 맞게 변환된다."},
			{"427×640.jpeg", 0, 200, 134, 200, 5, "세로에 대한 길이만 주어지면, 가로 길이는 비율에 맞게 변환된다."},
			{"427×640.jpeg", 214, 320, 214, 320, 5, "원본 이미지와 비율이 같은 가로 세로 길이가 주어지면, 해당 길이로 변환된다."},
			{"427×640.jpeg", 200, 200, 200, 200, 10, "원본 이미지와 비율이 다른 가로 세로 길이가 주어지면, 해당 길이 만큼 잘라내어 변환된다."},
		}

		for _, item := range testImgs {
			//fmt.Sprintf("%s 이미지가 %d x %d 사이즈로 리사이징 요청이 들어오면, %d x %d 로 변환된다.", item.imgName, item.reqWidth, item.reqHeight, item.expectedHeight, item.expectedWidth)
			g.It(item.description, func() {
				runLoadImageWithResize(g, item.imgName, item.reqWidth, item.reqHeight, item.expectedWidth, item.expectedHeight, item.distance)
			})
		}
	})

	g.Describe("올바르지 않은 이미지가 주어지면, 변환에 실패한다.", func() {

		g.BeforeEach(func() {
			err := config.SetupConfig()
			if err != nil {
				g.Fatalf("config.SetupConfig() error : %v", err)
				g.FailNow()
			}
			config.RootConfig.SavePath = path.Join("test", "data")
		})

		testImgs := []struct {
			imgName    string
			reqWidth   int
			reqHeight  int
			failReason string
		}{
			{"427×640.jpg", 500, 0, "jpg는 지원하지 않는 형식이므로"},
			{"378×504.png", 0, 500, "png는 지원하지 않는 형식이므로"},
			{"not_exist.jpeg", 0, 0, "파일을 찾을 수 없어"},
		}

		for _, item := range testImgs {
			g.It(fmt.Sprintf("%s 이미지가 주어지면, %s 변환에 실패한다.", item.imgName, item.failReason), func() {
				runFailLoadImageWithResize(g, item.imgName, item.reqWidth, item.reqHeight)
			})
		}
	})
}

func runLoadImageWithResize(g *goblin.G, targetFileName string, reqWidth, reqHeight, expectedWidth, expectedHeight, distance int) {

	result, err := LoadImageWithResize(targetFileName, reqWidth, reqHeight)
	g.Assert(err).IsNil(err)
	g.Assert(result.Bounds().Dx()).Equal(expectedWidth)
	g.Assert(result.Bounds().Dy()).Equal(expectedHeight)

	originImgHash, err := imageTestUtil.GetTestImageHash(targetFileName)
	if err != nil {
		g.Fatalf("config.SetupConfig() error : %v", err)
		g.FailNow()
	}

	newImgHash, err := goimagehash.AverageHash(result)
	g.Assert(err).IsNil(err)

	isSame, err := imageTestUtil.IsSame(originImgHash, newImgHash, distance)
	g.Assert(err).IsNil(err)
	g.Assert(isSame).IsTrue()
}

func runFailLoadImageWithResize(g *goblin.G, targetFileName string, reqWidth, reqHeight int) {

	_, err := LoadImageWithResize(targetFileName, reqWidth, reqHeight)
	g.Assert(err).IsNotNil(err)

}
