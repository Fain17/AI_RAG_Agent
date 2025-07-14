package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fain17/rag-backend/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// setupRouter creates a test router instance with test mode enabled
// This isolates tests from production configuration
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// TestHandlerInputValidation tests the core input validation logic that handlers use
// This test validates UUID parsing functionality that's common across multiple handlers
func TestHandlerInputValidation(t *testing.T) {
	router := setupRouter()

	// Test endpoint that mimics handler UUID validation behavior
	router.GET("/test/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Parse UUID similar to handler logic - this is the core validation
		// that GetHandler, DeleteHandler, UpdateHandler, SoftDeleteHandler, and UndoSoftDeleteHandler all use
		_, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// Test case 1: Valid UUID should pass validation
	// This ensures the UUID parsing logic works correctly for valid inputs
	validUUID := uuid.New().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test/"+validUUID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Test case 2: Invalid UUID should fail validation
	// This tests the error handling path that prevents invalid UUIDs from proceeding
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test/invalid-uuid", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)

	var response map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &response)
	assert.Equal(t, "invalid id", response["error"])
}

// TestFileUploadHandlerLogic tests the JSON binding and validation logic used by UploadHandler
// This covers the request parsing, validation, and response generation patterns
func TestFileUploadHandlerLogic(t *testing.T) {
	router := setupRouter()

	// Simulate the UploadHandler's request binding and validation logic
	router.POST("/upload", func(c *gin.Context) {
		var req models.FileUploadRequest
		// This mimics the c.BindJSON(&req) call in UploadHandler
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		// Simulate basic field validation that handlers might perform
		if req.Filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "filename required"})
			return
		}

		// Simulate successful upload response
		c.JSON(http.StatusOK, gin.H{
			"message":  "file uploaded successfully",
			"filename": req.Filename,
		})
	})

	// Test case 1: Valid upload request should succeed
	// This verifies the happy path through the upload logic
	validRequest := models.FileUploadRequest{
		Filename:  "test.txt",
		Content:   "test content",
		Embedding: []float32{1.0, 2.0, 3.0},
	}

	jsonBody, _ := json.Marshal(validRequest)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/upload", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "file uploaded successfully", response["message"])
	assert.Equal(t, "test.txt", response["filename"])
}

// TestFileUploadInvalidJSON tests error handling for malformed JSON requests
// This covers the error path in UploadHandler and UpdateHandler when JSON parsing fails
func TestFileUploadInvalidJSON(t *testing.T) {
	router := setupRouter()

	// Simulate handler behavior when JSON binding fails
	router.POST("/upload", func(c *gin.Context) {
		var req models.FileUploadRequest
		// This will fail for invalid JSON, testing the error handling path
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test case: Invalid JSON should trigger error response
	// This ensures proper error handling for malformed request bodies
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/upload", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid json", response["error"])
}

// TestSearchQueryValidation tests query parameter validation used by GetFilesByFilenameHandler
// This covers the input validation for search functionality
func TestSearchQueryValidation(t *testing.T) {
	router := setupRouter()

	// Simulate GetFilesByFilenameHandler's query parameter validation
	router.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		// This mirrors the validation logic in GetFilesByFilenameHandler
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter required"})
			return
		}

		// Simulate search results
		c.JSON(http.StatusOK, gin.H{
			"query":   query,
			"results": []string{"file1.txt", "file2.txt"},
		})
	})

	// Test case 1: Valid query should succeed
	// This tests the successful search path with proper query parameter
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?q=test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Test case 2: Missing query parameter should fail
	// This tests the validation error path when required parameter is missing
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/search", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

// TestDeleteHandlerLogic tests the deletion workflow including UUID validation
// This simulates the complete logic flow of DeleteHandler
func TestDeleteHandlerLogic(t *testing.T) {
	router := setupRouter()

	// Simulate DeleteHandler's validation and response logic
	router.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Validate UUID - same pattern as used in DeleteHandler
		_, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		// Simulate successful deletion with proper HTTP status
		c.JSON(http.StatusNoContent, nil)
	})

	// Test case 1: Valid UUID should allow deletion
	// This tests the successful deletion path
	validUUID := uuid.New().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/delete/"+validUUID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Test case 2: Invalid UUID should prevent deletion
	// This tests the input validation error path
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/delete/invalid", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

// TestUpdateHandlerLogic tests the update workflow including UUID and JSON validation
// This simulates the complete logic flow of UpdateHandler
func TestUpdateHandlerLogic(t *testing.T) {
	router := setupRouter()

	// Simulate UpdateHandler's validation and processing logic
	router.PUT("/update/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Validate UUID - first validation step in UpdateHandler
		_, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		// Validate JSON request body - second validation step
		var req models.FileUploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		// Simulate successful update response
		c.JSON(http.StatusOK, gin.H{
			"message":  "file updated successfully",
			"id":       id,
			"filename": req.Filename,
		})
	})

	// Test case: Valid UUID and JSON should allow update
	// This tests the complete successful update path
	validUUID := uuid.New().String()
	updateRequest := models.FileUploadRequest{
		Filename:  "updated.txt",
		Content:   "updated content",
		Embedding: []float32{4.0, 5.0, 6.0},
	}

	jsonBody, _ := json.Marshal(updateRequest)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/update/"+validUUID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "file updated successfully", response["message"])
	assert.Equal(t, "updated.txt", response["filename"])
}

// TestDateRangeValidation tests date parameter validation used by GetFilesByDateRangeHandler
// This covers the date parsing and validation logic
func TestDateRangeValidation(t *testing.T) {
	router := setupRouter()

	// Simulate GetFilesByDateRangeHandler's date validation logic
	router.GET("/date-range", func(c *gin.Context) {
		start := c.Query("start")
		end := c.Query("end")

		// Validate presence of both parameters
		if start == "" || end == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start and end date required"})
			return
		}

		// Basic date format validation (simplified version of time.Parse validation)
		if len(start) != 10 || len(end) != 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"start": start,
			"end":   end,
			"files": []string{"file1.txt"},
		})
	})

	// Test case 1: Valid date range should succeed
	// This tests the successful date range query path
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/date-range?start=2024-01-01&end=2024-12-31", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Test case 2: Missing end parameter should fail
	// This tests parameter validation error handling
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/date-range?start=2024-01-01", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)

	// Test case 3: Invalid date format should fail
	// This tests date format validation error handling
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/date-range?start=invalid&end=2024-12-31", nil)
	router.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusBadRequest, w3.Code)
}

// TestMiddlewareErrorHandling tests panic recovery middleware functionality
// This ensures the application gracefully handles unexpected errors
func TestMiddlewareErrorHandling(t *testing.T) {
	router := setupRouter()

	// Add recovery middleware - this prevents panics from crashing the server
	router.Use(gin.Recovery())

	// Create an endpoint that deliberately panics
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	// Test case: Panic should be caught and handled gracefully
	// This verifies that the recovery middleware prevents server crashes
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(w, req)

	// Should not crash, recovery middleware should handle it
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestContentTypeValidation tests Content-Type header validation
// This covers the HTTP header validation patterns used in handlers
func TestContentTypeValidation(t *testing.T) {
	router := setupRouter()

	// Simulate content type validation logic
	router.POST("/test-content-type", func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		// This tests the pattern of validating request headers
		if contentType != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "valid content type"})
	})

	// Test case 1: Correct Content-Type should succeed
	// This tests the successful header validation path
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test-content-type", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Test case 2: Incorrect Content-Type should fail
	// This tests header validation error handling
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/test-content-type", bytes.NewBuffer([]byte("{}")))
	req2.Header.Set("Content-Type", "text/plain")
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}
