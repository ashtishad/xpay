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
	PublicKey        *ecdsa.PublicKey
	AccessExpiration time.Duration

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
		PublicKey:        publicKey,
		AccessExpiration: config.AccessExpiration,
	}, nil
}

type JWTClaims struct {
	UserUUID string
	jwt.RegisteredClaims
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
