package handler

import (
	"github.com/MaximumTroubles/go-todo-app/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	// create instance of framework
	router := gin.New()

	auth := router.Group("/auth")
	{
		// it is kind of php router where we sign route path to controller method
		auth.POST("/sigh-up", h.sighUp)
		auth.POST("/sigh-in", h.sighIn)
	}

	api := router.Group("/api", h.userIndetity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateListById)
			lists.DELETE("/:id", h.deleteListById)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}

		items := api.Group("/items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItemById)
			items.DELETE("/:id", h.deleteItemById)
		}
	}

	return router
}
