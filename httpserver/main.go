package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "k8s.io/klog/v2"
)

var (
	VERSION           string
	VERSION_SECRET    string
	VERSION_CONFIGMAP string
	srv               *http.Server
)

func init() {
	VERSION = os.Getenv("VERSION")
	VERSION_SECRET = os.Getenv("VERSION_SECRET")
	VERSION_CONFIGMAP = os.Getenv("VERSION_CONFIGMAP")
	rand.Seed(time.Now().UnixNano())
	log.InitFlags(nil)
}

func main() {

	flag.Parse()

	// 测试：幽雅启动
	mockDelay()

	httpServerStart()

	// 测试：幽雅终止
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.V(1).Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatal("Server Shutdown:", err)
			}
			log.V(0).InfoS("Server exiting")
			log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func mockDelay() {
	log.V(1).InfoS("debug", "VERSION", VERSION,
		"VERSION_SECRET", VERSION_SECRET, "VERSION_CONFIGMAP", VERSION_CONFIGMAP)

	// Simulation of start-up elapsed time
	sleep := time.Duration(rand.Intn(5)+1) * time.Second
	log.V(1).InfoS("mock startup", "sleep", sleep)
	time.Sleep(sleep)
}

func httpServerStart() {
	log.V(1).InfoS("start httpserver")
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

	srv = &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func HandleGetAllData(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.V(2).Infof("---body/--- \r\n" + string(body))
	log.V(2).Infof("---header/--- \n")
	for k, v := range c.Request.Header {
		log.V(2).InfoS("Header", k, v)
		// The request header will be written to the response header
		c.Header(k, c.GetHeader(k))
	}

	c.JSON(200, gin.H{
		"hello": "world",
	})
}

func addVersionHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Version", VERSION)
		c.Header("Version-Secret", VERSION_SECRET)
		c.Header("Version-Configmap", VERSION_CONFIGMAP)
		c.Next()
	}
}
