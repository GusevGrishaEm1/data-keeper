package entity

import "time"

// User's stored file
type UserFile struct {
	// File UUID
	UUID string
	// File content
	Content []byte
	// Created at time
	CreatedAt time.Time
	// User who created file
	CreatedBy string
}
