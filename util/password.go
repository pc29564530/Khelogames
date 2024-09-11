package util

import (
	"khelogames/logger"

	"golang.org/x/crypto/bcrypt"
)

// Generate hashed password for particular user
func HashPassword(password string) (string, error) {
	logger := logger.NewLogger()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("failed to generate hashed password: ", err)
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
