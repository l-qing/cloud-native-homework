package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	VERSION string
)

func init() {
	VERSION = os.Getenv("VERSION")
}

func main() {
	// Creates a router without any middleware by default
	r := gin.New()

	// avoid auto redirection when path is like 'http://localhost:8080/123456'
	r.RedirectTrailingSlash = false

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	r.Use(addVersionHeader())

	// r.GET("/", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "ohai")
	// })

	r.NoRoute(func(c *gin.Context) {
		HandleGetAllData(c)
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "ok")
	})

	r.Run(":8081")
}

func HandleGetAllData(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("---body/--- \r\n" + string(body))
	fmt.Println("---header/--- ")
	for k, v := range c.Request.Header {
		fmt.Println(k, v)
		// The request header will be written to the response header
		c.Header(k, c.GetHeader(k))
	}

	c.JSON(200, gin.H{
		"hello": "world",
	})
}

func addVersionHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("VERSION", VERSION)
		c.Next()
	}
}
