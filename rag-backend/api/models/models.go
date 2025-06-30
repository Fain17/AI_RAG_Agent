package models

import "time"

// @Param	file	body	FileUploadRequest	true	"Upload data"
type FileUploadRequest struct {
	Filename  string    `json:"filename"`
	Content   string    `json:"content"`
	Embedding []float32 `json:"embedding"`
	CreatedAt time.Time `json:"created_at"`
	Deleted   bool      `json:"deleted"`
}

// @Param	file	body	FileMetaData	true	"Upload data"
type FileMetadata struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Size      int       `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
