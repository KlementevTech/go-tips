package todos

import (
	"context"

	"github.com/KlementevTech/gotips/api/gen/go/todos/v1"
)

var _ todosv1.TodosServiceServer = (*Service)(nil)

type Service struct {
	todosv1.UnimplementedTodosServiceServer
}

func NewService() *Service {
	return &Service{}
}

func (s Service) AddTodo(_ context.Context, _ *todosv1.AddTodoRequest) (*todosv1.AddTodoResponse, error) {
	return &todosv1.AddTodoResponse{}, nil
}

func (s Service) ListTodos(_ context.Context, _ *todosv1.ListTodosRequest) (*todosv1.ListTodosResponse, error) {
	return &todosv1.ListTodosResponse{}, nil
}

func (s Service) CompleteTodo(
	_ context.Context,
	_ *todosv1.CompleteTodoRequest,
) (*todosv1.CompleteTodoResponse, error) {
	return &todosv1.CompleteTodoResponse{}, nil
}

func (s Service) DeleteTodo(
	_ context.Context,
	_ *todosv1.DeleteTodoRequest,
) (*todosv1.DeleteTodoResponse, error) {
	return &todosv1.DeleteTodoResponse{}, nil
}
