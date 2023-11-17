package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sockheadrps/AiMouseMovement/mongoHandler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	dataArray []dataSet
}

func NewService() Service {
	return Service{
		dataArray: []dataSet{},
	}
}

func (c Service) AddDataHandler(context *gin.Context, mongoClient *mongoHandler.Client, db string, staging_col string) {
	var newDataSet dataSet

	if err := context.BindJSON(&newDataSet); err != nil {
		return
	}

	mongoClient.Insert(context, db, staging_col, newDataSet)
	context.IndentedJSON(http.StatusCreated, newDataSet)
}

func (c Service) GetDocumentCountHandler(context *gin.Context, mongoClient *mongoHandler.Client, db string, staging_col string) {
	collection := mongoClient.Collection(db, staging_col)

	filter := bson.D{}

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

	if err := cursor.Err(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error during cursor iteration"})
		return
	}

	count := len(results)

	context.JSON(http.StatusOK, gin.H{"count": count})
}

func (c Service) GetRandomDocumentHandler(context *gin.Context, mongoClient *mongoHandler.Client, verificationUUID string, db string, staging_col string) {
	collection := mongoClient.Collection(db, staging_col)

	filter := bson.D{}

	count, err := collection.CountDocuments(context, filter)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document count"})
		return
	}

	if count == 0 {
		context.JSON(http.StatusOK, gin.H{"error": "No documents", "documents": count})
		return
	}

	result := collection.FindOne(context, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve random document"})
			return
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
			return
		}
	}

	// Decode the result into a RandomDocument struct
	var randomDocument RandomDocument
	if err := result.Decode(&randomDocument); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode random document"})
		return
	}

	// Set the Uuid only if there are results
	randomDocument.Uuid = verificationUUID
	context.JSON(http.StatusOK, gin.H{"randomDocument": randomDocument, "documents": count})
}

func (c Service) AuthValidatorHandler(context *gin.Context, validationUser string, validationPwd string, verificationUUID string) {
	var jsonData validator

	// Bind the JSON data from the request body
	if err := context.ShouldBindJSON(&jsonData); err != nil {
		// Handle the error if the JSON binding fails
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate against the provided validation_user and validation_pwd
	if jsonData.User != validationUser || jsonData.Pwd != validationPwd {
		context.JSON(http.StatusUnauthorized, gin.H{"status": "Invalid credentials"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"uuid":   verificationUUID,
	})
}

func (c Service) UuidAuthHandler(context *gin.Context, verificationUUID string) {
	var jsonUuid clientUuid

	// Bind the JSON data from the request body
	if err := context.ShouldBindJSON(&jsonUuid); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"erdror": err.Error()})
		return
	}

	// Validate client UUID to server UUID
	if jsonUuid.Uuid != verificationUUID {
		// Return an error response if validation fails
		context.Redirect(http.StatusSeeOther, "/validate")
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (c Service) AddApprovedDataHandler(context *gin.Context, mongoClient *mongoHandler.Client, verificationUUID string, db string, validated_col string, staging_col string) {
	var newDataSet aprovedDataSet

	// Bind the JSON data from the request body
	if err := context.ShouldBindJSON(&newDataSet); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	aprovedDataSet := dataSet{
		WindowHeight: newDataSet.WindowHeight,
		WindowWidth:  newDataSet.WindowWidth,
		MouseArray:   newDataSet.MouseArray,
	}

	if newDataSet.Uuid == verificationUUID {
		// Insert data into MongoDB validated collection
		mongoClient.Insert(context, db, validated_col, aprovedDataSet)

		// Remove from staging db
		mongoClient.RemoveByID(context, db, staging_col, newDataSet.Id)

		context.IndentedJSON(http.StatusCreated, newDataSet)
	} else {
		context.Redirect(http.StatusSeeOther, "/validate")
		return
	}

}

func (c Service) RemoveDataHandler(context *gin.Context, mongoClient *mongoHandler.Client, verificationUUID string, db string, staging_col string) {
	var newDataSet removeDataSet

	// Bind the JSON data from the request body
	if err := context.ShouldBindJSON(&newDataSet); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newDataSet.Uuid == verificationUUID {
		mongoClient.RemoveByID(context, db, staging_col, newDataSet.Id)

		context.IndentedJSON(http.StatusCreated, newDataSet)
	} else {
		context.Redirect(http.StatusSeeOther, "/validate")
		return
	}
}

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Data Gathering",
	})
}

func ValidationHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "validation.html", gin.H{
		"title": "Data Validation",
	})
}

func ViewDataHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "dataview.html", gin.H{
		"title": "Data View",
	})
}
