package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"gocourse/config"
	"gocourse/internal/database"
	"time"
)

type RefreshTokenDTO struct {
	Id        int             `json:"id"`
	UserID    int             `json:"user_id"`
	Token     string          `json:"token"`
	ExpiresAt jwt.NumericDate `json:"expires_at"`
}

type JwtManager struct {
	pg               *database.DbPool
	secret           string
	AccessExpiresAt  time.Duration
	RefreshExpiresAt time.Duration
}

func NewJwtManager(cfg *config.Config, pg *database.DbPool) *JwtManager {
	return &JwtManager{
		pg:               pg,
		secret:           cfg.Secret,
		AccessExpiresAt:  cfg.AccessExpiresAt,
		RefreshExpiresAt: cfg.RefreshExpiresAt,
	}
}

func (m *JwtManager) GenerateJWT(userId string, ttl time.Duration) (string, error) {
	expirationTime := time.Now().Add(ttl)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   userId,
	})

	return token.SignedString([]byte(m.secret))
}

func (m *JwtManager) ValidateJWT(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
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

	return claims, nil
}

func (m *JwtManager) SaveRefreshToken(refreshToken string) error {
	claims, err := m.ValidateJWT(refreshToken)
	if err != nil {
		return fmt.Errorf("error get claims from token in saveRefreshToken: %v", err)
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return fmt.Errorf("error get user_id from token in saveRefreshToken: %v", err)
	}

	expiresAt, ok := claims["exp"].(jwt.NumericDate)
	if !ok {
		return fmt.Errorf("error get expires_at from token in saveRefreshToken: %v", err)
	}

	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (@user_id, @token, @expires_at)`
	args := pgx.NamedArgs{
		"user_id":    userID,
		"token":      refreshToken,
		"expires_at": expiresAt,
	}

	_, err = m.pg.Db.Exec(m.pg.Ctx, query, args)
	if err != nil {
		return fmt.Errorf("error saving refresh token: %v", err)
	}

	return nil
}

func (m *JwtManager) GetRefreshToken(userID int) (RefreshTokenDTO, error) {
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

func (m *JwtManager) DeleteRefreshToken(userID int) error {
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
