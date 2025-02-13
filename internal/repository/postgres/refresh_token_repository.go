package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/google/uuid"
)

// RefreshTokenRepository реализует интерфейс RefreshTokenRepository для работы с refresh токенами в PostgreSQL.
type RefreshTokenRepository struct {
	db *PostgresDB
}

// NewRefreshTokenRepository создает новый экземпляр RefreshTokenRepository.
func NewRefreshTokenRepository(db *PostgresDB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create создает новый refresh токен в базе данных.
func (r *RefreshTokenRepository) Create(ctx context.Context, refreshToken *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expiry_date)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.DB.ExecContext(ctx, query, refreshToken.ID, refreshToken.UserID, refreshToken.Token, refreshToken.ExpiryDate)
	if err != nil {
		return fmt.Errorf("ошибка при создании refresh токена: %w", err)
	}

	return nil
}

// GetByToken возвращает refresh токен по токену из базы данных.
func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expiry_date
		FROM refresh_tokens
		WHERE token = $1
	`

	row := r.db.DB.QueryRowContext(ctx, query, token)

	var refreshToken domain.RefreshToken
	if err := row.Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiryDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh токен не найден: %w", err)
		}
		return nil, fmt.Errorf("ошибка при получении refresh токена по токену: %w", err)
	}

	return &refreshToken, nil
}

// Delete удаляет refresh токен из базы данных.
func (r *RefreshTokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE id = $1
	`

	_, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении refresh токена: %w", err)
	}

	return nil
}

// DeleteAllByUserID удаляет все refresh токены пользователя из базы данных.
func (r *RefreshTokenRepository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1
	`

	_, err := r.db.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении всех refresh токенов пользователя: %w", err)
	}

	return nil
}
