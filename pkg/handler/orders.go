package handler

import (
	"WBTech_Level0/models"
	"fmt"
	"github.com/gin-gonic/gin"
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
			"error": fmt.Sprintf("error occured while parsing request body: %s", err.Error()),
		})
		return
	}

	err = h.services.CreateOrder(order)
	if err != nil {
		logrus.Errorf("error occured while creating order DB entry: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error occured while creating order DB entry: %s", err.Error()),
		})
		return
	}

	// запись в кэш
	err = h.services.PutOrder(order)
	if err != nil {
		logrus.Errorf("error occured while putting order data into cache: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error occured while putting order data into cache: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_uid": order.Uid,
	})
}

type orderRequest struct {
	Uid string `json:"uid"`
}

func (h *Handler) GetOrderById(c *gin.Context) {
	var req orderRequest

	err := c.BindJSON(&req)
	if err != nil {
		logrus.Errorf("[%s] error occured while parsing request body: %s", c.ClientIP(), err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if len(req.Uid) != 19 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "incorrect order uid",
		})
	}

	order, err := h.services.GetOrder(req.Uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("there is no order with uid %s", req.Uid),
		})
	}

	c.JSON(http.StatusOK, order)

}
