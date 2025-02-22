package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
	"github.com/weslleyrsr/auth-engine/account/model"
)

func TestNewPairFromUser(t *testing.T) {
	priv, err := os.ReadFile("../rsa_private_test.pem")
	if err != nil {
		t.Fatalf("failed to read private key file: %v", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		t.Fatalf("failed to parse private key: %v", err)
	}

	pub, err := os.ReadFile("../rsa_public_test.pem")
	if err != nil {
		t.Fatalf("failed to read public key file: %v", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		t.Fatalf("failed to parse public key: %v", err)
	}

	secret := "anotsorandomtestsecret"

	// instantiate a common token service to be used by all tests
	tokenService := NewTokenService(&TSConfig{
		PrivKey:       privKey,
		PubKey:        pubKey,
		RefreshSecret: secret,
	})

	if tokenService == nil {
		t.Fatal("tokenService is nil")
	}

	// include password to make sure it is not serialized
	// since json tag is "-"
	uid, _ := uuid.NewRandom()
	u := &model.User{
		UID:      uid,
		Email:    "bob@bob.com",
		Password: "blarghedymcblarghface",
	}

	t.Run("Returns a token pair with proper values", func(t *testing.T) {
		ctx := context.TODO()
		tokenPair, err := tokenService.NewPairFromUser(ctx, u, "")

		assert.NotEmpty(t, tokenPair.AccessToken, "AccessToken should not be empty")
		assert.NotEmpty(t, tokenPair.RefreshToken, "RefreshToken should not be empty")

		assert.NoError(t, err)

		var s string
		assert.IsType(t, s, tokenPair.AccessToken)

		// decode the Base64URL encoded string
		// simpler to use jwt library which is already imported
		idTokenClaims := &IDTokenCustomClaims{}

		_, err = jwt.ParseWithClaims(tokenPair.AccessToken, idTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		})

		assert.NoError(t, err)

		// assert claims on idToken
		expectedClaims := []interface{}{
			u.UID,
			u.Email,
			u.Name,
			u.ImageURL,
			u.Website,
		}
		actualIDClaims := []interface{}{
			idTokenClaims.User.UID,
			idTokenClaims.User.Email,
			idTokenClaims.User.Name,
			idTokenClaims.User.ImageURL,
			idTokenClaims.User.Website,
		}

		assert.ElementsMatch(t, expectedClaims, actualIDClaims)
		assert.Empty(t, idTokenClaims.User.Password) // password should never be encoded to json

		expiresAt := time.Unix(idTokenClaims.ExpiresAt.Unix(), 0)
		expectedExpiresAt := time.Now().Add(15 * time.Minute)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		refreshTokenClaims := &RefreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		assert.IsType(t, s, tokenPair.RefreshToken)

		// assert claims on refresh token
		assert.NoError(t, err)
		assert.Equal(t, u.UID, refreshTokenClaims.UID)

		expiresAt = time.Unix(refreshTokenClaims.ExpiresAt.Unix(), 0)
		expectedExpiresAt = time.Now().Add(3 * 24 * time.Hour)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)
	})
}
