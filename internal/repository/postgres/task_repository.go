package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/google/uuid"
)

// TaskRepository реализует интерфейс TaskRepository для работы с задачами в PostgreSQL.
type TaskRepository struct {
	db *PostgresDB
}

// NewTaskRepository создает новый экземпляр TaskRepository.
func NewTaskRepository(db *PostgresDB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, due_date, user_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.DB.ExecContext(ctx, query, task.ID, task.Title, task.Description, task.DueDate, task.UserID)
	if err != nil {
		return fmt.Errorf("ошибка при создании задачи: %w", err)
	}

	// TODO: Обработка label_ids (связь задачи и меток)

	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	query := `
		SELECT id, title, description, due_date, user_id
		FROM tasks
		WHERE id = $1
	`

	row := r.db.DB.QueryRowContext(ctx, query, id)

	var task domain.Task
	var dueDate time.Time
	if err := row.Scan(&task.ID, &task.Title, &task.Description, &dueDate, &task.UserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена: %w", err)
		}
		return nil, fmt.Errorf("ошибка при получении задачи по ID: %w", err)
	}

	task.DueDate = dueDate

	return &task, nil
}

func (r *TaskRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Task, error) {
	query := `
		SELECT id, title, description, due_date, user_id
		FROM tasks
		WHERE user_id = $1
	`

	rows, err := r.db.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении задач пользователя: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var task domain.Task
		var dueDate time.Time
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &dueDate, &task.UserID); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании задачи: %w", err)
		}
		task.DueDate = dueDate
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по задачам: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	query := `
		UPDATE tasks
		SET title = $2, description = $3, due_date = $4
		WHERE id = $1
	`

	_, err := r.db.DB.ExecContext(ctx, query, task.ID, task.Title, task.Description, task.DueDate)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении задачи: %w", err)
	}

	// TODO: Обработка label_ids (обновление связи задачи и меток)

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM tasks
		WHERE id = $1
	`

	_, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении задачи: %w", err)
	}

	return nil
}
