package init

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	checkEnv()
}

func checkEnv() {
	_ = godotenv.Load()
	needChecks := []string{
		"RECEIVE_URLS_AND_SEND_RESULT_API",
		"HEADLESS_URL_PREFIX",
	}

	for _, envKey := range needChecks {
		if os.Getenv(envKey) == "" {
			log.Fatalf("env %s missed", envKey)
		}
	}
}
