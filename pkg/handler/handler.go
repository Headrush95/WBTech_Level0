package handler

import (
	"WBTech_Level0/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	orders := router.Group("/orders")
	{
		orders.POST("/create", h.CreateOrder)
		orders.GET("/:id", h.GetOrderById)
		orders.GET("/", h.GetAllOrders)
	}

	return router
}
