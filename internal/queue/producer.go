package queue

import (
    "encoding/json"
    "time"
    "github.com/hibiken/asynq"
)

type JobPayload struct {
    ID      string `json:"id"`
    Payload string `json:"payload"`
}

func EnqueueJob(client *asynq.Client, payload JobPayload) error {
    data, _ := json.Marshal(payload)
    task := asynq.NewTask("process:job", data)
    _, err := client.Enqueue(task, asynq.MaxRetry(3), asynq.ProcessIn(1*time.Second))
    return err
}