package main

import _ "site-monitor/init"

import (
	"github.com/gin-gonic/gin"
	"log"
	"site-monitor/actions"
)

func main() {
	router := gin.Default()
	router.POST("check-urls", actions.CheckUrls)

	err := router.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
