package usecase

import (
	"github.com/adam-bunce/grpc-todo/domain"
	"github.com/adam-bunce/grpc-todo/repository"
)

type todoUsecase struct {
	todoRepository repository.TodoRepository
}

type TodoUsecase interface {
	CreateToDo(string) (*domain.Todo, error)
	GetToDo(int32) (*domain.Todo, error)
	GetAllTodos() (*[]domain.Todo, error)
	UpdateTodo(int32, string) (*domain.Todo, error)
	DeleteTodo(int32) (*domain.Todo, error)
}

func NewTodoUseCase(todoRepository repository.TodoRepository) TodoUsecase {
	return &todoUsecase{todoRepository: todoRepository}
}

func (t todoUsecase) CreateToDo(todo string) (*domain.Todo, error) {
	return t.todoRepository.CreateToDo(todo)
}

func (t todoUsecase) GetToDo(id int32) (*domain.Todo, error) {
	return t.todoRepository.GetToDoMessage(id)
}

func (t todoUsecase) GetAllTodos() (*[]domain.Todo, error) {

	resp, err := t.todoRepository.GetAllToDos()
	return resp, err
}

func (t todoUsecase) UpdateTodo(id int32, todo string) (*domain.Todo, error) {
	return t.todoRepository.UpdateToDo(id, todo)
}

func (t todoUsecase) DeleteTodo(id int32) (*domain.Todo, error) {
	return t.todoRepository.DeleteTodo(id)
}
