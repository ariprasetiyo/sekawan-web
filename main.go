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

	router.Static("/authentication-login", "./web/templates")
	router.Static("/href", "./web/templates")

	// router.GET("/authentication-login", controller.Contact)
	router.GET("/authentication-register", controller.Contact)
	router.GET("/charts", controller.Error404)
	// router.GET("/error-403", controller.Error404)
	router.GET("/error-405", controller.Error404)
	router.GET("/error-500", controller.Error404)
	router.GET("/form-basic", controller.Error404)
	router.GET("/form-wizard", controller.Error404)
	router.GET("/grid", controller.Error404)
	router.GET("/icon-fontawesome", controller.Error404)
	router.GET("/icon-meterial", controller.Error404)
	router.GET("/index2", controller.Error404)
	router.GET("/pages-buttons", controller.Error404)
	router.GET("/pages-calendar", controller.Error404)
	router.GET("/pages-chat", controller.Error404)
	router.GET("/pages-elements", controller.Error404)
	router.GET("/pages-gallery", controller.Error404)
	router.GET("/pages-invoice", controller.Error404)
	router.GET("/tables", controller.Error404)
	router.GET("/widgets", controller.Error404)

	router.LoadHTMLGlob("web/templates/*")

	router.GET("/", controller.Index)
	router.GET("/error-404", controller.Error404)

	router.Run(":" + server.HTTP_SERVER_PORT)
}
