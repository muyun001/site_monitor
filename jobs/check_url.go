package jobs

import (
	"log"
	"site-monitor/services"
)

func CheckUrl() {
	err := services.CheckUrl()
	if err != nil {
		log.Fatal(err.Error())
	}
}
