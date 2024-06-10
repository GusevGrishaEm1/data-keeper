package entity

import "time"

// Data type
type ContentType string

const (
	// Credit card
	CARD ContentType = "CARD"
	// Login/password
	LOG_PASS ContentType = "LOG_PASS"
	// Binary/Text file (file size < 5 MB)
	FILE ContentType = "FILE"
)

// User's stored data
type Data struct {
	// Data UUID
	UUID string
	// Info about data
	Content []byte
	// Content type
	ContentType ContentType
	// Created at time
	CreatedAt time.Time
	// User who created this data
	CreatedBy string
}
