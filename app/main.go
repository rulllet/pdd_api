package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"goApi/gen/data"
	"goApi/servicePdd"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var database *data.Queries

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Print("No .env file found")
	}
}

func databaseInit(dbName string) *data.Queries {
	db, err := sql.Open("sqlite3", "../"+dbName)
	if err != nil {
		panic(err)
	}
	queries := data.New(db)
	return queries
}

func routerInit(r *gin.Engine) *gin.Engine {
	r.StaticFile("/favicon.ico", "../static/favicon.ico")
	api := r.Group("/api")
	{
		//api.StaticFS("/static", http.Dir("../static"))
		api.GET("/", getDataQuestion)
		api.GET("/next", getNextQuestion)
	}
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return r
}

func getDataQuestion(c *gin.Context) {
	category := c.DefaultQuery("category", "AB")
	ticket := Convert64(c.Query("ticket"))
	question := Convert64(c.Query("question"))
	questionData, err := database.GetQuestion(c, data.GetQuestionParams{
		CategoryID: category,
		Ticket:     ticket,
		Number:     question,
	})
	if err != nil {
		panic(err)
	}

	answerData, err := database.GetAnswers(c, strconv.FormatInt(questionData.ID, 10))
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, gin.H{"question": questionData, "answers": answerData})
}

func getNextQuestion(c *gin.Context) {
	category := c.DefaultQuery("category", "AB")
	var body servicePdd.DataTicket
	c.Bind(&body)
	nextQ := servicePdd.NextQuestion(body)
	questionData, err := database.GetQuestion(c, data.GetQuestionParams{
		CategoryID: category,
		Ticket:     nextQ.Ticket,
		Number:     nextQ.Question,
	})
	if err != nil {
		panic(err)
	}
	answerData, err := database.GetAnswers(c, strconv.FormatInt(questionData.ID, 10))
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": nextQ, "question": questionData, "answers": answerData})
}

func Convert64(value string) int64 {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return result
}

func timeoutResponse(c *gin.Context) {
	c.IndentedJSON(http.StatusRequestTimeout, gin.H{"timeout": http.StatusRequestTimeout})
}

func timeoutMiddleware(timing int64) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(time.Duration(timing)*time.Millisecond),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(timeoutResponse),
	)
}

func main() {
	port := os.Getenv("PORT")
	dbName := os.Getenv("NAME")
	timeout := Convert64(os.Getenv("TIMEOUT"))

	database = databaseInit(dbName)
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(timeoutMiddleware(timeout))
	r.Use(static.Serve("/static", static.LocalFile("../static", true)))
	router := routerInit(r)
	if err := router.Run(port); err != nil {
		log.Fatal(err)
	}
}
