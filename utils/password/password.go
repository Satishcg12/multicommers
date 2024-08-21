package password

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given password using bcrypt with an increased cost factor.
func HashPassword(password string) (string, error) {
	const cost = 14 // Increase the cost factor for better security
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

// ComparePasswords compares a hashed password with a plain password using constant time comparison.
func ComparePasswords(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
