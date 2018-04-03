package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"fmt"
)

func initRoutes(port string) {
	router := gin.Default()
	gin.SetMode("debug")

	finCtrl := newFinCtrl()

	router.HEAD("/", ping) // this route for consul health check
	router.GET("/", ping)  // this route for consul health check

	router.GET("/balance", finCtrl.balance)
	router.GET("/deals", finCtrl.deals)

	router.GET("/duo-auth-on")

	router.Run(fmt.Sprintf(":%s", port))
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "")
}
