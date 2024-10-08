package secure

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	AccessExpiration time.Duration

	publicKey  *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
}

func NewJWTManager(config *common.JWTConfig) (*JWTManager, error) {
	privateKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(config.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	publicKey, err := jwt.ParseECPublicKeyFromPEM([]byte(config.PublicKey))
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	return &JWTManager{
		privateKey:       privateKey,
		publicKey:        publicKey,
		AccessExpiration: config.AccessExpiration,
	}, nil
}

type JWTClaims struct {
	UserUUID string
	jwt.RegisteredClaims
}

// GetPublicKey returns the public key used for token verification
func (jm *JWTManager) GetPublicKey() *ecdsa.PublicKey {
	return jm.publicKey
}

// GenerateAccessToken returns signed jwt token string
func (jm *JWTManager) GenerateAccessToken(userUUID string) (string, error) {
	if err := uuid.Validate(userUUID); err != nil {
		return "", errors.New("invalid user uuid")
	}

	claims := JWTClaims{
		UserUUID: userUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.AccessExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signedString, err := token.SignedString(jm.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to get signed string from the token: %w", err)
	}

	return signedString, nil
}

func ValidateToken(tokenString string, jwtPublicKey *ecdsa.PublicKey) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token string cannot be empty")
	}

	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return jwtPublicKey, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, keyFunc)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
