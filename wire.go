//go:build wireinject
// +build wireinject

package wireinject

import (
	"github.com/google/wire"
	"github.com/hrabit64/springnote-breezenote/auth"
	"github.com/hrabit64/springnote-breezenote/controller"
	"github.com/hrabit64/springnote-breezenote/middleware"
	"github.com/hrabit64/springnote-breezenote/service"
)

func InitImageController(itemService service.ItemService) controller.ImageController {
	wire.Build(controller.NewImageController)
	return nil
}

func InitAuthMiddleware(client auth.AuthenticateClient) middleware.AuthMiddleware {
	wire.Build(middleware.NewAuthMiddleware)
	return nil
}
