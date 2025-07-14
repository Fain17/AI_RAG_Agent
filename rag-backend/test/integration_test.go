package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fain17/rag-backend/api/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add basic middleware
	router.Use(gin.Recovery())

	return router
}

// TestBasicRouting tests that routes are properly defined
func TestBasicRouting(t *testing.T) {
	router := setupTestRouter()

	// Define a simple test route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

// TestJSONParsing tests JSON request parsing
func TestJSONParsing(t *testing.T) {
	router := setupTestRouter()

	router.POST("/test-json", func(c *gin.Context) {
		var req models.FileUploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"received": req.Filename})
	})

	// Test valid JSON
	validRequest := models.FileUploadRequest{
		Filename:  "test.txt",
		Content:   "test content",
		Embedding: []float32{1.0, 2.0, 3.0},
	}

	jsonBody, _ := json.Marshal(validRequest)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test-json", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "test.txt", response["received"])
}

func TestInvalidJSONParsing(t *testing.T) {
	router := setupTestRouter()

	router.POST("/test-json", func(c *gin.Context) {
		var req models.FileUploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"received": req.Filename})
	})

	// Test invalid JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test-json", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid json", response["error"])
}

// TestUUIDValidation tests UUID parameter validation
func TestUUIDValidation(t *testing.T) {
	router := setupTestRouter()

	router.GET("/test-uuid/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Basic UUID format check (simplified)
		if len(id) != 36 || id == "invalid-uuid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// Test valid UUID format
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-uuid/"+validUUID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Test invalid UUID
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test-uuid/invalid-uuid", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

// TestQueryParameters tests URL query parameter parsing
func TestQueryParameters(t *testing.T) {
	router := setupTestRouter()

	router.GET("/test-query", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing query parameter"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"query": query})
	})

	// Test with query parameter
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-query?q=test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "test", response["query"])

	// Test without query parameter
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test-query", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

// TestCORSHeaders tests that CORS headers can be set
func TestCORSHeaders(t *testing.T) {
	router := setupTestRouter()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	router.GET("/test-cors", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "cors test"})
	})

	// Test CORS headers
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-cors", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
}

// TestPanicRecovery tests that the recovery middleware works
func TestPanicRecovery(t *testing.T) {
	router := setupTestRouter()

	router.GET("/test-panic", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-panic", nil)
	router.ServeHTTP(w, req)

	// Should not crash the server, recovery middleware should handle it
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
