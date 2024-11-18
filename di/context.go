package di

import (
	"github.com/hrabit64/springnote-breezenote/auth"
	"github.com/hrabit64/springnote-breezenote/controller"
	"github.com/hrabit64/springnote-breezenote/middleware"
	"github.com/hrabit64/springnote-breezenote/router"
	"github.com/hrabit64/springnote-breezenote/service"
)

// ApplicationContext 는 의존성 주입을 위한 구조체입니다.
type ApplicationContext struct {
	ItemService service.ItemService

	AuthClient auth.AuthenticateClient

	AuthMiddleware middleware.AuthMiddleware

	ImageController controller.ImageController

	ImageRouter router.ImageRouter
}

func NewApplicationContext() *ApplicationContext {
	return &ApplicationContext{}
}
