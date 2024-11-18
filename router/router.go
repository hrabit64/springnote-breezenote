package router

import "github.com/gin-gonic/gin"

type Router interface {
	InitPath(r *gin.Engine) *gin.Engine
}
