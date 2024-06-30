package entity

import "time"

// UserFile stored file
type UserFile struct {
	// UUID
	UUID string
	// File content
	Content []byte
	// CreatedAt when created
	CreatedAt time.Time
	// CreatedBy created by user
	CreatedBy string
}
