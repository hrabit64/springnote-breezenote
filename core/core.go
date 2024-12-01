package core

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hrabit64/springnote-breezenote/config"
	"github.com/hrabit64/springnote-breezenote/di"
	"github.com/hrabit64/springnote-breezenote/middleware"
	"github.com/hrabit64/springnote-breezenote/router"
	"time"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.RootConfig.AllowOrigin},
		AllowMethods:     []string{"PUT", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
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
