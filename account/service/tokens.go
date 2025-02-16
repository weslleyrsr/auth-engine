package service

import (
	"crypto/rsa"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/weslleyrsr/auth-engine/account/model"
)

// IDTokenCustomClaims holds the structure of JWT claims for the ID token.
type IDTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.RegisteredClaims
}

// generateIDToken generates an ID token (JWT) with custom claims.
func generateIDToken(u *model.User, key *rsa.PrivateKey) (string, error) {
	now := time.Now()

	claims := IDTokenCustomClaims{
		User: u,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			// Optionally set other fields like Issuer, Subject, etc.
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		log.Println("Failed to sign ID token string:", err)
		return "", err
	}

	return ss, nil
}

// RefreshToken holds the signed JWT string along with its ID.
type RefreshToken struct {
	SS        string
	ID        string
	ExpiresIn time.Duration
}

// RefreshTokenCustomClaims holds the payload for a refresh token.
type RefreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.RegisteredClaims
}

// generateRefreshToken creates a refresh token that stores only the user's ID.
func generateRefreshToken(uid uuid.UUID, key string) (*RefreshToken, error) {
	now := time.Now()
	tokenExp := now.AddDate(0, 0, 3) // 3 days expiry

	tokenID, err := uuid.NewRandom()
	if err != nil {
		log.Println("Failed to generate refresh token ID:", err)
		return nil, err
	}

	claims := RefreshTokenCustomClaims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(tokenExp),
			ID:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		log.Println("Failed to sign refresh token string:", err)
		return nil, err
	}

	return &RefreshToken{
		SS:        ss,
		ID:        tokenID.String(),
		ExpiresIn: tokenExp.Sub(now),
	}, nil
}
