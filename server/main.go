package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/sockheadrps/AiMouseMovement/http"
	"github.com/sockheadrps/AiMouseMovement/mongo"
	"github.com/sockheadrps/AiMouseMovement/enviornment"
	"github.com/gin-gonic/gin"
	"time"
	"github.com/google/uuid"
)

var verificationUUID string

func generateUUID() string {
	return uuid.New().String()
}

func updateUUIDPeriodically() {
	for {
		verificationUUID = generateUUID()
		fmt.Println("Verification UUID updated:", verificationUUID)

		// Sleep for one hour before updating again
		time.Sleep(time.Hour)
	}
}


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
	go updateUUIDPeriodically()

	// init env variables
	development, mongo_url, validation_user, validation_pwd := environment.LoadEnv()
	// set variables for production vs development
	var htmlFilesPath string
	var staticFilesPath string

	if development {
		htmlFilesPath = "./templates/"
		staticFilesPath = "./assets/"
	} else {
		dir := getCurrentDirectory()
		htmlFilesPath = dir + "/templates/"
		staticFilesPath = dir + "/assets/"
		gin.SetMode(gin.ReleaseMode)
	}


	// initialize the mongo connection
	mongoClient := mongo.NewClient()

	opts := mongoClient.BuildMongoOptions(mongo_url)
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
	
	
	httpService := http.NewService()
	router := gin.Default()

	router.LoadHTMLGlob(htmlFilesPath + "*.html")
	router.Static("/assets", staticFilesPath)
	
	router.GET("/", http.IndexHandler)
	router.GET("/validate", http.ValidationHandler)
	router.GET("/view-data", http.ViewDataHandler)

	router.GET("/get-data-point", func(ctx *gin.Context) {
		httpService.GetRandomDocumentHandler(ctx, &mongoClient,  verificationUUID)
	})


	router.GET("/document_count", func(ctx *gin.Context) {
		httpService.GetDocumentCountHandler(ctx, &mongoClient)
	})
	router.POST("/add_data", func(ctx *gin.Context) {
		httpService.AddDataHandler(ctx, &mongoClient)
	})
	router.POST("/approve_data", func(ctx *gin.Context) {
		httpService.AddApprovedDataHandler(ctx, &mongoClient, verificationUUID)
	})
	router.POST("/remove_data", func(ctx *gin.Context) {
		httpService.RemoveDataHandler(ctx, &mongoClient, verificationUUID)
	})
	router.POST("/auth/validate", func(ctx *gin.Context) {
		httpService.AuthValidatorHandler(ctx, validation_user, validation_pwd, verificationUUID)
	})
	router.POST("/auth/uuid", func(ctx *gin.Context) {
		httpService.UuidAuthHandler(ctx, verificationUUID)
	})
	
	router.Run("0.0.0.0:9090")
}
