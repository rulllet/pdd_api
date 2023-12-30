package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"goApi/gen/data"
	"goApi/servicePdd"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var database *data.Queries

type ResponseQuestion struct {
	Question data.Question `json:"question"`
	Answers  []data.Answer `json:"answers"`
}

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	port := os.Getenv("PORT")
	dbName := os.Getenv("NAME")
	timeout := Convert64(os.Getenv("TIMEOUT"))
	database = databaseInit(dbName)

	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(timeoutMiddleware(timeout))
	r.StaticFile("/favicon.ico", "../static/favicon.ico")
	r.StaticFS("/static", http.Dir("../static"))
	r.LoadHTMLGlob("../templates/*")
	index := r.Group("/")
	{
		index.GET("/", func(c *gin.Context) {
			data := TodoPageData{
				PageTitle: "URL list",
				Todos: []Todo{
					{Title: "http://localhost:8080/api/questions/?&ticket=1&question=1", Done: false},
					{Title: "http://localhost:8080/api/questions/random?category=AB", Done: true},
					{Title: "http://localhost:8080/api/questions/random?category=CD", Done: true},
				},
			}
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": data.PageTitle,
				"todos": data.Todos,
			})
		})
	}
	api := r.Group("/api")
	{
		questions := api.Group("/questions")
		questions.GET("/", func(ctx *gin.Context) {
			category := ctx.DefaultQuery("category", "AB")
			ticket := Convert64(ctx.Query("ticket"))
			question := Convert64(ctx.Query("question"))
			res := getDataQuestion(category, ticket, question, ctx)
			ctx.IndentedJSON(http.StatusOK, gin.H{"data": res})
		})
		questions.GET("/random", func(ctx *gin.Context) {
			category := ctx.DefaultQuery("category", "AB")
			res := getDataQuestion(category, rand.Int63n(40), rand.Int63n(20), ctx)
			ctx.IndentedJSON(http.StatusOK, gin.H{"data": res})
		})
		questions.GET("/next", getNextQuestion)
		tickets := api.Group("/tickets")
		tickets.GET("/", func(ctx *gin.Context) {
			category := ctx.DefaultQuery("category", "AB")
			ticket := Convert64(ctx.Query("ticket"))
			ticketData, err := database.GetTicket(ctx, data.GetTicketParams{
				CategoryID: category,
				Ticket:     ticket,
			})
			if err != nil {
				panic(err)
			}
			ctx.IndentedJSON(http.StatusOK, gin.H{"data": ticketData})
		})

	}
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}

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

func getDataQuestion(category string, ticket int64, question int64, c *gin.Context) *ResponseQuestion {
	var questionData data.Question
	questionData, err := database.GetQuestion(c, data.GetQuestionParams{
		CategoryID: category,
		Ticket:     ticket,
		Number:     question,
	})
	if err != nil {
		panic(err)
	}
	var answerData []data.Answer
	answerData, err = database.GetAnswers(c, strconv.FormatInt(questionData.ID, 10))
	if err != nil {
		panic(err)
	}
	res := ResponseQuestion{Question: questionData, Answers: answerData}
	return &res
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
