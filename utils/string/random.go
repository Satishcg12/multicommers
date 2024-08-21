package randomString

import "math/rand"

func GenerateRandomString(length int) string {
	// Define the character set
	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a byte slice of the specified length
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[RandomInt(len(charset))]
	}

	return string(randomString)
}

func RandomInt(max int) int {
	return rand.Intn(max)
}
