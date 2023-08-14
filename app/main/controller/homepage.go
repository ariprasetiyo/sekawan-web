package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		// "views/indexx.html",
		// "matrix-Ip.html",
		"index.html",
		gin.H{
			"title":      "Geeksbeginner",
			"textRender": "Hallo ari",
		},
	)
}

func Error404(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"error-404.html",
		gin.H{},
	)
}
