package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/theghostmac/todo-api-with-gin/internal/models"

	_ "github.com/lib/pq"
)

// TodoRepository defines the operations available for a Todo repository.
type TodoRepository interface {
	GetAll(ctx context.Context) ([]models.Todo, error)
	GetByID(ctx context.Context, id int) (*models.Todo, error)
	Create(ctx context.Context, todo *models.Todo) error
	Update(ctx context.Context, id int, todo *models.Todo) error
	Delete(ctx context.Context, id int) error
}

type PostgresTodoRepository struct {
	DB *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &PostgresTodoRepository{DB: db}
}

func (r *PostgresTodoRepository) GetAll(ctx context.Context) ([]models.Todo, error) {
    todos := []models.Todo{}
    query := "SELECT id, title, completed FROM todos"

    rows, err := r.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var todo models.Todo
        if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed); err != nil {
            return nil, err
        }
        todos = append(todos, todo)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return todos, nil
}

func (r *PostgresTodoRepository) GetByID(ctx context.Context, id int) (*models.Todo, error) {
    todo := &models.Todo{}
    query := "SELECT id, title, completed FROM todos WHERE id = $1"

    row := r.DB.QueryRowContext(ctx, query, id)
    if err := row.Scan(&todo.ID, &todo.Title, &todo.Completed); err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("no todo found with id: %d", id)
        }
        return nil, err
    }

    return todo, nil
}

func (r *PostgresTodoRepository) Create(ctx context.Context, todo *models.Todo) error {
    query := "INSERT INTO todos (title, completed) VALUES ($1, $2) RETURNING id"
    return r.DB.QueryRowContext(ctx, query, todo.Title, todo.Completed).Scan(&todo.ID)
}

func (r *PostgresTodoRepository) Update(ctx context.Context, id int, todo *models.Todo) error {
    query := "UPDATE todos SET title = $1, completed = $2 WHERE id = $3"
    _, err := r.DB.ExecContext(ctx, query, todo.Title, todo.Completed, id)
    return err
}

func (r *PostgresTodoRepository) Delete(ctx context.Context, id int) error {
    query := "DELETE FROM todos WHERE id = $1"
    _, err := r.DB.ExecContext(ctx, query, id)
    return err
}
