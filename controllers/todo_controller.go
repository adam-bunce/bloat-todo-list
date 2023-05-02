package controllers

import (
	"context"
	"database/sql"
	"fmt"
	todo_service "github.com/adam-bunce/grpc-todo/domain/proto"
	"github.com/adam-bunce/grpc-todo/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TodoServiceServer struct {
	todo_service.UnimplementedTodoServiceServer
	UseCase usecase.TodoUsecase
}

// CreateToDo controller, use usecase in here
func (s *TodoServiceServer) CreateToDo(ctx context.Context, in *todo_service.CreateToDoMessage) (*todo_service.ToDo, error) {
	createdTodo, err := s.UseCase.CreateToDo(in.Todo)
	if err != nil {
		return &todo_service.ToDo{}, status.Error(codes.Internal, "Error creating ToDo")
	}

	return &todo_service.ToDo{
		Id:   createdTodo.Id,
		Todo: createdTodo.Todo,
	}, nil

}

func (s *TodoServiceServer) GetToDo(ctx context.Context, in *todo_service.GetToDoMessage) (*todo_service.ToDo, error) {
	gottenTodo, err := s.UseCase.GetToDo(in.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &todo_service.ToDo{}, status.Error(codes.NotFound, "No ToDos with that Id Found")
		}
		return &todo_service.ToDo{}, status.Error(codes.Internal, "Interal Error")
	}

	return &todo_service.ToDo{
		Id:   gottenTodo.Id,
		Todo: gottenTodo.Todo,
	}, nil
}

func (s *TodoServiceServer) GetAllToDos(ctx context.Context, in *todo_service.GetAllToDosMessage) (*todo_service.ToDos, error) {
	todos, err := s.UseCase.GetAllTodos()
	if err != nil {
		return &todo_service.ToDos{}, err
	}

	var responseMessage todo_service.ToDos

	for _, todo := range *todos {
		todoMsg := &todo_service.ToDo{
			Id:   todo.Id,
			Todo: todo.Todo,
		}
		responseMessage.Todos = append(responseMessage.Todos, todoMsg)
	}

	return &responseMessage, nil
}

func (s *TodoServiceServer) GetAllToDosStream(in *todo_service.GetAllToDosMessage, stream todo_service.TodoService_GetAllToDosStreamServer) error {
	// toy implementation, would only be useful if querying like 10k
	todos, err := s.UseCase.GetAllTodos()
	if err != nil {
		return err
	}

	for _, todo := range *todos {
		todoMsg := &todo_service.ToDo{
			Id:   todo.Id,
			Todo: todo.Todo,
		}
		stream.Send(todoMsg)
	}
	return nil
}

func (s *TodoServiceServer) UpdateToDo(ctx context.Context, in *todo_service.UpdateToDoMessage) (*todo_service.ToDo, error) {
	updatedTodo, err := s.UseCase.UpdateTodo(in.Id, in.UpdatedMessage)
	if err != nil {
		return &todo_service.ToDo{}, status.Error(codes.NotFound, "Error updating todo")
	}

	return &todo_service.ToDo{
		Id:   updatedTodo.Id,
		Todo: updatedTodo.Todo,
	}, nil
}

func (s *TodoServiceServer) DeleteToDo(ctx context.Context, in *todo_service.DeleteToDoMessage) (*todo_service.ToDo, error) {
	deletedTodo, err := s.UseCase.DeleteTodo(in.Id)
	if err != nil {
		return &todo_service.ToDo{}, status.Error(codes.Unknown, fmt.Sprintf("Error deleting ToDo %v", err))
	}

	return &todo_service.ToDo{
		Id:   deletedTodo.Id,
		Todo: deletedTodo.Todo,
	}, nil
}
