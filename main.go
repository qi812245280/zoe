package main

import (
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
	"zoe/config"
	"zoe/controller"
	"zoe/dao/db"
)

var (
	AppVersion   = "0.0.1"
	AppBuildTime = "2017-12-01T00:03:18+0800"
	AppGitHash   = "undefined"
)

func InfoHandler(c *gin.Context) {
	log.Info(c.Params)
	c.JSON(http.StatusOK, gin.H{
		"name":    "guldan",
		"version": AppVersion,
	})

}

func applyRoute(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	v1 := r.Group("/api")
	v1.GET("/info", InfoHandler)

	v1.POST("/org", controller.CreateOrgHandler)
	v1.PUT("/org/:org_id", controller.UpdateOrgHandler)
	v1.DELETE("/org/:org_id", controller.DeleteOrgHandler)
	v1.GET("/org", controller.ListOrgHandler)
	v1.GET("/org/:org_id", controller.SingleOrgHandler)
	v1.POST("/org/:org_id/authorize", controller.AuthorizeOrgHandler)
	v1.DELETE("/org/:org_id/authorize/:user_id", controller.DeleteAuthorizeOrgHandler)

	v1.PUT("/project", controller.CreateProjectHandler)
	v1.POST("/project/:project_id", controller.UpdateProjectHandler)
}

func guldanAccessLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		log.Infof("[GIN] \"%s %s\" %d %v %vus", method, path, statusCode, clientIP, latency.Nanoseconds()/1000.0)
	}
}

func main() {
	defer log.Flush()
	configFile := flag.String("config", "./config.yaml", "config file")
	version := flag.Bool("version", false, "print current version")
	help := flag.Bool("help", false, "show help")
	flag.Parse()
	if *version {
		fmt.Printf("Version: %v\n", AppVersion)
		fmt.Printf("Git Hash: %v\n", AppGitHash)
		fmt.Printf("Build Time: %v\n", AppBuildTime)
		os.Exit(0)
	}
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if err := config.LoadConfig(*configFile); err != nil {
		fmt.Printf("load %v fail: %v", *configFile, err.Error())
		os.Exit(1)
	}

	{
		logger, err := log.LoggerFromConfigAsFile(config.C.LogFormat)
		if err != nil {
			fmt.Printf("load %v fail: %v", config.C.LogFormat, err.Error())
			os.Exit(1)
		}
		_ = log.ReplaceLogger(logger)
	}

	if err := db.InitMysql(config.C); err != nil {
		_ = log.Criticalf("new middleware fail: %v", err)
		os.Exit(1)
	}
	defer db.Destroy()

	if config.C.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	}

	r := gin.New()
	r.Use(guldanAccessLogger())
	r.Use(gin.Recovery())
	pprof.Register(r) // 性能

	applyRoute(r)

	log.Infof("Listening and serving HTTP on %v", config.C.Listen)
	if err := r.Run(config.C.Listen); err != nil {
		_ = log.Errorf("http listen fail: %v", err)
	}
}
