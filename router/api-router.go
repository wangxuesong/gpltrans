package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"gpltrans/controller"
	_ "gpltrans/docs"
)

func InitApiRouter(r *gin.Engine) {
	// swagger router
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// api router
	apiRouter := r.Group("/api")

	v1Router := apiRouter.Group("/v1")
	v1Router.POST("/trans/create", controller.CreateTrans)
}
