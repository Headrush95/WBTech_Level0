package models

import (
	"fmt"
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

/*
	теоретически можно оставить просто поле в БД как bigint (для даты в unix формате),

но тогда данные таблицы будет сложно анализировать по дате.
*/
type paymentDate struct {
	time.Time
}

// UnmarshalJSON для занесения даты в БД
func (d *paymentDate) UnmarshalJSON(b []byte) (err error) {
	unixDate, err := strconv.Atoi(string(b))
	if err != nil {
		return
	}

	d.Time = time.Unix(int64(unixDate), 0)
	return
}

// SetDate костыль с преобразование даты из БД в paymentDate. Чую можно проще, но пока не знаю как.
func (p *Payment) SetDate(date time.Time) {
	p.PaymentDate = paymentDate{date}
}

// Scan для получения даты из БД
func (d *paymentDate) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	if t, ok := value.(time.Time); ok {
		d.Time = t
		return nil
	}
	return fmt.Errorf("failed to scan paymentDate: unexpected value type")
}
