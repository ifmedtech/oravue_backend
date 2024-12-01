package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"oravue_backend/internal/config"
	"time"
)

func GenerateJWT(userID string, config *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,                                // Add user-specific claims
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Set expiration time (24 hours)
		"iat":     time.Now().Unix(),                     // Issued at time
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(config.Jwt.Secret)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// VerifyJWT verifies the JWT token and returns the claims if valid
func VerifyJWT(tokenString string, config *config.Config) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return config.Jwt.Secret, nil
	})

	// Check if parsing the token resulted in an error
	if err != nil {
		return nil, err
	}

	// Verify the token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
