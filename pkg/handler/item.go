package handler

import (
	"net/http"
	"strconv"

	"github.com/MaximumTroubles/go-todo-app"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllItems(c *gin.Context) {
	// User id from auth
	id, err := getUserId(c)
	if err != nil {
		return
	}

	// list id we get from params (url)
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid list id param")
	}

	items, err := h.services.TodoItem.GetAll(id, listId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) createItem(c *gin.Context) {
	// user id we were set at userIdenity method in middleware file in Context object.
	id, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid list id param")
	}

	//write down in context user input data, bring error
	var input todo.TodoItem
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// here we are handle users input data. call method service.
	// notice that we have to make id type int. because c.Get() method return interface but we expect int
	itemId, err := h.services.TodoItem.Create(id, listId, input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": itemId,
	})
}

func (h *Handler) getItemById(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid item id param")
	}

	item, err := h.services.TodoItem.GetById(id, itemId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) updateItemById(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		return
	}

	//We have to get list id from URL (query params) and handle
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input todo.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// update method
	if err := h.services.TodoItem.Update(id, itemId, input); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteItemById(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid item id param")
	}

	err = h.services.TodoItem.Delete(id, itemId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}
