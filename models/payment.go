package models

import (
	"strconv"
	"time"
)

type Payment struct {
	Transaction  string      `json:"transaction" db:"transaction" binding:"required"`
	RequestId    string      `json:"request_id" db:"request_id"`
	Currency     string      `json:"currency" db:"currency" binding:"required"`
	Provider     string      `json:"provider" db:"provider" binding:"required"`
	Amount       int         `json:"amount" db:"amount" binding:"required"`
	PaymentDate  paymentDate `json:"payment_dt" db:"payment_dt" binding:"required"`
	Bank         string      `json:"bank" db:"bank" binding:"required"`
	DeliveryCost int         `json:"delivery_cost" db:"delivery_cost" binding:"required"`
	GoodsTotal   int         `json:"goods_total" db:"goods_total" binding:"required"`
	CustomFee    int         `json:"custom_fee" db:"custom_fee"`
}

type paymentDate struct {
	time.Time
}

func (d *paymentDate) UnmarshalJSON(b []byte) (err error) {
	unixDate, err := strconv.Atoi(string(b))
	if err != nil {
		return
	}

	d.Time = time.Unix(int64(unixDate), 0)
	return
}
