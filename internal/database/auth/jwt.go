package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"go-rest-api-auth/config"
	"go-rest-api-auth/internal/database"
	"time"
)

type RefreshTokenDTO struct {
	Id        int             `json:"id"`
	UserID    int             `json:"user_id"`
	Token     string          `json:"token"`
	ExpiresAt jwt.NumericDate `json:"expires_at"`
}

type JwtManagerImplementation struct {
	pg               *database.DbPool
	secret           string
	AccessExpiresAt  time.Duration
	RefreshExpiresAt time.Duration
}

//go:generate go run github.com/vektra/mockery/v2@v2.46.3 --name JwtManager --output ../../../testing/mocks
type JwtManager interface {
	GenerateJWT(userId string, tokenType string, ttl time.Duration) (string, error)
	ValidateJWT(reqToken string, expectedType string) (jwt.MapClaims, error)
	SaveRefreshToken(refreshToken string) error
	GetRefreshToken(userID int) (RefreshTokenDTO, error)
	IsRefreshTokenValid(refreshToken string) (bool, error)
	DeleteRefreshToken(userID int) error
	GetterAccessExpiresAt() time.Duration
	GetterRefreshExpiresAt() time.Duration
}

func NewJwtManager(cfg *config.Config, pg *database.DbPool) JwtManager {
	return &JwtManagerImplementation{
		pg:               pg,
		secret:           cfg.Secret,
		AccessExpiresAt:  cfg.AccessExpiresAt,
		RefreshExpiresAt: cfg.RefreshExpiresAt,
	}
}

type CustomClaims struct {
	jwt.RegisteredClaims
	TokenType string `json:"token_type"`
}

func (m *JwtManagerImplementation) GetterAccessExpiresAt() time.Duration {
	return m.AccessExpiresAt
}

func (m *JwtManagerImplementation) GetterRefreshExpiresAt() time.Duration {
	return m.RefreshExpiresAt
}

func (m *JwtManagerImplementation) GenerateJWT(userId string, tokenType string, ttl time.Duration) (string, error) {
	expirationTime := time.Now().Add(ttl)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   userId,
		},
		TokenType: tokenType,
	})

	return token.SignedString([]byte(m.secret))
}

func (m *JwtManagerImplementation) ValidateJWT(reqToken string, expectedType string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.secret), nil
	})

	if err != nil {
		return jwt.MapClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("error get user claims from token")
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != expectedType {
		return jwt.MapClaims{}, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

func (m *JwtManagerImplementation) SaveRefreshToken(refreshToken string) error {
	claims, err := m.ValidateJWT(refreshToken, "refresh")

	if err != nil {
		return fmt.Errorf("error get claims from token in saveRefreshToken: %v", err)
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return fmt.Errorf("error get user_id from token in saveRefreshToken: %v", err)
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("error get expires_at from token in saveRefreshToken: %v", err)
	}

	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (@user_id, @token, @expires_at)`
	args := pgx.NamedArgs{
		"user_id":    userID,
		"token":      refreshToken,
		"expires_at": time.Unix(int64(exp), 0),
	}

	_, err = m.pg.Db.Exec(m.pg.Ctx, query, args)
	if err != nil {
		return fmt.Errorf("error saving refresh token: %v", err)
	}

	return nil
}

func (m *JwtManagerImplementation) GetRefreshToken(userID int) (RefreshTokenDTO, error) {
	query := `SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE user_id = @user_id`
	args := pgx.NamedArgs{
		"user_id": userID,
	}

	row := m.pg.Db.QueryRow(m.pg.Ctx, query, args)
	token := RefreshTokenDTO{}
	err := row.Scan(&token.Id, &token.UserID, &token.Token, &token.ExpiresAt)
	if err != nil {
		return RefreshTokenDTO{}, fmt.Errorf("error getting refresh token: %v", err)
	}

	return token, nil
}

func (m *JwtManagerImplementation) IsRefreshTokenValid(refreshToken string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM refresh_tokens WHERE token = $1"
	err := m.pg.Db.QueryRow(m.pg.Ctx, query, refreshToken).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *JwtManagerImplementation) DeleteRefreshToken(userID int) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = @user_id`
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	_, err := m.pg.Db.Exec(m.pg.Ctx, query, args)
	if err != nil {
		return fmt.Errorf("error deleting refresh token: %v", err)
	}
	return nil
}
