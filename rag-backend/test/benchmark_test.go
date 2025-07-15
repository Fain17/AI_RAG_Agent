package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fain17/rag-backend/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Set Gin to test mode for benchmarks
	gin.SetMode(gin.TestMode)
}

// BenchmarkUUIDValidation benchmarks UUID validation performance
func BenchmarkUUIDValidation(b *testing.B) {
	validUUID := uuid.New().String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := uuid.Parse(validUUID)
		if err != nil {
			b.Fatal("UUID parsing failed")
		}
	}
}

// BenchmarkJSONMarshaling benchmarks JSON marshaling of FileUploadRequest
func BenchmarkJSONMarshaling(b *testing.B) {
	req := models.FileUploadRequest{
		Filename:  "benchmark.txt",
		Content:   "This is benchmark content that should be reasonably sized",
		Embedding: make([]float32, 512), // Common embedding size
	}

	// Fill embedding with sample data
	for i := range req.Embedding {
		req.Embedding[i] = float32(i) * 0.1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(req)
		if err != nil {
			b.Fatal("JSON marshaling failed")
		}
	}
}

// BenchmarkJSONUnmarshaling benchmarks JSON unmarshaling of FileUploadRequest
func BenchmarkJSONUnmarshaling(b *testing.B) {
	req := models.FileUploadRequest{
		Filename:  "benchmark.txt",
		Content:   "This is benchmark content that should be reasonably sized",
		Embedding: make([]float32, 512),
	}

	// Fill embedding with sample data
	for i := range req.Embedding {
		req.Embedding[i] = float32(i) * 0.1
	}

	jsonData, _ := json.Marshal(req)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var decoded models.FileUploadRequest
		err := json.Unmarshal(jsonData, &decoded)
		if err != nil {
			b.Fatal("JSON unmarshaling failed")
		}
	}
}

// BenchmarkRouterPerformance benchmarks basic router performance
func BenchmarkRouterPerformance(b *testing.B) {
	router := gin.New()
	router.GET("/benchmark", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "benchmark"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/benchmark", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatal("Router request failed")
		}
	}
}

// BenchmarkFileUploadEndpoint benchmarks the file upload endpoint simulation
func BenchmarkFileUploadEndpoint(b *testing.B) {
	router := gin.New()
	router.POST("/upload", func(c *gin.Context) {
		var req models.FileUploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "uploaded"})
	})

	uploadReq := models.FileUploadRequest{
		Filename:  "benchmark.txt",
		Content:   "Benchmark content",
		Embedding: make([]float32, 128),
	}

	jsonData, _ := json.Marshal(uploadReq)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/upload", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatal("Upload request failed")
		}
	}
}

// BenchmarkLargeEmbeddingProcessing benchmarks processing of large embeddings
func BenchmarkLargeEmbeddingProcessing(b *testing.B) {
	// Test with different embedding sizes
	sizes := []int{128, 512, 1024, 4096}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("EmbeddingSize-%d", size), func(b *testing.B) {
			req := models.FileUploadRequest{
				Filename:  "large-embedding.txt",
				Content:   "Content with large embedding",
				Embedding: make([]float32, size),
			}

			// Fill with sample data
			for i := range req.Embedding {
				req.Embedding[i] = float32(i) * 0.001
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				jsonData, err := json.Marshal(req)
				if err != nil {
					b.Fatal("Failed to marshal large embedding")
				}

				var decoded models.FileUploadRequest
				err = json.Unmarshal(jsonData, &decoded)
				if err != nil {
					b.Fatal("Failed to unmarshal large embedding")
				}
			}
		})
	}
}

// BenchmarkConcurrentRequests benchmarks concurrent request handling
func BenchmarkConcurrentRequests(b *testing.B) {
	router := gin.New()
	router.GET("/concurrent", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"id": c.Query("id")})
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/concurrent?id=test", nil)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				b.Fatal("Concurrent request failed")
			}
		}
	})
}

// BenchmarkMemoryAllocation tests memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("SmallRequest", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			req := models.FileUploadRequest{
				Filename:  "small.txt",
				Content:   "Small content",
				Embedding: []float32{1.0, 2.0, 3.0},
			}

			jsonData, _ := json.Marshal(req)
			var decoded models.FileUploadRequest
			json.Unmarshal(jsonData, &decoded)
		}
	})

	b.Run("LargeRequest", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			req := models.FileUploadRequest{
				Filename:  "large.txt",
				Content:   string(make([]byte, 10000)), // 10KB content
				Embedding: make([]float32, 1000),       // Large embedding
			}

			jsonData, _ := json.Marshal(req)
			var decoded models.FileUploadRequest
			json.Unmarshal(jsonData, &decoded)
		}
	})
}

// Performance test that can be run as a regular test
func TestPerformanceMetrics(t *testing.T) {
	// Test that demonstrates performance characteristics
	router := gin.New()
	router.POST("/perf", func(c *gin.Context) {
		var req models.FileUploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"received": len(req.Embedding)})
	})

	// Test with various payload sizes
	testSizes := []int{10, 100, 1000}

	for _, size := range testSizes {
		t.Run(fmt.Sprintf("PayloadSize-%d", size), func(t *testing.T) {
			req := models.FileUploadRequest{
				Filename:  "perf-test.txt",
				Content:   "Performance test content",
				Embedding: make([]float32, size),
			}

			jsonData, err := json.Marshal(req)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			httpReq, _ := http.NewRequest("POST", "/perf", bytes.NewBuffer(jsonData))
			httpReq.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, httpReq)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, float64(size), response["received"])
		})
	}
}
