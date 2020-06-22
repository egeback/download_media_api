package main

import (
	"fmt"
	"time"
	"github.com/egeback/download_media_api/internal/controllers"
	"github.com/egeback/download_media_api/internal/version"
	_ "github.com/egeback/download_media_api/internal/docs"
	

	"github.com/gin-gonic/gin"
	jsonp "github.com/tomwei7/gin-jsonp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Play Media API - Downloader
// @version 1.0
// @description API to download with svt-download

// @contact.name API Support
// @contact.url http://xxxx.xxx.xx
// @contact.email support@egeback.se

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @BasePath /api/v1/
func main() {
	fmt.Printf("%s Running Play Media API - Downloader version: %s (%s)\n", time.Now().Format("2006-01-02 15:04:05"), version.BuildVersion, version.BuildTime)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(jsonp.JsonP())
	c := controllers.NewController()
	v1 := r.Group("/api/v1")
	{
		jobs := v1.Group("/jobs")
		{
			jobs.POST("", c.AddJob)
			jobs.GET("/", c.Jobs)
			jobs.GET("/:uuid", c.GetJob)
		}
		common := v1.Group("/")
		{
			common.GET("ping", c.Ping)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8081")
}