package repository

import (
	"fmt"
	"strings"

	maksimzhashkevychtodoapp "github.com/LionJr/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list maksimzhashkevychtodoapp.TodoList) (int, error) {

	// transaction --> birnace operation yzygiderliligi bolup, bilelikde umumy bir isi yerine yetiryarler,
	// transaction-daky operationlar ya hemmesi yerine yetyar ya hic haysy
	// transaction doretmek ucin db.Begin() method ulanylyar

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	// QueryRow() method tablisa bir zat yazyp yzyna gaytarylmaly value bar bolsa ulanylyar
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		// Rollback() transaction saklayar
		tx.Rollback()
		return 0, nil
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	// Exec() method tablisa bir zat yazyp yzyna value gaytarmaly dal bolsa ulanylyar
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Commit() method db edilen uytgesmeleri save edyar we transaction tamamlayar
	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]maksimzhashkevychtodoapp.TodoList, error) {
	var lists []maksimzhashkevychtodoapp.TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)
	// Select() method edil Get() yaly yone birnace znaceniya almak ucin ulanylyar
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (maksimzhashkevychtodoapp.TodoList, error) {
	var list maksimzhashkevychtodoapp.TodoList
	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl 
					      INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable)
	// Select() method edil Get() yaly yone birnace znaceniya almak ucin ulanylyar
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListPostgres) Update(userId, listId int, input maksimzhashkevychtodoapp.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	// title=$1
	// description=$1
	// title=$1, description=$2
	setQuery := strings.Join(setValues, ", ")

	//bu yerde biz list update etjek bolyan usere sol listyn degislidigini barlayas (ul.list_id =$%d AND ul.user_id=$%d)
	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id =$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
