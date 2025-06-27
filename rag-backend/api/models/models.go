package models

//	@Param	file	body	FileUploadRequest	true	"Upload data"
type FileUploadRequest struct {
	Filename  string    `json:"filename"`
	Content   string    `json:"content"`
	Embedding []float32 `json:"embedding"`
}

//	@Param	query	body	SearchRequest	true	"Search query"
type SearchRequest struct {
	QueryEmbedding []float32 `json:"query_embedding"`
	Limit          int32     `json:"limit"`
}
