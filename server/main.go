package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/sockheadrps/AiMouseMovement/http"
	"github.com/sockheadrps/AiMouseMovement/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

	dir := getCurrentDirectory()

	// initialize the mongo connection
	mongoClient := mongo.NewClient()

	opts := mongoClient.BuildMongoOptions()
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
	
	// run the admin command and print the output
	res, err := mongoClient.Run(ctx, "admin", bson.D{{"ping", 1}})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Got result from run: %v", res)

	// start the http server
	httpService := http.NewService()

	router := gin.Default()
	router.LoadHTMLFiles(dir + "index.html")
	router.Static("/assets", dir + "./assets")

	router.GET("/", http.HTMLHandler)
	router.POST("/add_data", func(ctx *gin.Context) {
		httpService.AddDataHandler(ctx, &mongoClient)
	})

	router.Run("localhost:9090")
}
