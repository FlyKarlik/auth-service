package server

import (
	"net/http/pprof"

	_ "github.com/FlyKarlik/auth-service/api/docs"
	"github.com/FlyKarlik/auth-service/internal/gateways/http/handler"
	"github.com/FlyKarlik/auth-service/internal/gateways/http/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initRoutes(handler *handler.Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-type"},
	}))

	router.GET("/debug/pprof/", gin.WrapF(pprof.Index))
	router.GET("/debug/pprof/cmdline", gin.WrapF(pprof.Cmdline))
	router.GET("/debug/pprof/profile", gin.WrapF(pprof.Profile))
	router.GET("/debug/pprof/symbol", gin.WrapF(pprof.Symbol))
	router.GET("/debug/pprof/trace", gin.WrapF(pprof.Trace))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("v1")
	v1.Use(middleware.JSONMiddleware, middleware.BindClientIP)
	{
		auth := v1.Group("auth")
		auth.Use()
		{
			auth.POST("/:id", handler.Authenticate)
			auth.POST("/refresh/", handler.Refresh)
		}
	}

	return router
}
