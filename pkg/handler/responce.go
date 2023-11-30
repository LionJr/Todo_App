package handler

import (
	maksimzhashkevychtodoapp "github.com/LionJr/todo-app"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponce struct {
	Message string `json:"message"`
}

type statusResponce struct {
	Status string `json:"status"`
}

type getAllListsResponce struct {
	Data []maksimzhashkevychtodoapp.TodoList `json:"data"`
}

func newErrorResponce(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponce{Message: message})
}
