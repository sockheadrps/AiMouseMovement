package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service struct{
	dataArray []dataSet
}

func NewService() Service {
	return Service{
		dataArray: []dataSet{},
	}
}

func (c Service) AddDataHandler(context *gin.Context) {
	var newDataSet dataSet

	if err := context.BindJSON(&newDataSet); err != nil {
		return
	}

	c.dataArray = append(c.dataArray, newDataSet)
	context.IndentedJSON(http.StatusCreated, c.dataArray)
}