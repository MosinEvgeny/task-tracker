package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/google/uuid"
)

// UserRepository реализует интерфейс UserRepository для работы с пользователями в PostgreSQL.
type UserRepository struct {
	db *PostgresDB
}

// NewUserRepository создает новый экземпляр UserRepository.
func NewUserRepository(db *PostgresDB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, username, email, password)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.DB.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("ошибка в создании пользователя: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, username, email, password
		FROM users
		WHERE id = $1
	`

	row := r.db.DB.QueryRowContext(ctx, query, id)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь не найден: %w", err)
		}
		return nil, fmt.Errorf("ошибка в получении пользователя по ID: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password
		FROM users
		WHERE email = $1
	`

	row := r.db.DB.QueryRowContext(ctx, query, email)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь не найден: %w", err)
		}
		return nil, fmt.Errorf("ошибка в получении пользователя по email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET username = $2, email = $3, password = $4
		WHERE id = $1
	`

	_, err := r.db.DB.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении пользователя: %w", err)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}

	return nil
}
