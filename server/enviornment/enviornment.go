package environment

import (
	"github.com/joho/godotenv"
	"fmt"
	"os"
	"strconv"
)

func InitEnv() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
        fmt.Printf("Error loading environment variables: %s\n", envErr)
        os.Exit(1)
    }
}

func LoadEnv() (bool, string, string, string, string, string, string) {
	InitEnv()
	envMap, mapErr := godotenv.Read(".env")
	if mapErr != nil {
		fmt.Printf("Error loading environment into map[string]string\n")
		os.Exit(1)
	}

	// Converting string to boolean
	development, convErr := strconv.ParseBool(envMap["DEVELOPMENT"])
	if convErr != nil {
		fmt.Println("Error parsing DEVELOPMENT:", convErr)
		os.Exit(1)
	}
	mongo_url := envMap["MONGO_URI"]
	validation_user := envMap["VALIDATION_USER"]
	validation_pwd := envMap["VALIDATION_PWD"]
	db := envMap["DB"]
	staging_col := envMap["STAGING_COL"]
	validated_col := envMap["VALIDATED_COL"]

	return development, mongo_url, validation_user, validation_pwd, db, staging_col, validated_col
}