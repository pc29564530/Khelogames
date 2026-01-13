package handlers

import (
	"fmt"
	"io"
	errorhandler "khelogames/error_handler"
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
		fieldErrors := make(map[string]string)
		if uploadId == "" {
			fieldErrors["uploadId"] = "Upload ID is required"
		}
		if chunkIndexStr == "" {
			fieldErrors["chunkIndex"] = "Chunk index is required"
		}
		if totalChunksStr == "" {
			fieldErrors["totalChunks"] = "Total chunks is required"
		}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	//convert the string to int
	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		s.logger.Error("Invalid chunkIndex: ", err)
		fieldErrors := map[string]string{"chunkIndex": "Invalid chunk index format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	// Read a chunk form the request body
	file, _, err := ctx.Request.FormFile("chunk")
	if err != nil {
		s.logger.Error("Failed to get file chunk: ", err)
		fieldErrors := map[string]string{"chunk": "Failed to get file chunk"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	defer file.Close()

	// Save to temp dir
	tempDir := filepath.Join("/tmp/khelogames_tmp_uploads", uploadId)
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		s.logger.Error("Failed to create upload dir: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create upload directory",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Save the chunk to the temp dir
	chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", chunkIndex))
	out, err := os.Create(chunkPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to save chunk",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	defer out.Close()

	// Copy the chunk from the request body to the file
	io.Copy(out, file)

	// Save the chunk to the database
	ctx.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data": gin.H{
			"message":     "Chunk uploaded",
			"chunk_index": chunkIndex,
		},
	})
}

func (s *HandlersServer) CompletedChunkUploadFunc(ctx *gin.Context) {
	// Get the req params
	var req struct {
		UploadID    string `json:"upload_id" binding:"required"`
		TotalChunks int    `json:"total_chunks" binding:"required,min=1"`
		MediaType   string `json:"media_type" binding:"required"`
	}

	// Parse the request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	// Get the upload from the database
	chunkDir := filepath.Join("/tmp/khelogames_tmp_uploads", req.UploadID)
	// Create upload directory
	finalDir := "/tmp/khelogames_media_uploads"
	if err := os.MkdirAll(finalDir, os.ModePerm); err != nil {
		s.logger.Error("Failed to create final upload dir: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create upload directory",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	var finalPath string
	mediaType := req.MediaType
	if mediaType == "image/jpeg" || mediaType == "image/png" || mediaType == "image/jpg" {
		finalPath = filepath.Join(finalDir, req.UploadID+".jpg")
	} else if mediaType == "video/mp4" || mediaType == "video/quicktime" || mediaType == "video/mkv" {
		finalPath = filepath.Join(finalDir, req.UploadID+".mp4")
	} else {
		fieldErrors := map[string]string{"media_type": "Unsupported media type"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	// Create the final file
	finalFile, err := os.Create(finalPath)
	if err != nil {
		s.logger.Error("Failed to create file path: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create final file",
			},
			"request_id": ctx.GetString("request_id"),
		})
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
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to open chunk",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}

		// Copy the chunk to the final file
		if _, err := io.Copy(finalFile, chunkFile); err != nil {
			s.logger.Error("Failed to copy chunk: ", err)
			chunkFile.Close()
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to copy chunk",
				},
				"request_id": ctx.GetString("request_id"),
			})
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
		"success": true,
		"data": gin.H{
			"message":   "Upload complete",
			"file_url":  fileURL,
			"upload_id": req.UploadID,
		},
	})
}
