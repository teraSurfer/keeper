package database

import (
	"context"
	"log"
)

type Todo struct {
	Id          int
	Title       string
	Description string
	Completed   bool
}

type TodoModel interface {
	ListTodos(context.Context) (*[]Todo, error)
	GetTodo(context.Context, int) (*Todo, error)
	UpdateTodo(context.Context, int, *Todo) error
	DeleteTodo(context.Context, int) error
	Save(context.Context, *Todo) error
}

func (s *service) ListTodos(ctx context.Context) (*[]Todo, error) {
	rows, err := s.db.QueryContext(ctx, "select id, title, description, completed from todos")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos = make([]Todo, 0)

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed); err != nil {
			return &todos, err
		}

		todos = append(todos, todo)
	}

	return &todos, nil
}

func (s *service) GetTodo(ctx context.Context, id int) (*Todo, error) {
	row := s.db.QueryRowContext(ctx, "select id, title, description, completed from todos where id=?", id)

	var todo Todo

	if err := row.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed); err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s *service) UpdateTodo(ctx context.Context, id int, todo *Todo) error {
	result, err := s.db.ExecContext(ctx,
		`
			UPDATE todos
			SET title=?,
			description=?,
			completed=?
			WHERE id=?
		`, todo.Title, todo.Description, todo.Completed, id,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	log.Println("updated", rows)

	todo.Id = id

	return nil
}

func (s *service) DeleteTodo(ctx context.Context, id int) error {
	result, err := s.db.ExecContext(ctx,
		`
			DELETE FROM todos where id=?
		`, id,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	log.Println("deleted", rows)

	return nil
}

func (s *service) Save(ctx context.Context, todo *Todo) error {
	result, err := s.db.ExecContext(ctx,
		`
			INSERT INTO todos (title, description, completed)
			VALUES (?, ?, ?)
		`,
		todo.Title, todo.Description, todo.Completed,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	todo.Id = int(id)

	return nil
}
