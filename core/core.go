package core

import (
	"github.com/gin-gonic/gin"
	"github.com/hrabit64/springnote-breezenote/config"
	"github.com/hrabit64/springnote-breezenote/di"
	"github.com/hrabit64/springnote-breezenote/middleware"
	"github.com/hrabit64/springnote-breezenote/router"
)

// SetupRouter gin engine을 설정합니다.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	if config.RootConfig.Profile == "release" {
		r.TrustedPlatform = gin.PlatformCloudflare
	}

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.ResponseLogger())
	r.Use(middleware.CORSMiddleware(config.RootConfig.AllowOrigin))
	//setup auth middleware
	routers := configRouters()

	for _, targetRouter := range routers {
		r = targetRouter.InitPath(r)
	}

	return r
}

// configRouters router에 path를 설정합니다.
func configRouters() []router.Router {
	var routers []router.Router

	routers = append(routers, di.DependencyContext.ImageRouter)
	routers = append(routers, router.NewPingRouter())

	return routers
}
