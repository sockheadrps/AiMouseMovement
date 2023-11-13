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
	// set variables for production vs development
	var htmlFilesPath string
	var staticFilesPath string
	if development {
		htmlFilesPath = "./index.html"
		staticFilesPath = "./assets"
	} else {
		dir := getCurrentDirectory()
		htmlFilesPath = (dir + "/index.html") 
		staticFilesPath = (dir + "/assets")
	}


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

	router.LoadHTMLFiles(htmlFilesPath)
	router.Static("/assets", staticFilesPath)
	
	router.GET("/", http.HTMLHandler)
	router.POST("/add_data", func(ctx *gin.Context) {
		httpService.AddDataHandler(ctx, &mongoClient)
	})
	
	router.Run("0.0.0.0:9000")
}
