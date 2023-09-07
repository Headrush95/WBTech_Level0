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
			"error": fmt.Sprintf("error occurred while parsing request body: %s", err.Error()),
		})
		return
	}

	//// TODO DELETE
	//b, _ := json.MarshalIndent(order, "", "	")
	//fmt.Println(string(b))
	//// TODO DELETE

	err = h.services.CreateOrder(order)
	if err != nil {
		logrus.Errorf("error occurred while creating order DB entry: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
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

	//order, err := h.services.GetOrder(uid)
	order, err := h.services.GetOrderById(uid) // TODO поменять на поиск в кэше!!!
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error occurred while trying get order with uid %s: %s", uid, err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, order)

}
