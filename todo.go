package maksimzhashkevychtodoapp

import "errors"

type TodoList struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UsersList struct {
	Id     int
	UserId int
	ListId int
}

type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

type ListsItem struct {
	Id     int
	ListId int
	ItemId int
}

type UpdateListInput struct {
	// *string ptr etsen string nil bolup bilyar, sebabi string default "" bolyar
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

// update ucin berlen Title we Description nil bolup bilyanligi ucin olary validate etmeli
func (u UpdateListInput) Validate() error {
	if u.Title == nil && u.Description == nil {
		return errors.New("no values to update")
	}
	return nil
}

type UpdateItemInput struct {
	// *string ptr etsen string nil bolup bilyar, sebabi string default "" bolyar
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

// update ucin berlen Title, Description we Done nil bolup bilyanligi ucin olary validate etmeli
func (i UpdateItemInput) Validate() error {
	if i.Title == nil && i.Description == nil && i.Done != nil {
		return errors.New("no values to update")
	}
	return nil
}
