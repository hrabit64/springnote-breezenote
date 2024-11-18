package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hrabit64/springnote-breezenote/controller"
	"github.com/hrabit64/springnote-breezenote/middleware"
)

type ImageRouter interface {
	PostImage(r *gin.Engine) *gin.Engine
	GetImage(r *gin.Engine) *gin.Engine
	InitPath(r *gin.Engine) *gin.Engine
}

type imageRouter struct {
	imgController  controller.ImageController
	authMiddleware middleware.AuthMiddleware
}

func NewImageRouter(imgController controller.ImageController, authMiddleware middleware.AuthMiddleware) ImageRouter {
	return &imageRouter{imgController: imgController, authMiddleware: authMiddleware}
}

func (i *imageRouter) PostImage(r *gin.Engine) *gin.Engine {
	authorized := r.Group("/api/v1/image")
	authorized.Use(i.authMiddleware.UseAuthMiddleware())
	authorized.POST("", i.imgController.UploadImage)

	return r
}

func (i *imageRouter) GetImage(r *gin.Engine) *gin.Engine {
	r.GET("/:id", i.imgController.GetImage)
	return r
}

func (i *imageRouter) InitPath(r *gin.Engine) *gin.Engine {
	r = i.PostImage(r)
	r = i.GetImage(r)
	return r
}
