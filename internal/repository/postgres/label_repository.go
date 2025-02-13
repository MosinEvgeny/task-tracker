package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/google/uuid"
)

// LabelRepository реализует интерфейс LabelRepository для работы с метками в PostgreSQL.
type LabelRepository struct {
	db *PostgresDB
}

// NewLabelRepository создает новый экземпляр LabelRepository.
func NewLabelRepository(db *PostgresDB) *LabelRepository {
	return &LabelRepository{db: db}
}

func (r *LabelRepository) Create(ctx context.Context, label *domain.Label) error {
	query := `
		INSERT INTO labels (id, name, color, user_id)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.DB.ExecContext(ctx, query, label.ID, label.Name, label.Color, label.UserID)
	if err != nil {
		return fmt.Errorf("ошибка при создании метки: %w", err)
	}

	return nil
}

func (r *LabelRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Label, error) {
	query := `
		SELECT id, name, color, user_id
		FROM labels
		WHERE id = $1
	`

	row := r.db.DB.QueryRowContext(ctx, query, id)

	var label domain.Label
	if err := row.Scan(&label.ID, &label.Name, &label.Color, &label.UserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("метка не найдена: %w", err)
		}
		return nil, fmt.Errorf("ошибка при получении метки по ID: %w", err)
	}

	return &label, nil
}

func (r *LabelRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Label, error) {
	query := `
		SELECT id, name, color, user_id
		FROM labels
		WHERE user_id = $1
	`

	rows, err := r.db.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении меток пользователя: %w", err)
	}
	defer rows.Close()

	var labels []*domain.Label
	for rows.Next() {
		var label domain.Label
		if err := rows.Scan(&label.ID, &label.Name, &label.Color, &label.UserID); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании метки: %w", err)
		}
		labels = append(labels, &label)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по меткам: %w", err)
	}

	return labels, nil
}

func (r *LabelRepository) Update(ctx context.Context, label *domain.Label) error {
	query := `
		UPDATE labels
		SET name = $2, color = $3
		WHERE id = $1
	`

	_, err := r.db.DB.ExecContext(ctx, query, label.ID, label.Name, label.Color)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении метки: %w", err)
	}

	return nil
}

func (r *LabelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM labels
		WHERE id = $1
	`

	_, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении метки: %w", err)
	}

	return nil
}
