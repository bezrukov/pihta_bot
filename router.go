package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"fmt"
)

func initRoutes(port string) {
	router := gin.Default()
	gin.SetMode("debug")

	botCtrl := newBotCtrl()

	router.HEAD("/", ping) // this route for consul health check
	router.GET("/", ping)  // this route for consul health check

	router.GET("/test", botCtrl.testRoute)

	router.Run(fmt.Sprintf(":%s", port))
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "")
}
