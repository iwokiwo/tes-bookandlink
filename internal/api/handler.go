package api

import (
	"net/http"
	"queue-management/internal/queue"
	"queue-management/internal/store"
	Model "queue-management/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

var AsynqClient *asynq.Client

type JobRequest struct {
	Email string `json:"email"`
	URL   string `json:"url"`
}

func PostJob(c *gin.Context) {
	var req JobRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	job := store.Job{
		ID:        id,
		Email:     req.Email,
		URL:       req.URL,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	store.SaveJob(job)

	err := queue.EnqueueJob(AsynqClient, Model.Email{ID: id, Email: req.Email, URL: req.URL})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func GetJobs(c *gin.Context) {
	jobs := store.GetAllJobs()
	c.JSON(http.StatusOK, jobs)
}

func RetryJob(c *gin.Context) {
	id := c.Param("id")
	var req JobRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, exists := store.GetJob(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	job.Email = req.Email
	job.URL = req.URL
	job.Status = "pending"
	job.UpdatedAt = time.Now()
	store.SaveJob(job)

	err := queue.EnqueueJob(AsynqClient, Model.Email{ID: job.ID, Email: job.Email, URL: req.URL})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "re-enqueue failed"})
		return
	}
	c.JSON(http.StatusOK, job)
}

func ClearJobList(c *gin.Context) {
	store.ClearJobs()
	c.JSON(http.StatusOK, gin.H{"message": "all jobs cleared"})
}
