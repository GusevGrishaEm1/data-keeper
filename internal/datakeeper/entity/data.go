package entity

import "time"

// ContentType Data type
type ContentType string

const (
	// LogPass Login/password
	LogPass ContentType = "LOG_PASS"
	// File Binary/Text file (file size < 5 MB)
	File ContentType = "File"
)

// Data User's stored data
type Data struct {
	// UUID
	UUID string
	// Content data
	Content []byte
	// ContentType content type
	ContentType ContentType
	// CreatedAt Created at time
	CreatedAt time.Time
	// CreatedBy User who created this data
	CreatedBy string
}
