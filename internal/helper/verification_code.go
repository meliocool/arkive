package helper

import (
	"crypto/rand"
	"fmt"
)

func GenerateVerificationCode() (string, error) {
	num := make([]byte, 3) // 3 bytes = up to ~16M
	if _, err := rand.Read(num); err != nil {
		return "", err
	}
	code := int(num[0])<<16 | int(num[1])<<8 | int(num[2])
	return fmt.Sprintf("%06d", code%1000000), nil
}
