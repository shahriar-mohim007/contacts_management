package utilis

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Scope  string    `json:"scope"`
	jwt.RegisteredClaims
}

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

func GenerateJWT(userID uuid.UUID, scope string, secretKey string, ttl time.Duration) (string, error) {

	claims := Claims{
		UserID: userID,
		Scope:  scope,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
