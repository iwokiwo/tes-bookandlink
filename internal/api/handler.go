package api

import (
	"log"
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

func PostJob(c *gin.Context) {
	var req Model.JobRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	job := Model.Job{
		ID:        id,
		Email:     req.Email,
		URL:       req.URL,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	taskID, err := queue.EnqueueJob(AsynqClient, Model.Email{ID: id, Email: req.Email, URL: req.URL})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue"})
		return
	}
	job.TaskID = taskID
	store.SaveJob(job)

	c.JSON(http.StatusOK, job)
}

func GetJobs(c *gin.Context) {
	jobs := store.GetAllJobs()
	c.JSON(http.StatusOK, jobs)
}

func RetryJob(c *gin.Context) {
	id := c.Param("id")
	var req Model.JobRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, exists := store.GetJob(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	err := store.DeleteTaskFromRedis(job.TaskID)
	if err != nil {
		log.Printf("Gagal hapus task %s: %v", job.TaskID, err)
	} else {
		log.Printf("Task %s berhasil dihapus dari Redis", job.TaskID)
	}

	job.Email = req.Email
	job.URL = req.URL
	job.Status = "pending"
	job.UpdatedAt = time.Now()
	store.EditJob(job)

	taskID, err := queue.EnqueueJob(AsynqClient, Model.Email{ID: job.ID, Email: job.Email, URL: req.URL})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "re-enqueue failed"})
		return
	}
	job.TaskID = taskID
	store.EditJob(job)

	c.JSON(http.StatusOK, job)
}

func ClearJobList(c *gin.Context) {
	// 1. Ambil semua job dari memory store
	jobs := store.GetAllJobs() // Kamu harus buat fungsi ini di store.go

	// 2. Loop untuk hapus task dari Redis
	for _, job := range jobs {
		if job.TaskID != "" {
			err := store.DeleteTaskFromRedis(job.TaskID)
			if err != nil {
				log.Printf("Gagal hapus task %s: %v", job.TaskID, err)
			} else {
				log.Printf("Task %s berhasil dihapus dari Redis", job.TaskID)
			}
		}
	}
	store.ClearJobs()
	c.JSON(http.StatusOK, gin.H{"message": "all jobs cleared"})
}
