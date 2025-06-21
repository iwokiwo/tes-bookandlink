package queue

import (
	"encoding/json"
	Model "queue-management/model"
	"time"

	"github.com/hibiken/asynq"
)

//	func EnqueueJob(client *asynq.Client, payload Model.Email) error {
//		data, _ := json.Marshal(payload)
//		task := asynq.NewTask("process:job", data)
//		_, err := client.Enqueue(task, asynq.MaxRetry(3), asynq.ProcessIn(1*time.Second), asynq.Retention(1*time.Hour))
//		return err
//	}
func EnqueueJob(client *asynq.Client, payload Model.Email) (string, error) {
	data, _ := json.Marshal(payload)
	task := asynq.NewTask("process:job", data)

	info, err := client.Enqueue(task,
		asynq.MaxRetry(1),
		asynq.ProcessIn(1*time.Second),
		asynq.Retention(1*time.Hour), // penting agar bisa dihapus, simpan di radis semua data
	)
	if err != nil {
		return "", err
	}
	return info.ID, nil
}
