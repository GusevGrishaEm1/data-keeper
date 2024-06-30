package entity

import "time"

// FileRepo stored file
type FileRepo struct {
	// UUID
	UUID string
	// File content
	Content []byte
	// CreatedAt when created
	CreatedAt time.Time
	// CreatedBy created by user
	CreatedBy string
}
