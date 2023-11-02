package http

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/sockheadrps/AiMouseMovement/mongo"
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

func HTMLHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Data Gathering",
	})
}