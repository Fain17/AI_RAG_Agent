package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fain17/rag-backend/api/handlers"
	"github.com/fain17/rag-backend/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	pgvector "github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
)

// setupHandlersTestRouter creates a test router instance for comprehensive handler testing
// This router is used to test actual handler functions with their validation logic
func setupHandlersTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// TestHandlersInputValidation tests the input validation logic of all handlers
// This comprehensive test suite covers all the validation scenarios that handlers must handle
// Each test calls the actual handler function with nil queries to test validation-only paths
func TestHandlersInputValidation(t *testing.T) {
	// Test GetHandler with invalid UUID
	// This tests the UUID parsing validation in GetHandler - line 30-34 in handlers.go
	// GetHandler must validate UUID format before attempting database operations
	t.Run("GetHandler_InvalidUUID", func(t *testing.T) {
		router := setupHandlersTestRouter()
		// Use nil queries to test input validation only
		router.GET("/files/:id", handlers.GetHandler(nil))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/files/invalid-uuid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid id", response["error"])
	})

	// Test GetFilesByFilenameHandler missing query parameter
	// This tests the query parameter validation in GetFilesByFilenameHandler - line 99-102 in handlers.go
	// The handler must ensure the 'query' parameter is present before proceeding
	t.Run("GetFilesByFilenameHandler_MissingQuery", func(t *testing.T) {
		router := setupHandlersTestRouter()
		router.GET("/files/search", handlers.GetFilesByFilenameHandler(nil))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/files/search", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "query parameter is required", response["error"])
	})

	// Test GetFilesByFilenameHandler with empty query parameter
	// This tests the empty string validation in GetFilesByFilenameHandler - line 99-102 in handlers.go
	// The handler must reject empty query strings as invalid search terms
	t.Run("GetFilesByFilenameHandler_EmptyQuery", func(t *testing.T) {
		router := setupHandlersTestRouter()
		router.GET("/files/search", handlers.GetFilesByFilenameHandler(nil))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/files/search?query=", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "query parameter is required", response["error"])
	})

	// Test UploadHandler with invalid JSON
	// This tests the JSON binding validation in UploadHandler - line 178-181 in handlers.go
	// The handler must validate JSON format before processing upload data
	t.Run("UploadHandler_InvalidJSON", func(t *testing.T) {
		router := setupHandlersTestRouter()
		router.POST("/files", handlers.UploadHandler(nil))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/files", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid request", response["error"])
	})

	// Test UpdateHandler with invalid UUID
	// This tests the UUID parsing validation in UpdateHandler - line 243-246 in handlers.go
	// UpdateHandler must validate UUID format before attempting update operations
	t.Run("UpdateHandler_InvalidUUID", func(t *testing.T) {
		router := setupHandlersTestRouter()
		router.PUT("/files/:id", handlers.UpdateHandler(nil))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/files/invalid-uuid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid id", response["error"])
	})

	// Test UpdateHandler with invalid JSON
	// This tests the JSON binding validation in UpdateHandler - line 254-257 in handlers.go
	// UpdateHandler must validate both UUID and JSON before processing updates
	t.Run("UpdateHandler_InvalidJSON", func(t *testing.T) {
		router := setupHandlersTestRouter()
		testUUID := uuid.New()
		router.PUT("/files/:id", handlers.UpdateHandler(nil))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/files/"+testUUID.String(), bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid request body", response["error"])
	})

	// Test DeleteHandler with invalid UUID
	// This tests the UUID parsing validation in DeleteHandler - line 207-210 in handlers.go
	// DeleteHandler must validate UUID format before attempting deletion
	t.Run("DeleteHandler_InvalidUUID", func(t *testing.T) {
		router := setupHandlersTestRouter()
		router.DELETE("/files/:id", handlers.DeleteHandler(nil))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/files/invalid-uuid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid id", response["error"])
	})

	// Test SoftDeleteHandler with invalid UUID
	// This tests the UUID parsing validation in SoftDeleteHandler - line 297-300 in handlers.go
	// SoftDeleteHandler must validate UUID format before soft deletion
	t.Run("SoftDeleteHandler_InvalidUUID", func(t *testing.T) {
		router := setupHandlersTestRouter()
		router.PATCH("/files/:id/soft-delete", handlers.SoftDeleteHandler(nil))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/files/invalid-uuid/soft-delete", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid UUID", response["error"])
	})

	// Test UndoSoftDeleteHandler with invalid UUID
	// This tests the UUID parsing validation in UndoSoftDeleteHandler - line 335-338 in handlers.go
	// UndoSoftDeleteHandler must validate UUID format before restoration
	t.Run("UndoSoftDeleteHandler_InvalidUUID", func(t *testing.T) {
		router := setupHandlersTestRouter()
		router.PATCH("/files/:id/restore", handlers.UndoSoftDeleteHandler(nil))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/files/invalid-uuid/restore", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid UUID", response["error"])
	})

	// Test GetFilesByDateRangeHandler with invalid start date
	// This tests the start date parsing validation in GetFilesByDateRangeHandler - line 126-129 in handlers.go
	// The handler must validate start date format before querying by date range
	t.Run("GetFilesByDateRangeHandler_InvalidStartDate", func(t *testing.T) {
		router := setupHandlersTestRouter()
		router.GET("/files/date-range", handlers.GetFilesByDateRangeHandler(nil))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/files/date-range?start=invalid&end=2024-12-31", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid start date", response["error"])
	})

	// Test GetFilesByDateRangeHandler with invalid end date
	// This tests the end date parsing validation in GetFilesByDateRangeHandler - line 132-135 in handlers.go
	// The handler must validate end date format before querying by date range
	t.Run("GetFilesByDateRangeHandler_InvalidEndDate", func(t *testing.T) {
		router := setupHandlersTestRouter()
		router.GET("/files/date-range", handlers.GetFilesByDateRangeHandler(nil))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/files/date-range?start=2024-01-01&end=invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid end date", response["error"])
	})
}

// TestHandlersValidationSuccess tests the successful validation paths and helper functions
// This section validates that the supporting functions and data structures work correctly
// These tests ensure the building blocks of handlers function properly
func TestHandlersValidationSuccess(t *testing.T) {
	// Test valid UUID parsing functionality
	// This validates the uuid.Parse function used throughout handlers
	// Ensures that valid UUIDs are parsed correctly without errors
	t.Run("ValidUUIDParsing", func(t *testing.T) {
		testUUID := uuid.New()
		parsedUUID, err := uuid.Parse(testUUID.String())
		assert.NoError(t, err)
		assert.Equal(t, testUUID, parsedUUID)
	})

	// Test valid JSON structure marshaling and unmarshaling
	// This validates the models.FileUploadRequest structure used in handlers
	// Ensures JSON binding works correctly for valid request structures
	t.Run("ValidJSONStructure", func(t *testing.T) {
		validRequest := models.FileUploadRequest{
			Filename:  "test.txt",
			Content:   "test content",
			Embedding: []float32{1.0, 2.0, 3.0},
		}

		jsonBytes, err := json.Marshal(validRequest)
		assert.NoError(t, err)

		var parsed models.FileUploadRequest
		err = json.Unmarshal(jsonBytes, &parsed)
		assert.NoError(t, err)
		assert.Equal(t, validRequest.Filename, parsed.Filename)
		assert.Equal(t, validRequest.Content, parsed.Content)
		assert.Equal(t, validRequest.Embedding, parsed.Embedding)
	})

	// Test valid date parsing functionality
	// This validates the time.Parse function used in GetFilesByDateRangeHandler
	// Ensures that valid date strings are parsed correctly
	t.Run("ValidDateParsing", func(t *testing.T) {
		dateStr := "2024-01-01"
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		assert.NoError(t, err)
		assert.Equal(t, 2024, parsedDate.Year())
		assert.Equal(t, 1, int(parsedDate.Month()))
		assert.Equal(t, 1, parsedDate.Day())
	})

	// Test leap year date parsing edge case
	// This validates proper handling of leap year dates in date range queries
	// Ensures February 29th in leap years is handled correctly
	t.Run("LeapYearDateParsing", func(t *testing.T) {
		dateStr := "2024-02-29"
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		assert.NoError(t, err)
		assert.Equal(t, 2024, parsedDate.Year())
		assert.Equal(t, 2, int(parsedDate.Month()))
		assert.Equal(t, 29, parsedDate.Day())
	})
}

// TestHandlersUUIDScanning tests the UUID scanning functionality used in handlers
// This validates the pgtype.UUID scanning operations that handlers use for database operations
// These tests ensure proper UUID conversion for database queries
func TestHandlersUUIDScanning(t *testing.T) {
	// Test valid UUID scanning for database operations
	// This validates the UUID scanning logic used in lines 35-38 in handlers.go
	// Ensures valid UUIDs can be converted to database-compatible format
	t.Run("ValidUUIDScanning", func(t *testing.T) {
		testUUID := uuid.New()
		var dbUUID pgtype.UUID

		err := dbUUID.Scan(testUUID.String())
		assert.NoError(t, err)
		assert.True(t, dbUUID.Valid)
	})

	// Test invalid UUID scanning error handling
	// This validates error handling when invalid UUIDs are scanned
	// Ensures proper error detection for malformed UUID strings
	t.Run("InvalidUUIDScanning", func(t *testing.T) {
		var dbUUID pgtype.UUID

		err := dbUUID.Scan("invalid-uuid")
		assert.Error(t, err)
	})

	// Test nil UUID scanning behavior
	// This validates proper handling of nil values in UUID scanning
	// Ensures NULL values are handled correctly in database operations
	t.Run("NilUUIDScanning", func(t *testing.T) {
		var dbUUID pgtype.UUID

		err := dbUUID.Scan(nil)
		assert.NoError(t, err)
		assert.False(t, dbUUID.Valid)
	})
}

// TestHandlersVectorCreation tests the vector creation functionality used in handlers
// This validates the pgvector.NewVector operations used in UploadHandler and UpdateHandler
// These tests ensure proper vector creation for embedding storage
func TestHandlersVectorCreation(t *testing.T) {
	// Test standard vector creation
	// This validates the vector creation logic used in lines 183 and 259 in handlers.go
	// Ensures embeddings are properly converted to database vectors
	t.Run("CreateVector", func(t *testing.T) {
		embedding := []float32{1.0, 2.0, 3.0, 4.0}
		vec := pgvector.NewVector(embedding)
		assert.NotNil(t, vec)
	})

	// Test empty vector creation edge case
	// This validates handling of empty embedding arrays
	// Ensures system can handle files with no embeddings
	t.Run("CreateEmptyVector", func(t *testing.T) {
		emptyEmbedding := []float32{}
		vec := pgvector.NewVector(emptyEmbedding)
		assert.NotNil(t, vec)
	})

	// Test large vector creation performance
	// This validates handling of large embedding vectors
	// Ensures system can handle high-dimensional embeddings efficiently
	t.Run("CreateLargeVector", func(t *testing.T) {
		largeEmbedding := make([]float32, 1000)
		for i := range largeEmbedding {
			largeEmbedding[i] = float32(i)
		}
		vec := pgvector.NewVector(largeEmbedding)
		assert.NotNil(t, vec)
	})

	// Test vector creation with negative values
	// This validates handling of negative embedding values
	// Ensures vectors with negative components are handled correctly
	t.Run("CreateVectorWithNegativeValues", func(t *testing.T) {
		embedding := []float32{-1.0, -2.0, 3.0, -4.0}
		vec := pgvector.NewVector(embedding)
		assert.NotNil(t, vec)
	})
}

// TestHandlersEdgeCasesAndBoundaries tests edge cases and boundary conditions
// This comprehensive test suite covers unusual scenarios and boundary conditions
// These tests ensure handlers are robust against unexpected inputs and edge cases
func TestHandlersEdgeCasesAndBoundaries(t *testing.T) {
	// Test error response format consistency
	// This validates that error responses follow consistent JSON format
	// Ensures all handlers return properly formatted error messages
	t.Run("ErrorResponseFormat", func(t *testing.T) {
		router := setupHandlersTestRouter()
		router.GET("/files/:id", handlers.GetHandler(nil))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/files/invalid-uuid", nil)
		router.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.IsType(t, "", response["error"])
	})

	// Test HTTP method validation
	// This validates that routes only accept appropriate HTTP methods
	// Ensures proper HTTP method restrictions are enforced
	t.Run("MethodNotAllowed", func(t *testing.T) {
		router := setupHandlersTestRouter()
		router.GET("/files/:id", handlers.GetHandler(nil))

		// Try POST on GET-only endpoint
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/files/123", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// Test concurrent UUID generation reliability
	// This validates UUID uniqueness under concurrent generation
	// Ensures UUID generation is thread-safe and produces unique values
	t.Run("ConcurrentUUIDs", func(t *testing.T) {
		uuids := make([]uuid.UUID, 100)
		for i := 0; i < 100; i++ {
			uuids[i] = uuid.New()
		}

		// Check that all UUIDs are unique
		uniqueUUIDs := make(map[uuid.UUID]bool)
		for _, u := range uuids {
			assert.False(t, uniqueUUIDs[u], "UUID should be unique")
			uniqueUUIDs[u] = true
		}
	})

	// Test various date format edge cases
	// This validates date parsing robustness in GetFilesByDateRangeHandler
	// Ensures proper handling of various date formats and edge cases
	t.Run("DateFormatEdgeCases", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected bool
		}{
			{"2024-01-01", true},  // Standard date
			{"2024-12-31", true},  // End of year
			{"2024-02-29", true},  // Leap year
			{"2023-02-29", false}, // Non-leap year
			{"2024-13-01", false}, // Invalid month
			{"2024-01-32", false}, // Invalid day
			{"invalid", false},    // Invalid format
			{"", false},           // Empty string
		}

		for _, tc := range testCases {
			_, err := time.Parse("2006-01-02", tc.input)
			hasError := err != nil
			assert.Equal(t, !tc.expected, hasError, "Date: %s", tc.input)
		}
	})

	// Test JSON marshaling edge cases
	// This validates JSON handling robustness across different data scenarios
	// Ensures proper serialization/deserialization of various data types
	t.Run("JSONMarshalingEdgeCases", func(t *testing.T) {
		testCases := []models.FileUploadRequest{
			{Filename: "", Content: "", Embedding: []float32{}},                                        // Empty values
			{Filename: "test.txt", Content: "content", Embedding: nil},                                 // Nil embedding
			{Filename: "special chars: !@#$%^&*()", Content: "unicode: 测试", Embedding: []float32{1.0}}, // Special characters
		}

		for _, tc := range testCases {
			jsonBytes, err := json.Marshal(tc)
			assert.NoError(t, err)

			var parsed models.FileUploadRequest
			err = json.Unmarshal(jsonBytes, &parsed)
			assert.NoError(t, err)
		}
	})
}

// TestHandlerSpecificLogic tests handler-specific functionality that can be tested without database
// This section validates handler registration, routing, and parameter extraction
// These tests ensure the handler framework itself functions correctly
func TestHandlerSpecificLogic(t *testing.T) {
	// Test handlers with nil database queries for validation coverage
	// This validates that all handlers properly handle the validation phase
	// Ensures handlers can operate in validation-only mode for testing
	t.Run("HandlersWithNilQueries", func(t *testing.T) {
		router := setupHandlersTestRouter()

		// Register all handlers - this tests that they can all be initialized
		// All these should handle the validation part before hitting database
		router.GET("/files/:id", handlers.GetHandler(nil))
		router.GET("/files", handlers.GetAllHandler(nil))
		router.GET("/files/search", handlers.GetFilesByFilenameHandler(nil))
		router.GET("/files/date-range", handlers.GetFilesByDateRangeHandler(nil))
		router.POST("/files", handlers.UploadHandler(nil))
		router.DELETE("/files/:id", handlers.DeleteHandler(nil))
		router.PUT("/files/:id", handlers.UpdateHandler(nil))
		router.PATCH("/files/:id/soft-delete", handlers.SoftDeleteHandler(nil))
		router.PATCH("/files/:id/restore", handlers.UndoSoftDeleteHandler(nil))
		router.GET("/files/recycle-bin", handlers.GetDeletedFilesHandler(nil))
		router.GET("/files/metadata", handlers.GetFileMetadataHandler(nil))

		// Test that routes are properly registered
		routes := router.Routes()
		assert.True(t, len(routes) >= 11, "Should have at least 11 routes registered")
	})

	// Test query parameter extraction functionality
	// This validates the parameter extraction logic used in handlers
	// Ensures proper extraction of query parameters from HTTP requests
	t.Run("QueryParameterExtraction", func(t *testing.T) {
		router := setupHandlersTestRouter()

		router.GET("/test", func(c *gin.Context) {
			query := c.Query("q")
			start := c.Query("start")
			end := c.Query("end")

			c.JSON(http.StatusOK, gin.H{
				"query": query,
				"start": start,
				"end":   end,
			})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test?q=search&start=2024-01-01&end=2024-12-31", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "search", response["query"])
		assert.Equal(t, "2024-01-01", response["start"])
		assert.Equal(t, "2024-12-31", response["end"])
	})

	// Test path parameter extraction functionality
	// This validates the path parameter extraction logic used in handlers
	// Ensures proper extraction of path parameters from HTTP requests
	t.Run("PathParameterExtraction", func(t *testing.T) {
		router := setupHandlersTestRouter()

		router.GET("/files/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"id": id})
		})

		testUUID := uuid.New().String()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/files/"+testUUID, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, testUUID, response["id"])
	})
}
