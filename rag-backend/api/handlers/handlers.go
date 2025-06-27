package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fain17/rag-backend/api/models"
	"github.com/fain17/rag-backend/db"
	"github.com/gin-gonic/gin"
	"github.com/pgvector/pgvector-go"
)

// GetHandler godoc
//
//	@Summary	Get file by ID
//	@Tags		files
//	@Produce	json
//	@Param		id	path		int	true	"File ID"
//
// @Success 200 {object} models.FileUploadRequest
//
//	@Failure	404	{object}	map[string]string
//	@Router		/files/{id} [get]
func GetHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		file, err := q.GetFile(c, int32(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		c.JSON(http.StatusOK, file)
	}
}

// UploadHandler godoc
//
//	@Summary		Upload a file
//	@Description	Store a file and its embedding vector
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			file	body   models.FileUploadRequest true "Upload Input"
//
// @Success 200 {object} models.FileUploadRequest
//
//	@Failure		400		{object}	map[string]string
//	@Router			/files [post]
func UploadHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req models.FileUploadRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		vec := pgvector.NewVector(req.Embedding)
		file, err := q.CreateFile(c, db.CreateFileParams{
			Filename:  req.Filename,
			Content:   req.Content,
			Embedding: vec,
		})
		fmt.Print(err)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create file"})
			return
		}
		c.JSON(http.StatusOK, file)
	}
}

// DeleteHandler godoc
//
//	@Summary	Delete a file
//	@Tags		files
//	@Param		id	path		int	true	"File ID"
//	@Success	204	{object}	nil
//	@Failure	400	{object}	map[string]string
//	@Router		/files/{id} [delete]
func DeleteHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		err = q.DeleteFile(c, int32(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// SearchHandler godoc
//
//	@Summary	Search similar files
//	@Tags		search
//	@Accept		json
//	@Produce	json
//	@Param		query	body	models.SearchRequest true "Search Query"
//
// @Success 200 {object} models.SearchRequest
//
//	@Failure	400		{object}	map[string]string
//	@Router		/search [post]
func SearchHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SearchRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		vec := pgvector.NewVector(req.QueryEmbedding)
		results, err := q.SearchFiles(c, db.SearchFilesParams{
			Embedding: vec,
			Limit:     req.Limit,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

// UpdateHandler godoc
//
//	@Summary	Update a file
//	@Tags		files
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int																true	"File ID"
//	@Param		file	body		models.FileUploadRequest true "Update Input"
//
// @Success 200 {object} models.FileUploadRequest
//
//	@Failure	400		{object}	map[string]string
//	@Router		/files/{id} [put]
func UpdateHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var req models.FileUploadRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		vec := pgvector.NewVector(req.Embedding)
		updated, err := q.UpdateFile(c, db.UpdateFileParams{
			ID:        int32(id),
			Filename:  req.Filename,
			Content:   req.Content,
			Embedding: vec,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}
