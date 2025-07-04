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
	uploadId := ctx.PostForm("uploadId")
	chunkIndexStr := ctx.PostForm("chunkIndex")
	totalChunksStr := ctx.PostForm("totalChunks")
	fmt.Println("Line no 1838i ")
	fmt.Println("Upload ID: ", uploadId)
	fmt.Println("Chunk Index: ", chunkIndexStr)
	fmt.Println("Total Chunks; ", totalChunksStr)

	if uploadId == "" || chunkIndexStr == "" || totalChunksStr == "" {
		s.logger.Error("Missing uploadId, chunkIndex, or totalChunks")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing uploadId, chunkIndex, or totalChunks"})
		return
	}

	fmt.Println("Line no 28")

	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		s.logger.Error("Invalid chunkIndex: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chunkIndex"})
		return
	}

	// Read chunk
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

	chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", chunkIndex))
	out, err := os.Create(chunkPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chunk"})
		return
	}
	defer out.Close()
	io.Copy(out, file)

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Chunk uploaded",
		"chunk_index": chunkIndex,
	})
}

func (s *HandlersServer) CompletedChunkUploadFunc(ctx *gin.Context) {
	var req struct {
		UploadID    string `json:"upload_id"`
		TotalChunks int    `json:"total_chunks"`
		MediaType   string `json:"media_type"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
	}

	// uploadId := ctx.Query("uploadId")
	// mediaType := ctx.Query("mediaType")
	// fmt.Println("Upload Id: ", uploadId)
	// fmt.Println("Media Type: ", mediaType)
	// totalChunksStr := ctx.Query("totalChunks")
	// totalChunks, err := strconv.Atoi(totalChunksStr)
	// if err != nil {
	// 	s.logger.Error("Failed to convert to int: ", err)
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chunk"})
	// 	return
	// }

	fmt.Println("Lin eno media type: ", req.MediaType)
	chunkDir := filepath.Join("/tmp/khelogames_tmp_uploads", req.UploadID)
	finalDir := "/tmp/khelogames_media_uploads"
	os.MkdirAll(finalDir, os.ModePerm)

	var finalPath string
	if req.MediaType == "image" {
		finalPath = filepath.Join(finalDir, req.UploadID+".jpg")
	} else if req.MediaType == "video" {
		finalPath = filepath.Join(finalDir, req.UploadID+".mp4")
	}
	finalFile, err := os.Create(finalPath)
	if err != nil {
		s.logger.Error("Failed to create file path: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create final file"})
	}

	defer finalFile.Close()

	for i := 0; i < req.TotalChunks; i++ {
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("chunk_%d", i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			s.logger.Error("Failed to open chunks: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open chunk"})
		}

		io.Copy(finalFile, chunkFile)

		chunkFile.Close()
	}

	_ = os.RemoveAll(chunkDir)

	var fileExt string
	if req.MediaType == "image" {
		fileExt = "jpg"
	} else if req.MediaType == "video" {
		fileExt = "mp4"
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported media type"})
		return
	}

	fileURL := fmt.Sprintf("http://192.168.1.3:8080/media/%s.%s", req.UploadID, fileExt)

	fmt.Println("Url: ", fileURL)
	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Upload complete",
		"file_url":  fileURL,
		"upload_id": req.UploadID,
	})

}
