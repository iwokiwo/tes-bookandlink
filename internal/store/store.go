package store

import (
	"sync"

	Model "queue-management/model"
)

var (
	jobStore = make(map[string]Model.Job)
	mu       sync.RWMutex
)

func InitMemoryStore() {
	jobStore = make(map[string]Model.Job)
}

func SaveJob(job Model.Job) {
	mu.Lock()
	defer mu.Unlock()
	jobStore[job.ID] = job
}

func GetJob(id string) (Model.Job, bool) {
	mu.RLock()
	defer mu.RUnlock()
	job, ok := jobStore[id]
	return job, ok
}

func ClearJobs() {
	mu.Lock()
	defer mu.Unlock()
	jobStore = make(map[string]Model.Job)
}

func GetAllJobs() []Model.Job {
	mu.RLock()
	defer mu.RUnlock()
	jobs := make([]Model.Job, 0, len(jobStore))
	for _, job := range jobStore {
		jobs = append(jobs, job)
	}
	return jobs
}

func EditJob(updatedJob Model.Job) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := jobStore[updatedJob.ID]; !exists {
		return false
	}

	jobStore[updatedJob.ID] = updatedJob
	return true
}
