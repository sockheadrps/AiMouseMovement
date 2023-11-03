package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/sockheadrps/AiMouseMovement/http"
	"github.com/sockheadrps/AiMouseMovement/mongo"
	"github.com/sockheadrps/AiMouseMovement/enviornment"
	"github.com/gin-gonic/gin"
)

func getCurrentDirectory() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.Fatal(err)
    }
	fmt.Println("Current directory is:", dir)
    return dir
}

func main() {
	ctx := context.Background()
	logger := log.New(os.Stdout, "", log.Lshortfile)

	// init env variables
	development, mongo_url := environment.LoadEnv()
	fmt.Printf("Type: %T, Value: %v\n", development, development)


	// initialize the mongo connection
	mongoClient := mongo.NewClient()

	opts := mongoClient.BuildMongoOptions(mongo_url)
	err := mongoClient.Connect(ctx, opts)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		err = mongoClient.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	httpService := http.NewService()
	router := gin.Default()

	// set env variables for production vs development
	var endpoint string
	if development {
		fmt.Println("dev mode enabled")
		router.LoadHTMLFiles("./index.html")
		router.Static("/assets", "./assets")
		endpoint = "localhost:9090"
	} else {
		fmt.Println("production mode enabled")
		dir := getCurrentDirectory()
		router.LoadHTMLFiles(dir + "/index.html")
		router.Static("/assets", dir + "/assets")
		endpoint = "0.0.0.0:9090"
	}
	
	
	router.GET("/", http.HTMLHandler)
	router.POST("/add_data", func(ctx *gin.Context) {
		httpService.AddDataHandler(ctx, &mongoClient)
	})

	router.Run(endpoint)
}
