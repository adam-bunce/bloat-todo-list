package repository

import (
	"database/sql"
	"github.com/adam-bunce/grpc-todo/domain"
	"github.com/adam-bunce/grpc-todo/variables"
	"log"
)

type todoRepository struct {
	database *sql.DB
}

type TodoRepository interface {
	CreateToDo(string) (*domain.Todo, error)
	GetToDoMessage(int32) (*domain.Todo, error)
	GetAllToDos() (*[]domain.Todo, error)
	UpdateToDo(int32, string) (*domain.Todo, error)
	DeleteTodo(int32) (*domain.Todo, error)
}

func NewTaskRepository(db *sql.DB) TodoRepository {
	return &todoRepository{database: db}
}

func (t todoRepository) CreateToDo(todo string) (*domain.Todo, error) {
	var createdTodo domain.Todo

	// try to create row
	result, err := variables.DB.Query("INSERT INTO newTable (todo) VALUES ($1) returning id, todo", todo)
	if err != nil {
		return &domain.Todo{}, err
	}
	defer result.Close()

	// read row back in from returned row
	result.Next()
	err = result.Scan(&createdTodo.Id, &createdTodo.Todo)
	if err != nil {
		return &domain.Todo{}, err
	}

	return &createdTodo, nil
}

func (t todoRepository) GetToDoMessage(id int32) (*domain.Todo, error) {
	var foundTodo domain.Todo

	err := variables.DB.QueryRow("SELECT id, todo FROM newTable WHERE id = $1", id).Scan(&foundTodo.Id, &foundTodo.Todo)
	if err != nil {
		return &domain.Todo{}, err
	}

	if err != nil {
		return &domain.Todo{}, err
	}

	return &foundTodo, nil
}

func (t todoRepository) GetAllToDos() (*[]domain.Todo, error) {
	var foundToDos []domain.Todo

	rows, err := variables.DB.Query("SELECT id, todo From newTable")
	if err != nil {
		return &foundToDos, err
	}

	for rows.Next() != false {
		var row domain.Todo
		err := rows.Scan(&row.Id, &row.Todo)
		if err != nil {
			log.Println("Failed to read row from GetAllToDos")
		}

		foundToDos = append(foundToDos, row)
	}

	return &foundToDos, nil
}

func (t todoRepository) UpdateToDo(id int32, updatedToDo string) (*domain.Todo, error) {
	var todo domain.Todo
	err := t.database.QueryRow("UPDATE newTable SET todo = $1 WHERE id = $2 RETURNING id, todo", updatedToDo, id).Scan(&todo.Id, &todo.Todo)
	if err != nil {
		return &domain.Todo{}, err
	}
	return &todo, nil
}

func (t todoRepository) DeleteTodo(id int32) (*domain.Todo, error) {
	var todo domain.Todo
	err := t.database.QueryRow("DELETE FROM newTable where id = $1 returning id, todo", id).Scan(&todo.Id, &todo.Todo)
	if err != nil {
		return &domain.Todo{}, err
	}
	return &todo, nil
}
