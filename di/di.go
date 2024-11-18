package di

import (
	"github.com/hrabit64/springnote-breezenote/auth"
	"github.com/hrabit64/springnote-breezenote/controller"
	"github.com/hrabit64/springnote-breezenote/middleware"
	"github.com/hrabit64/springnote-breezenote/router"
	"github.com/hrabit64/springnote-breezenote/service"
)

// DependencyContext 는 의존성 주입을 위한 구조체입니다.
var DependencyContext = NewApplicationContext()

// InitApplicationContext 는 의존성 주입을 초기화합니다.
func InitApplicationContext() {
	//auth
	DependencyContext.AuthClient = auth.NewFirebaseAuthClient()

	//middleware
	DependencyContext.AuthMiddleware = middleware.NewAuthMiddleware(DependencyContext.AuthClient)

	//services
	DependencyContext.ItemService = service.NewItemService()

	//controllers
	DependencyContext.ImageController = controller.NewImageController(DependencyContext.ItemService)

	//routers
	DependencyContext.ImageRouter = router.NewImageRouter(DependencyContext.ImageController, DependencyContext.AuthMiddleware)
}
