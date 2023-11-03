package environment

import (
	"github.com/joho/godotenv"
)

func InitEnv() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Printf("Error loading environment variables")
		os.Exit(1)
	}

	EnvMap, mapErr := godotenv.Read(".env")
	if mapErr != nil {
		fmt.Printf("Error loading environment into map[string]string\n")
		os.Exit(1)
	}
	return EnvMap
}