package main

import (
	"queue-management/internal/api"
	"queue-management/internal/logstream"
	"queue-management/internal/queue"
	"queue-management/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

const MaxConcurrency = 3

func main() {
	store.InitMemoryStore()
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"POST, GET, OPTIONS, PUT, DELETE"}
	config.AllowHeaders = []string{"Authorization", "Content-Type", "X-CSRF-Token"}
	r.Use(cors.New(config))

	redisConn := asynq.RedisClientOpt{Addr: "localhost:6379"}
	api.AsynqClient = asynq.NewClient(redisConn)

	r.POST("/api/queue", api.PostJob)
	r.GET("/api/queue", api.GetJobs)
	r.POST("/api/queue/:id/retry", api.RetryJob)
	r.DELETE("/api/queue", api.ClearJobList)
	r.GET("/ws/logs", logstream.Handler)

	go startWorker(redisConn)

	r.Run(":8080")
}

// func startWorker(redis asynq.RedisClientOpt) {
// 	srv := asynq.NewServer(redis, asynq.Config{Concurrency: 10})
// 	mux := asynq.NewServeMux()
// 	mux.HandleFunc("process:job", queue.NewJobHandler())
// 	if err := srv.Run(mux); err != nil {
// 		panic(err)
// 	}
// }

func startWorker(redis asynq.RedisClientOpt) {
	srv := asynq.NewServer(redis, asynq.Config{
		Concurrency: MaxConcurrency, //Gunakan batas maksimal di sini
		Queues: map[string]int{
			"default": 1, // Optional: atur prioritas queue jika kamu pakai lebih dari satu
		},
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc("process:job", queue.NewJobHandler())

	if err := srv.Run(mux); err != nil {
		panic(err)
	}
}
