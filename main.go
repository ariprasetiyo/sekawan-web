package main

import (
	"context"
	controller "sekawan-web/app/main/controller"
	repository "sekawan-web/app/main/repository"
	"sekawan-web/app/main/server"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-contrib/cors"
)

var db = make(map[string]string)

func main() {

	server.InitConfig()
	server.InitLogrusFormat()

	// running open telemetry
	cleanup := server.InitTracer()
	defer cleanup(context.Background())

	gin.SetMode(server.GIN_MODE)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// router.Use(server.RequestResponseLogger())
	router.Use(otelgin.Middleware(server.APP_NAME))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	router.Use(cors.New(corsConfig))

	initPosgresSQL := server.InitPostgreSQL()
	_ = repository.NewDatabase(initPosgresSQL.DB, initPosgresSQL)

	router.Static("/assets", "./web/assets")
	router.Static("/dist", "./web/dist")
	router.Static("/src", "./web/src")
	router.Static("/href", "./web/templates")

	router.LoadHTMLGlob("web/templates/*")

	router.GET("/", controller.Index)
	router.GET("/home", controller.Index)
	router.GET("/error-404", controller.Error404)

	router.Run(":" + server.HTTP_SERVER_PORT)
}
