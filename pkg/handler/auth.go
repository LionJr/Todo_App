package handler

import (
	"net/http"

	maksimzhashkevychtodoapp "github.com/LionJr/todo-app"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input maksimzhashkevychtodoapp.User

	if err := c.BindJSON(&input); err != nil {
		// status code == 400 yalnys data berilse
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		// status code == 500 server error
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	// status code == 200 OK
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		// status code == 400 yalnys data berilse
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		// status code == 500 server error
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	// status code == 200 OK
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
