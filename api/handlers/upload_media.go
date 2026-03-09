package handlers

import (
	"fmt"
	"io"
	errorhandler "khelogames/error_handler"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	if _, err := io.Copy(out, file); err != nil {
		s.logger.Error("Failed to write chunk: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to write chunk",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Save the chunk to the database
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message":     "Chunk uploaded",
			"chunk_index": chunkIndex,
		},
	})
}

func (s *HandlersServer) CompletedChunkUploadFunc(ctx *gin.Context) {

	// Request body
	var req struct {
		UploadID    string `json:"upload_id" binding:"required"`
		TotalChunks int    `json:"total_chunks" binding:"required,min=1"`
		MediaType   string `json:"media_type" binding:"required"`
	}

	// Parse JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	// Temp chunk directory
	chunkDir := filepath.Join("/tmp/khelogames_tmp_uploads", req.UploadID)

	// Final local file directory
	finalDir := "/tmp/khelogames_media_uploads"

	if err := os.MkdirAll(finalDir, os.ModePerm); err != nil {
		s.logger.Error("Failed to create final upload dir: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create upload directory",
		})
		return
	}

	// Determine extension and folder
	var fileExt string
	var folder string

	switch req.MediaType {
	case "image/jpeg", "image/png", "image/jpg":
		fileExt = "jpg"
		folder = "images"
	case "video/mp4", "video/quicktime", "video/mkv":
		fileExt = "mp4"
		folder = "videos"
	default:
		errorhandler.ValidationErrorResponse(ctx, map[string]string{
			"media_type": "Unsupported media type",
		})
		return
	}

	// Local merged file path
	finalPath := filepath.Join(finalDir, req.UploadID+"."+fileExt)

	// Create merged file
	finalFile, err := os.Create(finalPath)
	if err != nil {
		s.logger.Error("Failed to create final file: ", err)
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

	// Merge chunks
	for i := 0; i < req.TotalChunks; i++ {

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

		if _, err := io.Copy(finalFile, chunkFile); err != nil {
			chunkFile.Close()
			s.logger.Error("Failed to merge chunk: ", err)
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

	// Remove chunk directory
	os.RemoveAll(chunkDir)

	// Create R2 object key
	fileKey := fmt.Sprintf("media/%s/%s.%s", folder, req.UploadID, fileExt)

	// Open merged file for upload
	file, err := os.Open(finalPath)
	if err != nil {
		s.logger.Error("Failed to open final file: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to read merged file",
		})
		return
	}
	defer file.Close()

	// Upload to R2
	_, err = s.r2Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.config.R2BucketName),
		Key:         aws.String(fileKey),
		Body:        file,
		ContentType: aws.String(req.MediaType),
	})

	if err != nil {
		s.logger.Error("R2 upload failed:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to upload media",
		})
		return
	}

	// Remove local merged file
	os.Remove(finalPath)

	mediaUrl := fmt.Sprintf("%s/%s", s.config.R2BasePublicUrl, fileKey)

	// Return response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message":   "Upload complete",
			"media_url": mediaUrl,
		},
	})
}
