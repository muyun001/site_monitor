package util

import (
	"math/rand"
	"time"
)

func RandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36 Edge/15.15063",
		"Mozilla/5.0 (Windows NT 10.0; …) Gecko/20100101 Firefox/59.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	randIndex := r.Intn(len(userAgents))
	ret := userAgents[randIndex : randIndex+1]
	return ret[0]
}
