package store

import (
	"sync"
	"time"
)

type Job struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	jobStore = make(map[string]Job)
	mu       sync.RWMutex
)

func InitMemoryStore() {
	jobStore = make(map[string]Job)
}

func SaveJob(job Job) {
	mu.Lock()
	defer mu.Unlock()
	jobStore[job.ID] = job
}

func GetJob(id string) (Job, bool) {
	mu.RLock()
	defer mu.RUnlock()
	job, ok := jobStore[id]
	return job, ok
}

func ClearJobs() {
	mu.Lock()
	defer mu.Unlock()
	jobStore = make(map[string]Job)
}

func GetAllJobs() []Job {
	mu.RLock()
	defer mu.RUnlock()
	jobs := make([]Job, 0, len(jobStore))
	for _, job := range jobStore {
		jobs = append(jobs, job)
	}
	return jobs
}
