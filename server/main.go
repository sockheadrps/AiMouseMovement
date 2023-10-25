package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	router := gin.Default()
	router.POST("/add_data", addData)
	router.Run("localhost:9090")
}
