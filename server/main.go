package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mousePoint struct {
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Time int64   `json:"time"`
}

type dataSet struct {
	WindowHeight int32      `json:"window-height"`
	WindowWidth  int32      `json:"window-width"`
	MouseArray   []mousePoint `json:"mouse-array"`
}

var dataArray = []dataSet{}

func addData(context *gin.Context) {
	var newDataSet dataSet

	if err := context.BindJSON(&newDataSet); err != nil {
		return
	}

	dataArray = append(dataArray, newDataSet)
	context.IndentedJSON(http.StatusCreated, dataArray)
}

func main() {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// pls dont bomb my free tier mongo db
	opts := options.Client().ApplyURI("mongodb+srv://sockheadrps:s3eUDQdQqO82UYH2@cluster0.g5abd.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	coll := client.Database("mousedb").Collection("mouse")

	router := gin.Default()
	router.POST("/add_data", addData)
	router.Run("localhost:9090")
}
