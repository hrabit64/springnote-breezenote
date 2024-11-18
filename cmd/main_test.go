package main

import (
	"github.com/corona10/goimagehash"
	"github.com/franela/goblin"
	"github.com/hrabit64/springnote-breezenote/config"
	"github.com/hrabit64/springnote-breezenote/di"
	"github.com/hrabit64/springnote-breezenote/dto"
	"github.com/hrabit64/springnote-breezenote/test/mock"
	testUtils "github.com/hrabit64/springnote-breezenote/test/utils"
	dirTestUtil "github.com/hrabit64/springnote-breezenote/test/utils/dir"
	imageTestUtil "github.com/hrabit64/springnote-breezenote/test/utils/image"
	"github.com/hrabit64/springnote-breezenote/test/utils/setup"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/jpeg"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Ping 테스트", func() {
		g.It("서버가 정상적으로 동작하면, pong을 반환한다.", func() {
			err := setupTestUtils.SetupConfig()
			g.Assert(err).IsNil(err)

			di.InitApplicationContext()

			w := testUtils.RunTestApiReq("GET", "/ping", nil, nil)
			g.Assert(w.Code).Equal(200)
			g.Assert(w.Body.String()).Equal("pong")
		})
	})
}

func TestImageUpload(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("올바른 이미지 파일이 주어지면, 업로드에 성공한다.", func() {
		items := []struct {
			imgName     string
			description string
		}{
			{"427×640.jpg", "jpeg 이미지가 주어지면, webp, jpeg로 변환되어 저장된다."},
			{"378×504.webp", "webp 이미지가 주어지면, webp, jpeg로 변환되어 저장된다."},
			{"378×504.png", "png 이미지가 주어지면, webp, jpeg로 변환되어 저장된다."},
			{"427×640.jpeg", "jpeg 이미지가 주어지면, webp, jpeg로 변환되어 저장된다."},
		}

		for _, item := range items {
			g.It(item.description, func() {
				runTestImageUpload(g, item.imgName)
			})
		}
	})

	g.Describe("지원하지 않는 이미지가 주어지면, 400 에러를 반환한다.", func() {
		notimage := "7KCE7Jet7ZWY6rOg7Iu264ukLg=="

		g.It("pdf 파일이 주어지면, 400 에러를 반환한다.", func() {
			testBody := &dto.ImageCreateRequest{
				Image: notimage,
				Name:  "notimage.txt",
			}

			w := testUtils.RunTestApiReq("POST",
				"/api/v1/image",
				testUtils.ConvertToJsonReader(testBody),
				http.Header{
					"Authorization": {"Bearer " + mock.AdminToken},
				},
			)
			g.Assert(w.Code).Equal(400)
		})

		g.It("이름은 이미지파일이나, 실제 파일이 지원하지 않는 형식이면, 400 에러를 반환한다.", func() {
			testBody := &dto.ImageCreateRequest{
				Image: notimage,
				Name:  "haha.jpg",
			}

			w := testUtils.RunTestApiReq("POST",
				"/api/v1/image",
				testUtils.ConvertToJsonReader(testBody),
				http.Header{
					"Authorization": {"Bearer " + mock.AdminToken},
				},
			)
			g.Assert(w.Code).Equal(400)
		})

	})

	g.Describe("권한이 없는 사용자가 이미지를 업로드하려고 하면, 401 에러를 반환한다.", func() {
		g.It("유효한 사용자이나, 권한이 없는 사용자가 이미지를 업로드하려고 하면, 401 에러를 반환한다.", func() {
			testBody := &dto.ImageCreateRequest{
				Image: "test",
				Name:  "test.jpg",
			}

			w := testUtils.RunTestApiReq("POST",
				"/api/v1/image",
				testUtils.ConvertToJsonReader(testBody),
				http.Header{
					"Authorization": {"Bearer " + mock.UserToken},
				},
			)
			g.Assert(w.Code).Equal(401)
		})

		g.It("유효하지 않은 사용자가 이미지를 업로드하려고 하면, 401 에러를 반환한다.", func() {
			testBody := &dto.ImageCreateRequest{
				Image: "test",
				Name:  "test.jpg",
			}

			w := testUtils.RunTestApiReq("POST",
				"/api/v1/image",
				testUtils.ConvertToJsonReader(testBody),
				http.Header{
					"Authorization": {"Bearer " + "lol"},
				},
			)
			g.Assert(w.Code).Equal(401)
		})

		g.It("토큰이 주어지지 않으면, 401 에러를 반환한다.", func() {
			testBody := &dto.ImageCreateRequest{
				Image: "test",
				Name:  "test.jpg",
			}

			w := testUtils.RunTestApiReq("POST",
				"/api/v1/image",
				testUtils.ConvertToJsonReader(testBody),
				nil,
			)
			g.Assert(w.Code).Equal(401)
		})
	})
}

func runTestImageUpload(g *goblin.G, imgName string) {
	dirName, err := dirTestUtil.CreateTmpDir()
	if err != nil {
		g.Fatalf("dirTestUtil.CreateTmpDir() error : %v", err)
		g.FailNow()
	}
	defer func() {
		err := dirTestUtil.RemoveTmpDir(dirName)
		if err != nil {
			g.Fatalf("dirTestUtil.RemoveTmpDir() error : %v", err)
			g.FailNow()
		}
	}()

	err = setupTestUtils.SetupTestGin()
	g.Assert(err).IsNil(err)

	config.RootConfig.SavePath = path.Join("tmp", dirName)

	imgBase64, err := imageTestUtil.GetTestImageBase64(imgName)
	if err != nil {
		g.Fatalf("imageTestUtil.GetTestImageBase64() error : %v", err)
		g.FailNow()
	}

	testBody := &dto.ImageCreateRequest{
		Image: imgBase64,
		Name:  imgName,
	}

	w := testUtils.RunTestApiReq("POST",
		"/api/v1/image",
		testUtils.ConvertToJsonReader(testBody),
		http.Header{
			"Authorization": {"Bearer " + mock.AdminToken},
		},
	)
	loggingBody := w.Body.String()
	log.Printf("loggingBody : %v", loggingBody)
	g.Assert(w.Code).Equal(200)

	bodyMap := testUtils.ConvertResToMap(w)
	g.Assert(bodyMap["file_id"]).IsNotNil()

	newImgName := bodyMap["file_id"].(string)

	originHash, err := imageTestUtil.GetTestImageHash(imgName)
	if err != nil {
		g.Fatalf("imageTestUtil.GetTestImageHash() error : %v", err)
		g.FailNow()
	}

	webpHash, err := imageTestUtil.GetTestImageHashFromOtherDir(newImgName+".webp", config.RootConfig.SavePath)
	g.Assert(err).IsNil(err)

	jpegHash, err := imageTestUtil.GetTestImageHashFromOtherDir(newImgName+".jpeg", config.RootConfig.SavePath)
	g.Assert(err).IsNil(err)

	isSame, err := imageTestUtil.IsSame(originHash, webpHash, 5)
	g.Assert(err).IsNil(err)
	g.Assert(isSame).IsTrue()

	isSame, err = imageTestUtil.IsSame(originHash, jpegHash, 5)
	g.Assert(err).IsNil(err)
	g.Assert(isSame).IsTrue()
}

// 이미지 조회 테스트
func TestGetImage(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("이미지 이름이 주어지면, 해당 이미지를 반환한다.", func() {

		testImgs := []struct {
			imgName     string
			description string
		}{
			{"427×640.jpeg", "jpeg 이미지가 주어지면, 해당 이미지를 반환한다."},
			{"378×504.webp", "webp 이미지가 주어지면, 해당 이미지를 반환한다."},
			{"427×640.jpg", "jpg 이미지가 주어지면, 확장자를 jpeg로 변경하여 해당 이미지를 반환한다."},
		}
		for _, item := range testImgs {
			g.It(item.description, func() {
				runTestGetImage(g, item.imgName)
			})
		}
	})

	g.Describe("이미지 이름과, 리사이징 조건이 주어지면, 리사이징 된 이미지를 반환한다.", func() {
		testImgs := []struct {
			imgName     string
			reqW        int
			reqH        int
			resW        int
			resH        int
			distance    int
			description string
		}{
			{"427×640.jpeg", 200, 300, 200, 300, 10, "jpeg 이미지와 가로 세로 길이가 모두 주어지면, 해당 길이로 변환된다."},
			{"427×640.jpeg", 200, 0, 200, 300, 10, "jpeg 이미지와 가로 길이만 주어지면, 가로 길이가 주어진 길이로 변환되고, 세로 길이는 비율에 맞춰 변환된다."},
			{"427×640.jpeg", 0, 300, 200, 300, 10, "jpeg 이미지와 세로 길이만 주어지면, 세로 길이가 주어진 길이로 변환되고, 가로 길이는 비율에 맞춰 변환된다."},
			{"378×504.webp", 189, 252, 189, 252, 10, "webp 이미지와 세로 길이만 주어지면, 세로 길이가 주어진 길이로 변환되고, 가로 길이는 비율에 맞춰 변환된다."},
			{"378×504.webp", 189, 0, 189, 252, 10, "webp 이미지와 가로 길이만 주어지면, 가로 길이가 주어진 길이로 변환되고, 세로 길이는 비율에 맞춰 변환된다."},
			{"378×504.webp", 0, 252, 189, 252, 10, "webp 이미지와 가로 세로 길이가 모두 주어지면, 해당 길이로 변환된다."},
		}

		for _, item := range testImgs {
			g.It(item.description, func() {
				runTestGetImageResize(g, item.imgName, item.reqW, item.reqH, item.resW, item.resH, item.distance)
			})

		}
	})

	g.Describe("올바르지 않은 요청이 주어지면, 400 에러를 반환한다.", func() {
		testImgs := []struct {
			imgName     string
			reqW        int
			reqH        int
			status      int
			description string
		}{
			{"427×640.jpeg", -1, 300, 400, "가로 길이가 음수이면, 400 에러를 반환한다."},
			{"427×640.jpeg", 200, -1, 400, "세로 길이가 음수이면, 400 에러를 반환한다."},
			{"427×640.jpeg", config.RootConfig.MaxImageLen + 1, config.RootConfig.MaxImageLen + 1, 400, "최대 길이보다 큰 길이가 주어지면, 400 에러를 반환한다."},
			{"not_exist.jpeg", 200, 300, 404, "존재하지 않는 이미지 이름이 주어지면, 404 에러를 반환한다."},
			{"not_support_format.bmp", 200, 300, 400, "지원하지 않는 이미지 형식이 주어지면, 400 에러를 반환한다."},
		}

		for _, item := range testImgs {
			g.It(item.description, func() {
				w := testUtils.RunTestApiReq("GET", "/"+item.imgName+"?width="+strconv.Itoa(item.reqW)+"&height="+strconv.Itoa(item.reqH), nil, nil)
				g.Assert(w.Code).Equal(item.status)
			})
		}
	})
}

func runTestGetImage(g *goblin.G, imgName string) {
	err := setupTestUtils.SetupConfig()
	if err != nil {
		g.Fatalf("config.SetupConfig() error : %v", err)
		g.FailNow()
	}

	di.InitApplicationContext()

	config.RootConfig.SavePath = path.Join("test", "data")

	w := testUtils.RunTestApiReq("GET", "/"+imgName, nil, nil)

	g.Assert(w.Code).Equal(200)
	ext := strings.Split(imgName, ".")[1]
	if ext == "jpg" {
		ext = "jpeg"
	}
	g.Assert(w.Header().Get("Content-Type")).Equal("image/" + ext)

	originHash, err := imageTestUtil.GetTestImageHash(imgName)
	if err != nil {
		g.Fatalf("imageTestUtil.GetTestImageHash() error : %v", err)
		g.FailNow()
	}

	img, _, err := image.Decode(w.Body)
	g.Assert(err).IsNil(err)

	resHash, err := goimagehash.AverageHash(img)
	g.Assert(err).IsNil(err)

	isSame, err := imageTestUtil.IsSame(originHash, resHash, 5)
	g.Assert(err).IsNil(err)
	g.Assert(isSame).IsTrue()
}

func runTestGetImageResize(g *goblin.G, imgName string, reqW int, reqH int, resW int, resH int, distance int) {
	err := setupTestUtils.SetupConfig()
	if err != nil {
		g.Fatalf("config.SetupConfig() error : %v", err)
		g.FailNow()
	}

	di.InitApplicationContext()

	config.RootConfig.SavePath = path.Join("test", "data")

	url := "/" + imgName

	if reqW > 0 {
		url += "?width=" + strconv.Itoa(reqW)
	}
	if reqH > 0 {
		if strings.ContainsAny(url, "?") {
			url += "&"
		} else {
			url += "?"
		}
		url += "height=" + strconv.Itoa(reqH)
	}

	w := testUtils.RunTestApiReq("GET", url, nil, nil)
	g.Assert(w.Code).Equal(200)
	g.Assert(w.Header().Get("Content-Type")).Equal("image/" + strings.Split(imgName, ".")[1])

	originHash, err := imageTestUtil.GetTestImageHash(imgName)
	if err != nil {
		g.Fatalf("imageTestUtil.GetTestImageHash() error : %v", err)
		g.FailNow()
	}

	img, _, err := image.Decode(w.Body)
	g.Assert(err).IsNil(err)

	resHash, err := goimagehash.AverageHash(img)
	g.Assert(err).IsNil(err)

	isSame, err := imageTestUtil.IsSame(originHash, resHash, distance)
	g.Assert(err).IsNil(err)
	g.Assert(isSame).IsTrue()

	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	g.Assert(imgWidth).Equal(resW)
	g.Assert(imgHeight).Equal(resH)
}
