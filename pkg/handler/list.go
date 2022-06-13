package handler

import (
	"net/http"
	"strconv"

	"github.com/MaximumTroubles/go-todo-app"
	"github.com/gin-gonic/gin"
)

// additional struct for wrapping response with data
type getAllListsResponse struct {
	Data []todo.TodoList `json:"data"`
}

func (h *Handler) getAllLists(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// return list of todoLists. Call new reponse struct
	c.JSON(http.StatusOK, getAllListsResponse{
		Data: lists,
	})
}

func (h *Handler) createList(c *gin.Context) {
	// user id we were set at userIdenity method in middleware file in Context object.
	id, err := getUserId(c)
	if err != nil {
		return
	}

	//write down in context user input data, bring error
	var input todo.TodoList
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// here we are handle users input data. call method service.
	// notice that we have to make id type int. because c.Get() method return interface but we expect int
	listId, err := h.services.TodoList.Create(id, input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": listId,
	})
}

func (h *Handler) getListById(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		return
	}

	//We have to get list id from URL (query params) and handle
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return 
	}

	list, err := h.services.TodoList.GetById(id, listId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) updateListById(c *gin.Context) {

}

func (h *Handler) deleteListById(c *gin.Context) {

}
