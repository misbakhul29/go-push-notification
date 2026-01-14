package models

import "time"

type NotificationJob struct {
	ID           string    `json:"id"`
	TargetUserID string    `json:"target_user_id"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at"`
}
