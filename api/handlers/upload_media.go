package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *HandlersServer) CreateUploadMediaFunc(ctx *gin.Context) {
	//getting the data from form
	uploadId := ctx.PostForm("uploadId")
	chunkIndexStr := ctx.PostForm("chunkIndex")
	totalChunksStr := ctx.PostForm("totalChunks")

	//checking if the data is valid
	if uploadId == "" || chunkIndexStr == "" || totalChunksStr == "" {
		s.logger.Error("Missing uploadId, chunkIndex, or totalChunks")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing uploadId, chunkIndex, or totalChunks"})
		return
	}

	//convert the string to int
	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		s.logger.Error("Invalid chunkIndex: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chunkIndex"})
		return
	}

	// Read a chunk form the request body
	file, _, err := ctx.Request.FormFile("chunk")
	if err != nil {
		s.logger.Error("Failed to get file chunk: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file chunk"})
		return
	}

	defer file.Close()

	// Save to temp dir
	tempDir := filepath.Join("/tmp/khelogames_tmp_uploads", uploadId)
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		s.logger.Error("Failed to create upload dir: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload dir"})
		return
	}

	// Save the chunk to the temp dir
	chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", chunkIndex))
	out, err := os.Create(chunkPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chunk"})
		return
	}
	defer out.Close()

	// Copy the chunk from the request body to the file
	io.Copy(out, file)

	// Save the chunk to the database
	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Chunk uploaded",
		"chunk_index": chunkIndex,
	})
}

func (s *HandlersServer) CompletedChunkUploadFunc(ctx *gin.Context) {
	// Get the req params
	var req struct {
		UploadID    string `json:"upload_id"`
		TotalChunks int    `json:"total_chunks"`
		MediaType   string `json:"media_type"`
	}

	// Parse the request body
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate required fields
	if req.UploadID == "" || req.TotalChunks <= 0 || req.MediaType == "" {
		s.logger.Error("Missing required fields")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Get the upload from the database
	chunkDir := filepath.Join("/tmp/khelogames_tmp_uploads", req.UploadID)
	// Create upload directory
	finalDir := "/tmp/khelogames_media_uploads"
	if err := os.MkdirAll(finalDir, os.ModePerm); err != nil {
		s.logger.Error("Failed to create final upload dir: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	var finalPath string
	mediaType := req.MediaType
	if mediaType == "image/jpeg" || mediaType == "image/png" || mediaType == "image/jpg" {
		finalPath = filepath.Join(finalDir, req.UploadID+".jpg")
	} else if mediaType == "video/mp4" || mediaType == "video/quicktime" || mediaType == "video/mkv" {
		finalPath = filepath.Join(finalDir, req.UploadID+".mp4")
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported media type"})
		return
	}

	// Create the final file
	finalFile, err := os.Create(finalPath)
	if err != nil {
		s.logger.Error("Failed to create file path: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create final file"})
		return
	}
	defer finalFile.Close()

	// Copy the chunks from the temporary directory to the final file
	for i := 0; i < req.TotalChunks; i++ {
		// Get the chunk from the database
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("chunk_%d", i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			s.logger.Error("Failed to open chunk: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open chunk"})
			return
		}

		// Copy the chunk to the final file
		if _, err := io.Copy(finalFile, chunkFile); err != nil {
			s.logger.Error("Failed to copy chunk: ", err)
			chunkFile.Close()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy chunk"})
			return
		}

		chunkFile.Close()
	}

	// Remove the temp chunks
	if err := os.RemoveAll(chunkDir); err != nil {
		s.logger.Error("Failed to remove temp chunks: ", err)
		// Continue even if cleanup fails, as the main operation succeeded
	}

	var fileExt string
	if mediaType == "image/jpeg" || mediaType == "image/png" || mediaType == "image/jpg" {
		fileExt = "jpg"
	} else if mediaType == "video/mp4" || mediaType == "video/quicktime" || mediaType == "video/mkv" {
		fileExt = "mp4"
	}

	// Return the final file path - use relative path or configurable base URL
	// For now, use a relative path that can be served by a static file server
	fileURL := fmt.Sprintf("http://192.168.1.3:8080/media/%s.%s", req.UploadID, fileExt)

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Upload complete",
		"file_url":  fileURL,
		"upload_id": req.UploadID,
	})
}
