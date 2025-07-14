package test

import (
	"encoding/json"
	"testing"

	"github.com/fain17/rag-backend/api/models"
	"github.com/stretchr/testify/assert"
)

// Test FileUploadRequest model
func TestFileUploadRequestModel(t *testing.T) {
	t.Run("ValidFileUploadRequest", func(t *testing.T) {
		req := models.FileUploadRequest{
			Filename:  "test.txt",
			Content:   "test content",
			Embedding: []float32{1.0, 2.0, 3.0},
		}

		assert.Equal(t, "test.txt", req.Filename)
		assert.Equal(t, "test content", req.Content)
		assert.Equal(t, []float32{1.0, 2.0, 3.0}, req.Embedding)
		assert.Len(t, req.Embedding, 3)
	})

	t.Run("EmptyFileUploadRequest", func(t *testing.T) {
		req := models.FileUploadRequest{}

		assert.Equal(t, "", req.Filename)
		assert.Equal(t, "", req.Content)
		assert.Nil(t, req.Embedding)
	})

	t.Run("FileUploadRequestWithLargeEmbedding", func(t *testing.T) {
		largeEmbedding := make([]float32, 1000)
		for i := range largeEmbedding {
			largeEmbedding[i] = float32(i)
		}

		req := models.FileUploadRequest{
			Filename:  "large.txt",
			Content:   "large content",
			Embedding: largeEmbedding,
		}

		assert.Equal(t, "large.txt", req.Filename)
		assert.Equal(t, "large content", req.Content)
		assert.Len(t, req.Embedding, 1000)
		assert.Equal(t, float32(999), req.Embedding[999])
	})

	t.Run("FileUploadRequestJSONMarshaling", func(t *testing.T) {
		original := models.FileUploadRequest{
			Filename:  "test.json",
			Content:   "json content",
			Embedding: []float32{1.5, 2.5, 3.5},
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(original)
		assert.NoError(t, err)
		assert.NotEmpty(t, jsonData)

		// Unmarshal back
		var unmarshaled models.FileUploadRequest
		err = json.Unmarshal(jsonData, &unmarshaled)
		assert.NoError(t, err)

		assert.Equal(t, original.Filename, unmarshaled.Filename)
		assert.Equal(t, original.Content, unmarshaled.Content)
		assert.Equal(t, original.Embedding, unmarshaled.Embedding)
	})

	t.Run("FileUploadRequestJSONUnmarshalingEdgeCases", func(t *testing.T) {
		testCases := []struct {
			name     string
			jsonStr  string
			hasError bool
		}{
			{
				name:     "ValidJSON",
				jsonStr:  `{"filename":"test.txt","content":"content","embedding":[1.0,2.0]}`,
				hasError: false,
			},
			{
				name:     "JSONWithNullEmbedding",
				jsonStr:  `{"filename":"test.txt","content":"content","embedding":null}`,
				hasError: false,
			},
			{
				name:     "JSONWithEmptyEmbedding",
				jsonStr:  `{"filename":"test.txt","content":"content","embedding":[]}`,
				hasError: false,
			},
			{
				name:     "InvalidJSON",
				jsonStr:  `{"filename":"test.txt","content":"content","embedding":}`,
				hasError: true,
			},
			{
				name:     "EmptyJSON",
				jsonStr:  `{}`,
				hasError: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var req models.FileUploadRequest
				err := json.Unmarshal([]byte(tc.jsonStr), &req)

				if tc.hasError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("FileUploadRequestWithSpecialCharacters", func(t *testing.T) {
		req := models.FileUploadRequest{
			Filename:  "файл.txt", // Cyrillic
			Content:   "内容测试",     // Chinese
			Embedding: []float32{-1.0, 0.0, 1.0},
		}

		// Should handle unicode correctly
		jsonData, err := json.Marshal(req)
		assert.NoError(t, err)

		var unmarshaled models.FileUploadRequest
		err = json.Unmarshal(jsonData, &unmarshaled)
		assert.NoError(t, err)

		assert.Equal(t, req.Filename, unmarshaled.Filename)
		assert.Equal(t, req.Content, unmarshaled.Content)
		assert.Equal(t, req.Embedding, unmarshaled.Embedding)
	})

	t.Run("FileUploadRequestWithNegativeEmbeddings", func(t *testing.T) {
		req := models.FileUploadRequest{
			Filename:  "negative.txt",
			Content:   "negative embeddings",
			Embedding: []float32{-1.0, -2.5, -3.14, 0.0, 1.0, 2.5, 3.14},
		}

		assert.Contains(t, req.Embedding, float32(-1.0))
		assert.Contains(t, req.Embedding, float32(-2.5))
		assert.Contains(t, req.Embedding, float32(-3.14))
		assert.Contains(t, req.Embedding, float32(0.0))
		assert.Contains(t, req.Embedding, float32(1.0))
	})

	t.Run("FileUploadRequestWithVeryLongContent", func(t *testing.T) {
		longContent := make([]byte, 10000)
		for i := range longContent {
			longContent[i] = byte('a' + (i % 26))
		}

		req := models.FileUploadRequest{
			Filename:  "long.txt",
			Content:   string(longContent),
			Embedding: []float32{1.0},
		}

		assert.Equal(t, "long.txt", req.Filename)
		assert.Equal(t, 10000, len(req.Content))
		assert.Equal(t, byte('a'), req.Content[0])
		assert.Equal(t, byte('z'), req.Content[25])
	})

	t.Run("FileUploadRequestFieldTypes", func(t *testing.T) {
		req := models.FileUploadRequest{
			Filename:  "type-test.txt",
			Content:   "content",
			Embedding: []float32{1.0, 2.0},
		}

		// Test field types
		assert.IsType(t, "", req.Filename)
		assert.IsType(t, "", req.Content)
		assert.IsType(t, []float32{}, req.Embedding)
	})

	t.Run("FileUploadRequestCopy", func(t *testing.T) {
		original := models.FileUploadRequest{
			Filename:  "original.txt",
			Content:   "original content",
			Embedding: []float32{1.0, 2.0, 3.0},
		}

		// Create a copy
		copy := models.FileUploadRequest{
			Filename:  original.Filename,
			Content:   original.Content,
			Embedding: append([]float32(nil), original.Embedding...),
		}

		assert.Equal(t, original.Filename, copy.Filename)
		assert.Equal(t, original.Content, copy.Content)
		assert.Equal(t, original.Embedding, copy.Embedding)

		// Modify original shouldn't affect copy
		original.Filename = "modified.txt"
		original.Embedding[0] = 999.0

		assert.NotEqual(t, original.Filename, copy.Filename)
		assert.NotEqual(t, original.Embedding[0], copy.Embedding[0])
	})
}
