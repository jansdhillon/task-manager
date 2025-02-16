package main

import "time"

type Task struct {
	Id          int
	Name        string
	Description string
	IsDone      bool      `json:"is_done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
