package store

import (
	"log"

	"github.com/hibiken/asynq"
)

var asynqRedisOpt asynq.RedisClientOpt

func DeleteTaskFromRedis(taskID string) error {
	inspector := asynq.NewInspector(asynqRedisOpt)
	err := inspector.DeleteTask("default", taskID)
	if err != nil {
		log.Printf("Gagal hapus task %s: %v", taskID, err)
		return err
	} else {
		log.Printf("Task %s berhasil dihapus dari Redis", taskID)
		return nil
	}
}
