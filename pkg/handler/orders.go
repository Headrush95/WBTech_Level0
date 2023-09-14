package handler

import (
	"WBTech_Level0/models"
	"WBTech_Level0/pkg/repository"
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
		// перезаписывать статус код в ответе не требуется, так как BindJSON при ошибке устанавливает код 400
		resp := fmt.Sprintf("error occurred while parsing request body: %s", err.Error())
		logrus.Error(resp)
		c.Writer.Write([]byte(resp))
		return
	}

	// запись в БД
	err = h.services.CreateOrder(order)
	if err != nil {
		logrus.Errorf("error occurred while creating order DB entry: %s", err.Error())
		var status int
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

	// разрешаем всем пользователям доступ к телу ответа
	c.Header("Access-Control-Allow-Origin", "*")
	// тело ответа для удобства проверки работоспособности
	c.JSON(http.StatusOK, gin.H{
		"order_uid": order.Uid,
	})
}

func (h *Handler) GetOrderById(c *gin.Context) {
	uid := c.Param("id")
	// разрешаем всем пользователям доступ к данным заказа
	c.Header("Access-Control-Allow-Origin", "*")
	if len(uid) != 19 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "incorrect order uid",
		})
		return
	}

	order, err := h.services.GetOrder(uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error occurred while trying get order with uid %s: %s", uid, err.Error()),
		})
		return
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, order)
}

func (h *Handler) GetAllOrders(c *gin.Context) {
	orders, err := h.services.GetAllOrders()
	if err != nil {
		if errors.Is(err, repository.EmptyDB) {
			// 200, так как отсутствие заказов не говорит о неправильном запросе или ошибке на сервере
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"info": "there are no orders in DB",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error occurred while trying get all orders: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, orders)
}
