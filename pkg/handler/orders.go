package handler

import (
	"WBTech_Level0/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) CreateOrder(c *gin.Context) {
	var order models.Order
	err := c.BindJSON(&order)
	if err != nil {
		logrus.Errorf("error occured while parsing request body: %s", err.Error())

		//TODO можно ли как-то записать тело ответа не трогая код статуса ответа, так как BindJSON и так записывает код 400??
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("error occurred while parsing request body: %s", err.Error()),
		})
		return
	}

	// запись в БД
	err = h.services.CreateOrder(order)
	if err != nil {
		logrus.Errorf("error occurred while creating order DB entry: %s", err.Error())
		status := 0
		if ok := errors.As(err, &validator.ValidationErrors{}); ok {
			status = http.StatusBadRequest
		} else {
			status = http.StatusInternalServerError
		}

		c.AbortWithStatusJSON(status, gin.H{
			"error": fmt.Sprintf("error occurred while creating order DB entry: %s", err.Error()),
		})
		return
	}

	// запись в кэш
	err = h.services.PutOrder(order)
	if err != nil {
		logrus.Errorf("error occurred while putting order data into cache: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error occured while putting order data into cache: %s", err.Error()),
		})
		return
	}

	// тело ответа для удобства проверки работоспособности
	c.JSON(http.StatusOK, gin.H{
		"order_uid": order.Uid,
	})
}

func (h *Handler) GetOrderById(c *gin.Context) {
	uid := c.Param("id")
	if len(uid) != 19 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "incorrect order uid",
		})
		return
	}

	order, err := h.services.GetOrder(uid)
	//order, err := h.services.GetOrderById(uid) // TODO поменять на поиск в кэше!!!
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error occurred while trying get order with uid %s: %s", uid, err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, order)

}

func (h *Handler) GetAllOrders(c *gin.Context) {
	orders, err := h.services.GetAllOrders()
	if err != nil {
		// формально, не всегда код ошибки 500
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error occurred while trying get all orders: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, orders)
}
