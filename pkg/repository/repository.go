package repository

import (
	maksimzhashkevychtodoapp "github.com/LionJr/todo-app"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user maksimzhashkevychtodoapp.User) (int, error)
	GetUser(username, password string) (maksimzhashkevychtodoapp.User, error)
}

type TodoList interface {
	Create(userId int, list maksimzhashkevychtodoapp.TodoList) (int, error)
	GetAll(userId int) ([]maksimzhashkevychtodoapp.TodoList, error)
	GetById(userId, listId int) (maksimzhashkevychtodoapp.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input maksimzhashkevychtodoapp.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, item maksimzhashkevychtodoapp.TodoItem) (int, error)
	GetAll(userId, listId int) ([]maksimzhashkevychtodoapp.TodoItem, error)
	GetById(userId, itemId int) (maksimzhashkevychtodoapp.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input maksimzhashkevychtodoapp.UpdateItemInput) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
