package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initRoutes() {
	router := gin.Default()
	gin.SetMode("debug")

	botCtrl := newBotCtrl()

	router.HEAD("/", ping) // this route for consul health check
	router.GET("/", ping)  // this route for consul health check

	router.GET("/test", botCtrl.testRoute)

	router.Run(":2000")
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "")
}
