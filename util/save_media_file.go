package util

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

func SaveImageToFile(data []byte, mediaType string) (string, error) {
	randomString, err := generateRandomString(12)
	if err != nil {
		fmt.Printf("Error generating random string: %v\n", err)
		return "", err
	}

	var mediaFolder string
	switch mediaType {
	case "image":
		mediaFolder = "images"
	case "video":
		mediaFolder = "videos"
	default:
		return "", fmt.Errorf("unsupported media type for inserting in mediaFolder: %s", mediaType)
	}

	filePath := fmt.Sprintf("/Users/pawan/database/Khelogames/%s/%s", mediaFolder, randomString)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	path := convertLocalPathToURL(filePath, mediaFolder)
	return path, nil
}

func generateRandomString(length int) (string, error) {
	if length%2 != 0 {
		return "", fmt.Errorf("length must be even for generating hex string")
	}

	randomBytes := make([]byte, length/2)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
}

func convertLocalPathToURL(localPath string, mediaFolder string) string {
	baseURL := fmt.Sprintf("http://10.0.2.2:8080/%s/", mediaFolder)
	imagePath := baseURL + strings.TrimPrefix(localPath, fmt.Sprintf("/Users/pawan/database/Khelogames/%s/", mediaFolder))
	filePath := imagePath
	return filePath
}
