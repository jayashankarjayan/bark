package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"

	"github.com/techrail/bark/controllers"
	"github.com/techrail/bark/resources"
	"github.com/techrail/bark/services/dbLogWriter"
)

func Index(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("Welcome to Bark!")
}

func Hello(ctx *fasthttp.RequestCtx) {
	_, _ = fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}

// Init performs prerequisite tasks - like loading env variables
func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	Init()
	r := router.New()
	r.GET("/", Index)
	r.GET("/hello/{name}", Hello)
	r.POST("/insertSingle", controllers.SendSingleToChannel)
	r.POST("/insertMultiple", controllers.SendMultipleToChannel)
	r.POST("/shutdownServiceAsap", controllers.ShutdownService)

	err := resources.InitDatabase()
	if err != nil {
		log.Fatal("E#1KDZRP - " + err.Error())
	}
	go dbLogWriter.StartWritingLogs()
	log.Printf("Running on port: %v", os.Getenv("APPPORT"))
	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APPPORT")), r.Handler))
}
