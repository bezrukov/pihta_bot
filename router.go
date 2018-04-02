package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initRoutes() {
	router := gin.Default()
	gin.SetMode("develop")

	botCtrl := newBotCtrl()

	router.HEAD("/", ping) // this route for consul health check
	router.GET("/", ping)  // this route for consul health check

	router.GET("/test", botCtrl.testRoute)
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "")
}
