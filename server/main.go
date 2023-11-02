package main

import (
	"context"
	"log"
	"os"

	"github.com/sockheadrps/AiMouseMovement/http"
	"github.com/sockheadrps/AiMouseMovement/mongo"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	ctx := context.Background()
	logger := log.New(os.Stdout, "", log.Lshortfile)

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

	// get and print the mouse collection
	mouseColl := mongoClient.Collection("mousedb", "mouse")
	log.Println(mouseColl)

	// start the http server
	httpService := http.NewService()

	router := gin.Default()
	router.POST("/add_data", httpService.AddDataHandler)
	router.Run("localhost:9090")
}
