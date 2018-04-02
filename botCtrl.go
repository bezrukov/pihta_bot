package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type botCtrl struct {
}

func newBotCtrl() *botCtrl {
	return &botCtrl{}
}

func (ctrl *botCtrl) testRoute(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}
