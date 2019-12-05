package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"site-monitor/global"
	"site-monitor/services"
)

func CheckUrls(c *gin.Context) {
	if global.IsRunning == true {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"err": "the check url program is running, please wait",
		})
		return
	}

	go func() {
		global.IsRunning = true
		defer func() {
			global.IsRunning = false
		}()

		err := services.CheckUrl()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err": err.Error(),
			})
			return
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"msg": "checking url",
	})
}
