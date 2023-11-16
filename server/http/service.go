package http

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sockheadrps/AiMouseMovement/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
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

func (c Service) GetRandomDocumentHandler(context *gin.Context, mongoClient *mongo.Client, verificationUUID string) {
	// Specify the collection
	collection := mongoClient.Collection("mousedb", "mouse")

	// Create a filter (empty in this example, you can specify conditions)
	filter := bson.D{}

	// Execute the query to get the count of documents
	count, err := collection.CountDocuments(context, filter)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document count"})
		return
	}

	// Generate a random index within the range of available documents
	randomIndex := rand.Int63n(count)

	// Execute the query to find the document at the random index
	cursor, err := collection.Find(context, filter, options.Find().SetSkip(randomIndex).SetLimit(1))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve random document"})
		return
	}
	defer cursor.Close(context)

	// Iterate through the results (should be only one result)
	var result RandomDocument
	if cursor.Next(context) {
		if err := cursor.Decode(&result); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode random document"})
			return
		}
		result.Uuid = verificationUUID
	}

	// Check for errors from iterating over the cursor
	if err := cursor.Err(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error during cursor iteration"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"randomDocument": result})
}

func (c Service) AuthValidatorHandler(context *gin.Context, validationUser string, validationPwd string, verificationUUID string) {

	// Create an instance of the Validator struct
	var jsonData validator

	// Bind the JSON data from the request body
	if err := context.ShouldBindJSON(&jsonData); err != nil {
		// Handle the error if the JSON binding fails
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate against the provided validation_user and validation_pwd
	if jsonData.User != validationUser || jsonData.Pwd != validationPwd {
		// Return an error response if validation fails
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"uuid":   verificationUUID,
	})
}

func (c Service) UuidAuthHandler(context *gin.Context, verificationUUID string) {

	// Create an instance of the Validator struct
	var jsonUuid clientUuid

	// Bind the JSON data from the request body
	if err := context.ShouldBindJSON(&jsonUuid); err != nil {
		// Handle the error if the JSON binding fails
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(jsonUuid.Uuid, verificationUUID)

	// Validate against the provided validation_user and validation_pwd
	if jsonUuid.Uuid != verificationUUID {
		// Return an error response if validation fails
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid UUID"})
		context.Redirect(http.StatusSeeOther, "/validate")
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (c Service) AddApprovedDataHandler(context *gin.Context, mongoClient *mongo.Client, verificationUUID string) {
	var newDataSet aprovedDataSet

	// Bind the JSON data from the request body
	if err := context.ShouldBindJSON(&newDataSet); err != nil {
		// Handle the error if the JSON binding fails
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	aprovedDataSet := dataSet{
		WindowHeight: newDataSet.WindowHeight,
		WindowWidth:  newDataSet.WindowWidth,
		MouseArray:   newDataSet.MouseArray,
	}
	fmt.Println(aprovedDataSet)

	if newDataSet.Uuid == verificationUUID {
		// Insert data into MongoDB using the existing mongoClient variable
		mongoClient.Insert(context, "mousedb", "validdata", aprovedDataSet)

		// Remove from staging db
		fmt.Println("newDataSet.Id")
		fmt.Println(newDataSet.Id)
		mongoClient.RemoveByID(context, "mousedb", "mouse", newDataSet.Id)

		context.IndentedJSON(http.StatusCreated, newDataSet)
	} else {
		context.Redirect(http.StatusSeeOther, "/validate")
		return
	}

}

func (c Service) RemoveDataHandler(context *gin.Context, mongoClient *mongo.Client, verificationUUID string) {
	var newDataSet removeDataSet

	// Bind the JSON data from the request body
	if err := context.ShouldBindJSON(&newDataSet); err != nil {
		// Handle the error if the JSON binding fails
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newDataSet.Uuid == verificationUUID {
		// Remove from staging db
		mongoClient.RemoveByID(context, "mousedb", "mouse", newDataSet.Id)

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
