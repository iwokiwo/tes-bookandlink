package Model

import "time"

type Email struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

type Job struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TaskID    string    `json:"task_id"`
}

type JobRequest struct {
	Email string `json:"email"`
	URL   string `json:"url"`
}
