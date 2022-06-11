package handler

import (
	"net/http"

	"github.com/MaximumTroubles/go-todo-app"
	"github.com/gin-gonic/gin"
)

func (h *Handler) sighUp(c *gin.Context) {
	// take out User struct and give pointer to BindJson method that expect object, and will parse income json to our properties with same keys
	var input todo.User

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// here we are need login and password only. So struct User we can't use due to field "username" are required
// so we create a new struct 
type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) sighIn(c *gin.Context) {
	// here we use new struct
	var input SignInInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}


