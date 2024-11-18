package setupTestUtils

import (
	"github.com/hrabit64/springnote-breezenote/config"
	"github.com/hrabit64/springnote-breezenote/database"
	"github.com/hrabit64/springnote-breezenote/di"
	"github.com/hrabit64/springnote-breezenote/middleware"
	"github.com/hrabit64/springnote-breezenote/router"
	"github.com/hrabit64/springnote-breezenote/test/mock"
)

// SetupConfig 설정을 초기화합니다.
func SetupConfig() error {
	return config.SetupConfig()
}

// SetupDb DB를 초기화합니다.
func SetupDb() error {
	return database.RunSetup()
}

// SetupTestGin 테스트를 위한 gin을 설정합니다. db,config,di를 초기화합니다. 또한 테스트용 mock auth client를 설정합니다.
func SetupTestGin() error {
	err := SetupConfig()
	if err != nil {
		return err
	}
	err = SetupDb()
	if err != nil {
		return err
	}

	di.InitApplicationContext()
	di.DependencyContext.AuthMiddleware = middleware.NewAuthMiddleware(mock.NewMockAuthClient())
	di.DependencyContext.AuthClient = mock.NewMockAuthClient()
	di.DependencyContext.ImageRouter = router.NewImageRouter(di.DependencyContext.ImageController, di.DependencyContext.AuthMiddleware)
	return nil
}
