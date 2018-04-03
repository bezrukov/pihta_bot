package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type finCtrl struct {
}

func newFinCtrl() *finCtrl {
	return &finCtrl{}
}

func (ctrl *finCtrl) balance(c *gin.Context) {
	c.JSON(http.StatusOK, "2035.00")
}

func (ctrl *finCtrl) deals(c *gin.Context) {
	deals := []struct {
		Assets  string
		Percent int
		Amount  float32
		Date    string
		Success bool
	}{
		{"GBP USD", 80, 1000, "01.04.2018 14:40", true},
		{"USD CHF", 75, 2000, "01.04.2018 14:46", false},
		{"GBP USD", 80, 300, "01.04.2018 14:51", true},
	}

	c.JSON(http.StatusOK, deals)
}
