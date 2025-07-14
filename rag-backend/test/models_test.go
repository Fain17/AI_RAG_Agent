package test

import (
	"encoding/json"
	"testing"

	"github.com/fain17/rag-backend/api/models"
	"github.com/stretchr/testify/assert"
)

// TestFileUploadRequest tests the FileUploadRequest model
func TestFileUploadRequest(t *testing.T) {
	// Test valid file upload request
	req := models.FileUploadRequest{
		Filename:  "test.txt",
		Content:   "This is test content",
		Embedding: []float32{1.0, 2.0, 3.0, 4.0, 5.0},
	}

	assert.Equal(t, "test.txt", req.Filename)
	assert.Equal(t, "This is test content", req.Content)
	assert.Equal(t, 5, len(req.Embedding))
	assert.Equal(t, float32(1.0), req.Embedding[0])
}

// TestFileUploadRequestJSON tests JSON marshaling/unmarshaling
func TestFileUploadRequestJSON(t *testing.T) {
	original := models.FileUploadRequest{
		Filename:  "document.pdf",
		Content:   "PDF content here",
		Embedding: []float32{0.1, 0.2, 0.3},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Unmarshal from JSON
	var decoded models.FileUploadRequest
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)

	// Verify values
	assert.Equal(t, original.Filename, decoded.Filename)
	assert.Equal(t, original.Content, decoded.Content)
	assert.Equal(t, len(original.Embedding), len(decoded.Embedding))

	for i, val := range original.Embedding {
		assert.Equal(t, val, decoded.Embedding[i])
	}
}

// TestFileUploadRequestValidation tests different validation scenarios
func TestFileUploadRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     models.FileUploadRequest
		expectValid bool
	}{
		{
			name: "Valid request",
			request: models.FileUploadRequest{
				Filename:  "valid.txt",
				Content:   "Valid content",
				Embedding: []float32{1.0, 2.0},
			},
			expectValid: true,
		},
		{
			name: "Empty filename",
			request: models.FileUploadRequest{
				Filename:  "",
				Content:   "Content without filename",
				Embedding: []float32{1.0, 2.0},
			},
			expectValid: true, // Empty filename might be allowed depending on business logic
		},
		{
			name: "Empty content",
			request: models.FileUploadRequest{
				Filename:  "empty-content.txt",
				Content:   "",
				Embedding: []float32{1.0, 2.0},
			},
			expectValid: true, // Empty content might be allowed
		},
		{
			name: "Empty embedding",
			request: models.FileUploadRequest{
				Filename:  "no-embedding.txt",
				Content:   "Content without embedding",
				Embedding: []float32{},
			},
			expectValid: true, // Empty embedding might be allowed
		},
		{
			name: "Nil embedding",
			request: models.FileUploadRequest{
				Filename:  "nil-embedding.txt",
				Content:   "Content with nil embedding",
				Embedding: nil,
			},
			expectValid: true, // Nil embedding should be handled gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling doesn't fail
			jsonData, err := json.Marshal(tt.request)
			assert.NoError(t, err)
			assert.NotEmpty(t, jsonData)

			// Test JSON unmarshaling doesn't fail
			var decoded models.FileUploadRequest
			err = json.Unmarshal(jsonData, &decoded)
			assert.NoError(t, err)
		})
	}
}

// TestFileUploadRequestEdgeCases tests edge cases
func TestFileUploadRequestEdgeCases(t *testing.T) {
	// Test very long filename
	longFilename := make([]byte, 1000)
	for i := range longFilename {
		longFilename[i] = 'a'
	}

	req := models.FileUploadRequest{
		Filename:  string(longFilename),
		Content:   "Content",
		Embedding: []float32{1.0},
	}

	jsonData, err := json.Marshal(req)
	assert.NoError(t, err)

	var decoded models.FileUploadRequest
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, string(longFilename), decoded.Filename)
}

// TestFileUploadRequestLargeEmbedding tests handling of large embeddings
func TestFileUploadRequestLargeEmbedding(t *testing.T) {
	// Create a large embedding array
	largeEmbedding := make([]float32, 1000)
	for i := range largeEmbedding {
		largeEmbedding[i] = float32(i) * 0.1
	}

	req := models.FileUploadRequest{
		Filename:  "large-embedding.txt",
		Content:   "Content with large embedding",
		Embedding: largeEmbedding,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test JSON unmarshaling
	var decoded models.FileUploadRequest
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, len(largeEmbedding), len(decoded.Embedding))
	assert.Equal(t, largeEmbedding[0], decoded.Embedding[0])
	assert.Equal(t, largeEmbedding[999], decoded.Embedding[999])
}

// TestFileUploadRequestSpecialCharacters tests handling of special characters
func TestFileUploadRequestSpecialCharacters(t *testing.T) {
	specialChars := "file-name_with.special@chars#123!.txt"
	unicodeContent := "Content with unicode: 你好世界 🌍 émojis"

	req := models.FileUploadRequest{
		Filename:  specialChars,
		Content:   unicodeContent,
		Embedding: []float32{-1.0, 0.0, 1.0},
	}

	// Test JSON marshaling preserves special characters
	jsonData, err := json.Marshal(req)
	assert.NoError(t, err)

	var decoded models.FileUploadRequest
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, specialChars, decoded.Filename)
	assert.Equal(t, unicodeContent, decoded.Content)
}

// TestFileUploadRequestNumericValues tests various numeric values in embeddings
func TestFileUploadRequestNumericValues(t *testing.T) {
	req := models.FileUploadRequest{
		Filename: "numeric-test.txt",
		Content:  "Testing numeric values",
		Embedding: []float32{
			0.0,          // zero
			1.0,          // positive
			-1.0,         // negative
			0.123456789,  // decimal
			-0.987654321, // negative decimal
			1e-10,        // very small
			1e10,         // very large
		},
	}

	jsonData, err := json.Marshal(req)
	assert.NoError(t, err)

	var decoded models.FileUploadRequest
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, len(req.Embedding), len(decoded.Embedding))

	// Check each value with appropriate tolerance for floating point comparison
	for i, expected := range req.Embedding {
		assert.InDelta(t, expected, decoded.Embedding[i], 1e-6)
	}
}
