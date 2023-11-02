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

	// Insert data into MongoDB using the existing mongoClient variable
    mongoClient.Insert(context, "mousedb", "mouse", newDataSet)
	context.IndentedJSON(http.StatusCreated, newDataSet)
}