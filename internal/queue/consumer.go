package queue

import (
	"context"
	"encoding/json"
	"log"
	"queue-management/internal/store"
	"time"

	"github.com/hibiken/asynq"
)

func NewJobHandler() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload JobPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			return err
		}

		log.Printf("Processing job: %s", payload.ID)

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
