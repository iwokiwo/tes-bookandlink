package main

import (
	"queue-management/internal/api"
	"queue-management/internal/queue"
	"queue-management/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

func main() {
	store.InitMemoryStore()
	r := gin.Default()

	redisConn := asynq.RedisClientOpt{Addr: "localhost:6379"}
	api.AsynqClient = asynq.NewClient(redisConn)

	r.POST("/api/queue", api.PostJob)
	r.GET("/api/queue", api.GetJobs)
	r.POST("/api/queue/:id/retry", api.RetryJob)

	go startWorker(redisConn)

	r.Run(":8080")
}

func startWorker(redis asynq.RedisClientOpt) {
	srv := asynq.NewServer(redis, asynq.Config{Concurrency: 10})
	mux := asynq.NewServeMux()
	mux.HandleFunc("process:job", queue.NewJobHandler())
	if err := srv.Run(mux); err != nil {
		panic(err)
	}
}
