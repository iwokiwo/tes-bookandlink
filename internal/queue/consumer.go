package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"queue-management/internal/store"
	Model "queue-management/model"
	"regexp"
	"time"

	"github.com/hibiken/asynq"
)

// type JobPayloadConsumer struct {
// 	ID    string `json:"id"`
// 	Email string `json:"email"`
// }

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func NewJobHandler() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
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
			}
			return fmt.Errorf("invalid email")
		}

		job, exists := store.GetJob(payload.ID)
		if !exists {
			return nil
		}

		job.Status = "success"
		job.UpdatedAt = time.Now()
		store.SaveJob(job)
		return nil
	}
}
