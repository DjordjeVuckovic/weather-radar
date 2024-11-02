package env

import (
	"fmt"
)
import (
	"github.com/joho/godotenv"
)

func Load() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Skipping .env file ...")
	}
}
