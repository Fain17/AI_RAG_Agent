package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fain17/rag-backend/api/models"
	"github.com/fain17/rag-backend/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
		id := c.Param("id")
		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var dbUUID pgtype.UUID
		if err := dbUUID.Scan(parsedUUID.String()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert UUID"})
			return
		}

		file, err := q.GetFile(c, dbUUID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		c.JSON(http.StatusOK, file)
	}
}

// GetAllHandler godoc
//
//	@Summary		Get all files
//	@Tags			files
//	@Produce		json
//
//	@Success		200	{array}	models.FileUploadRequest
//	@Failure		500	{object}	map[string]string
//	@Router			/files [get]
func GetAllHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		files, err := q.GetAllFiles(c)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		c.JSON(http.StatusOK, files)
	}
}

// GetFilesByFilenameHandler godoc
//
//	@Summary	Search files by filename
//	@Tags		files
//	@Produce	json
//	@Param		query	query	string	true	"Search keyword"
//	@Success	200		{array}	models.FileUploadRequest
//	@Failure	400		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/files/search [get]
func GetFilesByFilenameHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("query")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
			return
		}

		files, err := q.GetFilesByFilename(c, pgtype.Text{String: query, Valid: true})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
			return
		}

		c.JSON(http.StatusOK, files)
	}
}

// GetFilesByDateRangeHandler godoc
//
//	@Summary	Get files within a date range
//	@Tags		files
//	@Produce	json
//	@Param		start	query	string	true	"Start date (YYYY-MM-DD)"
//	@Param		end		query	string	true	"End date (YYYY-MM-DD)"
//	@Success	200		{array}	models.FileUploadRequest
//	@Failure	400		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/files/date-range [get]
func GetFilesByDateRangeHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := c.Query("start")
		end := c.Query("end")

		startDate, err := time.Parse("2006-01-02", start)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
			return
		}

		endDate, err := time.Parse("2006-01-02", end)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
			return
		}

		var startTS, endTS pgtype.Timestamptz

		_ = startTS.Scan(startDate)
		_ = endTS.Scan(endDate)

		params := db.GetFilesByDateRangeParams{
			CreatedAt:   startTS,
			CreatedAt_2: endTS,
		}

		files, err := q.GetFilesByDateRange(c, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get files by date"})
			return
		}

		c.JSON(http.StatusOK, files)
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
		id := c.Param("id")
		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var dbUUID pgtype.UUID
		if err := dbUUID.Scan(parsedUUID.String()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert UUID"})
			return
		}

		err = q.DeleteFile(c, dbUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
			return
		}

		c.Status(http.StatusNoContent)
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
		id := c.Param("id")
		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var dbUUID pgtype.UUID
		if err := dbUUID.Scan(parsedUUID.String()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert UUID"})
			return
		}

		var req models.FileUploadRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		vec := pgvector.NewVector(req.Embedding)
		updated, err := q.UpdateFile(c, db.UpdateFileParams{
			ID:        dbUUID,
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

// SoftDeleteHandler godoc
//
//	@Summary	Soft delete a file
//	@Description	Marks a file as deleted without removing it from the database.
//	@Tags			files
//	@Param			id	path	string	true	"File UUID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/files/soft-delete/{id} [delete]
func SoftDeleteHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		parsedUUID, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
			return
		}

		var dbUUID pgtype.UUID
		if err := dbUUID.Scan(parsedUUID.String()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID conversion failed"})
			return
		}

		err = q.SoftDeleteFile(c, dbUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not soft delete file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "file soft-deleted successfully"})
	}
}

// UndoSoftDeleteHandler godoc
//
//	@Summary	Restore a soft-deleted file
//	@Description	Sets the file's deleted flag back to false.
//	@Tags			files
//	@Param			id	path	string	true	"File UUID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/files/restore/{id} [patch]
func UndoSoftDeleteHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		parsedUUID, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
			return
		}

		var dbUUID pgtype.UUID
		if err := dbUUID.Scan(parsedUUID.String()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID conversion failed"})
			return
		}

		err = q.UndoSoftDelete(c, dbUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not restore file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "file restored successfully"})
	}
}

// GetDeletedFilesHandler godoc
//
//	@Summary	Get all soft-deleted files
//	@Description	Retrieves files that have been soft-deleted (Recycle Bin).
//	@Tags			files
//	@Produce		json
//	@Success		200	{array}	models.FileUploadRequest
//	@Failure		500	{object}	map[string]string
//	@Router			/files/recycle-bin [get]
func GetDeletedFilesHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		files, err := q.GetDeletedFiles(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch deleted files"})
			return
		}

		c.JSON(http.StatusOK, files)
	}
}

// GetFileMetadataHandler godoc
//
//	@Summary	Get lightweight file metadata
//	@Tags		files
//	@Produce	json
//	@Success	200	{array}	models.FileMetadata
//	@Failure	500	{object}	map[string]string
//	@Router		/files/metadata [get]
func GetFileMetadataHandler(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		files, err := q.GetFileMetadata(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get metadata"})
			return
		}

		c.JSON(http.StatusOK, files)
	}
}
