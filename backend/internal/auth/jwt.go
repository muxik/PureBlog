package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Token type discriminators carried in the "typ" claim.
const (
	TokenAccess  = "access"
	TokenRefresh = "refresh"
)

// ErrInvalidToken is returned for malformed, expired, or wrongly-signed tokens.
var ErrInvalidToken = errors.New("invalid token")

// Claims is the JWT payload for both access and refresh tokens.
type Claims struct {
	jwt.RegisteredClaims
	Type     string `json:"typ"`
	Username string `json:"username,omitempty"`
}

// UserID parses the subject claim back into a numeric user id.
func (c *Claims) UserID() int64 {
	id, _ := strconv.ParseInt(c.Subject, 10, 64)
	return id
}

// TokenManager issues and verifies HS256 JWTs.
type TokenManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewTokenManager builds a TokenManager.
func NewTokenManager(secret string, accessTTL, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (m *TokenManager) issue(userID int64, username, typ string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
		Type:     typ,
		Username: username,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

// Access issues a short-lived access token.
func (m *TokenManager) Access(userID int64, username string) (string, error) {
	return m.issue(userID, username, TokenAccess, m.accessTTL)
}

// Refresh issues a long-lived refresh token.
func (m *TokenManager) Refresh(userID int64, username string) (string, error) {
	return m.issue(userID, username, TokenRefresh, m.refreshTTL)
}

// Parse verifies a token string and returns its claims.
func (m *TokenManager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
