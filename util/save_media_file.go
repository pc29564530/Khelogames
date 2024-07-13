package util

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"khelogames/logger"
	"math/rand"
	"os"
	"strings"
)

type SaveImageStruct struct {
	logger *logger.Logger
}

func NewSaveImageStruct(logger *logger.Logger) *SaveImageStruct {
	return &SaveImageStruct{logger: logger}
}

func (s *SaveImageStruct) SaveImageToFile(data []byte, mediaType string) (string, error) {
	s.logger.Info("Starting SaveImageToFile function")

	randomString, err := s.generateRandomString(12)
	if err != nil {
		s.logger.Error("Error generating random string: %v", err)
		return "", err
	}
	s.logger.Debug("Generated random string: %s", randomString)

	var mediaFolder string
	switch mediaType {
	case "image":
		mediaFolder = "images"
	case "video":
		mediaFolder = "videos"
	default:
		err := fmt.Errorf("unsupported media type: %s", mediaType)
		s.logger.Error("Unsupported media type: %v", err)
		return "", err
	}

	filePath := fmt.Sprintf("/Users/pawan/database/Khelogames/%s/%s", mediaFolder, randomString)
	file, err := os.Create(filePath)
	if err != nil {
		s.logger.Error("Failed to create file: %v", err)
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(data))
	if err != nil {
		s.logger.Error("Failed to copy data to file: %v", err)
		return "", err
	}

	s.logger.Debug("File created successfully at path: %s", filePath)

	path := s.convertLocalPathToURL(filePath, mediaFolder)
	s.logger.Info("Image saved successfully, URL: %s", path)

	return path, nil
}

func (s *SaveImageStruct) generateRandomString(length int) (string, error) {
	s.logger.Info("Starting generateRandomString function")

	if length%2 != 0 {
		err := fmt.Errorf("length must be even for generating hex string")
		s.logger.Error("Invalid length for random string generation: %v", err)
		return "", err
	}

	randomBytes := make([]byte, length/2)
	_, err := rand.Read(randomBytes)
	if err != nil {
		s.logger.Error("Failed to read random bytes: %v", err)
		return "", err
	}

	randomString := hex.EncodeToString(randomBytes)
	s.logger.Debug("Generated random hex string: %s", randomString)
	return randomString, nil
}

func (s *SaveImageStruct) convertLocalPathToURL(localPath string, mediaFolder string) string {
	s.logger.Info("Starting convertLocalPathToURL function")

	baseURL := fmt.Sprintf("http://10.0.2.2:8080/%s/", mediaFolder)
	imagePath := baseURL + strings.TrimPrefix(localPath, fmt.Sprintf("/Users/pawan/database/Khelogames/%s/", mediaFolder))
	s.logger.Debug("Converted local path to URL: %s", imagePath)

	return imagePath
}
