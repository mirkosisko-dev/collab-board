package utils

import (
	"fmt"
	"math/rand"
)

var domains = []string{"example.com", "test.com", "yourdomain.com"}

// `GenerateRandomEmail` generates a random email address.
func GenerateRandomEmail(mailLength int) string {
	username := GenerateRandomString(mailLength)
	domain := domains[rand.Intn(len(domains))]
	email := fmt.Sprintf("%s@%s", username, domain)
	return email
}

func GenerateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func GenerateRandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
