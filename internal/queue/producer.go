package queue

import (
	"encoding/json"
	Model "queue-management/model"
	"time"

	"github.com/hibiken/asynq"
)

// type JobPayloadProducer struct {
// 	ID    string `json:"id"`
// 	Email string `json:"email"`
// }

func EnqueueJob(client *asynq.Client, payload Model.Email) error {
	data, _ := json.Marshal(payload)
	task := asynq.NewTask("process:job", data)
	_, err := client.Enqueue(task, asynq.MaxRetry(3), asynq.ProcessIn(1*time.Second))
	return err
}
