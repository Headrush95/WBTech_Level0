package models

import (
	"fmt"
	"strconv"
	"time"
)

type Payment struct {
	Transaction  string      `json:"transaction" db:"transaction" validate:"required,max=19"`
	RequestId    string      `json:"request_id" db:"request_id"`
	Currency     string      `json:"currency" db:"currency" validate:"required,max=5"`
	Provider     string      `json:"provider" db:"provider" validate:"required,max=255"`
	Amount       int         `json:"amount" db:"amount" validate:"required,gt=0"` // не использую uint во избежания конфликтов с БД
	PaymentDate  paymentDate `json:"payment_dt" db:"payment_dt" validate:"required"`
	Bank         string      `json:"bank" db:"bank" validate:"required,max=255"`
	DeliveryCost int         `json:"delivery_cost" db:"delivery_cost" validate:"required,gt=0"`
	GoodsTotal   int         `json:"goods_total" db:"goods_total" validate:"required,gt=0"`
	CustomFee    int         `json:"custom_fee" db:"custom_fee" validate:"gte=0"`
}

/*
Теоретически можно оставить поле в БД как bigint (для даты в unix формате),
а в структуре оставить int, но тогда данные таблицы будет сложно анализировать по дате.
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
