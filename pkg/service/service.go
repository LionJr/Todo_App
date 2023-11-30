package service

import (
	maksimzhashkevychtodoapp "github.com/LionJr/todo-app"
	"github.com/LionJr/todo-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user maksimzhashkevychtodoapp.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list maksimzhashkevychtodoapp.TodoList) (int, error)
	GetAll(userId int) ([]maksimzhashkevychtodoapp.TodoList, error)
	GetById(userId, listId int) (maksimzhashkevychtodoapp.TodoList, error)
	Delete(userId, listIs int) error
	Update(userId, listId int, input maksimzhashkevychtodoapp.UpdateListInput) error
}

type TodoItem interface {
	Create(userId, listId int, item maksimzhashkevychtodoapp.TodoItem) (int, error)
	GetAll(userId, listId int) ([]maksimzhashkevychtodoapp.TodoItem, error)
	GetById(userId, itemId int) (maksimzhashkevychtodoapp.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, listId int, input maksimzhashkevychtodoapp.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
