package img

import (
	"github.com/franela/goblin"
	"github.com/hrabit64/springnote-breezenote/config"
	dirTestUtil "github.com/hrabit64/springnote-breezenote/test/utils/dir"
	fileTestUtil "github.com/hrabit64/springnote-breezenote/test/utils/file"
	imageTestUtil "github.com/hrabit64/springnote-breezenote/test/utils/image"
	"path"
	"testing"
)

func TestSaveGif(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("올바른 Gif이미지가 주어지면, 변환에 성공한다.", func() {
		g.It("정상적인 gif base64 string이 주어지면, 변환하여 저장한다.", func() {
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

			gifb64, err := imageTestUtil.GetTestImageBase64("256×256.gif")
			if err != nil {
				g.Fatalf("imageTestUtil.GetTestImageBase64() error : %v", err)
				g.FailNow()
			}

			err = SaveGif(gifb64, "test")
			g.Assert(err).IsNil(err)

			res, err := fileTestUtil.CompareFiles(
				path.Join(config.RootConfig.SavePath, "test.gif"),
				path.Join("test", "data", "256×256.gif"),
			)
			g.Assert(err).IsNil(err)
			g.Assert(res).IsTrue()
		})
	})
}

func TestLoadGif(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Gif 파일이 주어지면, 해당 파일을 로드한다.", func() {
		g.It("정상적인 파일이름이 주어지면, 해당 파일을 로드한다.", func() {
			err := config.SetupConfig()
			if err != nil {
				g.Fatalf("config.SetupConfig() error : %v", err)
				g.FailNow()
			}

			config.RootConfig.SavePath = path.Join("test", "data")

			_, err = LoadGif("256×256.gif")
			g.Assert(err).IsNil(err)

			// TODO : gif 이미지의 hash값을 비교하는 테스트 코드 추가
		})

		g.It("존재하지 않는 gif 파일이 주어지면, 에러를 반환한다.", func() {
			err := config.SetupConfig()
			if err != nil {
				g.Fatalf("config.SetupConfig() error : %v", err)
				g.FailNow()
			}

			config.RootConfig.SavePath = path.Join("test", "data")

			_, err = LoadGif("notExist.gif")
			g.Assert(err).IsNotNil(err)
		})
	})
}
