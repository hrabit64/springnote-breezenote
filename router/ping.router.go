package router

import "github.com/gin-gonic/gin"

type PingRouter interface {
	GetPing(r *gin.Engine) *gin.Engine
	InitPath(r *gin.Engine) *gin.Engine
}

type pingRouter struct{}

func NewPingRouter() PingRouter {
	return &pingRouter{}
}

func (p *pingRouter) GetPing(r *gin.Engine) *gin.Engine {
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func (p *pingRouter) InitPath(r *gin.Engine) *gin.Engine {
	r = p.GetPing(r)
	return r
}
