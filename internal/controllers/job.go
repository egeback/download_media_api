package controllers

import (
	"net/http"
	"net/url"

	"github.com/google/uuid"

	"github.com/egeback/download_media_api/internal/actions"
	"github.com/gin-gonic/gin"
)

//Job ...
type Job struct {
	UUID     string
	Download *actions.Downloader
}

var jobs = make(map[string]Job)

// AddJob ...
func (c *Controller) AddJob(ctx *gin.Context) {
	u := ctx.DefaultQuery("url", "")
	_, err := url.ParseRequestURI(u)

	if u == "" || err != nil {
		c.createErrorResponse(ctx, http.StatusBadRequest, 100, "no valid url provided")
		return
	}

	download := actions.AddDownload(u)
	//download := actions.AddDownload("https://www.svtplay.se/video/21868842/palmegruppen-tar-langlunch")
	//download := actions.AddDownload("https://www.svtplay.se/video/26987573/you-were-never-really-here")
	download.Start()
	id := uuid.New()
	uuid := id.String()
	jobs[uuid] = Job{UUID: uuid, Download: &download}

	ctx.JSON(http.StatusAccepted, gin.H{
		"job_id": uuid,
	})
	return
}

//Jobs ...
func (c *Controller) Jobs(ctx *gin.Context) {
	retVal := make([]Job, 0, len(jobs))
	for _, job := range jobs {
		retVal = append(retVal, job)
	}
	ctx.JSON(http.StatusOK, retVal)
}

//GetJob ...
func (c *Controller) GetJob(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	_, exists := jobs[uuid]
	if !exists {
		//ctx.JSON(http.StatusNotFound, gin.H{})
		c.createErrorResponse(ctx, http.StatusNotFound, 101, "job id does not exist")
		return
	}

	ctx.JSON(http.StatusAccepted, jobs[uuid])
	return
}
