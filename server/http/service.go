package http

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/sockheadrps/AiMouseMovement/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct{
	dataArray []dataSet
}

func NewService() Service {
	return Service{
		dataArray: []dataSet{},
	}
}

func (c Service) AddDataHandler(context *gin.Context, mongoClient *mongo.Client) {
	var newDataSet dataSet

	if err := context.BindJSON(&newDataSet); err != nil {
		return
	}

	// Check for required fields in the dataSet struct
	if newDataSet.WindowHeight == 0 || newDataSet.WindowWidth == 0 || len(newDataSet.MouseArray) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Required field(s) are missing"})
		return
	}

	// Insert data into MongoDB using the existing mongoClient variable
    mongoClient.Insert(context, "mousedb", "mouse", newDataSet)
	context.IndentedJSON(http.StatusCreated, newDataSet)
}

func (c Service) GetDocumentCountHandler(context *gin.Context, mongoClient *mongo.Client) {
	// Specify the collection
	collection := mongoClient.Collection("mousedb", "mouse")

	// Create a filter (empty in this example, you can specify conditions)
	filter := bson.D{}

	// Execute the query
	cursor, err := collection.Find(context, filter)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}
	defer cursor.Close(context)

	// Iterate through the results
	var results []bson.M
	for cursor.Next(context) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode result"})
			return
		}
		results = append(results, result)
	}

	// Check for errors from iterating over the cursor
	if err := cursor.Err(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error during cursor iteration"})
		return
	}


	count := len(results)

	context.JSON(http.StatusOK, gin.H{"count": count})
}

func HTMLHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Data Gathering",
	})
}