package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"queue-management/internal/store"
	Model "queue-management/model"
	"regexp"
	"time"

	"github.com/hibiken/asynq"
)

func checkAndDeleteFailedTask(ctx context.Context, taskID string) {
	if taskID == "" {
		return
	}

	retry, ok1 := asynq.GetRetryCount(ctx)
	max, ok2 := asynq.GetMaxRetry(ctx)

	if ok1 && ok2 && retry >= max {
		log.Printf("Retry habis (%d/%d), hapus task %s", retry, max, taskID)
		go func() {
			time.Sleep(2 * time.Second)
			store.DeleteTaskFromRedis(taskID)
		}()
	}
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func pingURL(url string) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func NewJobHandler() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		log.Println("Simulating network delay...")
		time.Sleep(1 * time.Second) // delay simulasi lambat
		log.Println("Finished simulated delay.")

		var payload Model.Email
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			return err
		}

		log.Printf("Processing job: %s", payload.ID)

		if !isValidEmail(payload.Email) {
			log.Printf("Invalid email format: %s", payload.Email)
			job, exists := store.GetJob(payload.ID)
			if exists {
				job.Status = "failed"
				job.UpdatedAt = time.Now()
				store.SaveJob(job)
				//logstream.BroadcastJob(job) // Send job update websocket
			}
			return fmt.Errorf("invalid email")
		}

		// Ping ke URL yang dikirim dari Postman
		if err := pingURL(payload.URL); err != nil {
			log.Printf("Ping failed: %v", err)
			if job, exists := store.GetJob(payload.ID); exists {
				job.Status = "failed"
				job.UpdatedAt = time.Now()
				checkAndDeleteFailedTask(ctx, job.TaskID)
				store.SaveJob(job)
				//logstream.BroadcastJob(job) // Send job update websocket
			}
			return fmt.Errorf("ping failed: %v", err)
		}

		job, exists := store.GetJob(payload.ID)
		if !exists {
			return nil
		}

		job.Status = "success"
		job.URL = payload.URL
		job.UpdatedAt = time.Now()
		store.SaveJob(job)
		return nil
	}
}
