package api

import (
	"net/http"
	"queue-management/internal/queue"
	"queue-management/internal/store"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

var AsynqClient *asynq.Client

func PostJob(c *gin.Context) {
	var req struct {
		Payload string `json:"payload"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	job := store.Job{
		ID:        id,
		Payload:   req.Payload,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	store.SaveJob(job)

	err := queue.EnqueueJob(AsynqClient, queue.JobPayload{ID: id, Payload: req.Payload})
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
	job, exists := store.GetJob(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	job.Status = "pending"
	job.UpdatedAt = time.Now()
	store.SaveJob(job)

	err := queue.EnqueueJob(AsynqClient, queue.JobPayload{ID: job.ID, Payload: job.Payload})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "re-enqueue failed"})
		return
	}
	c.JSON(http.StatusOK, job)
}
