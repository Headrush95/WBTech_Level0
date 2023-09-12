package models

import (
	"time"
)

// TODO добавить валидацию

type Order struct {
	Uid               string    `json:"order_uid" db:"order_uid" validate:"required,max=19"`
	TrackNumber       string    `json:"track_number" db:"track_number" validate:"required,max=255"`
	Entry             string    `json:"entry" db:"entry" validate:"required,max=50"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,dive,required"`
	Locale            string    `json:"locale" db:"locale" validate:"required,max=2"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature,max=255"`
	CustomerId        string    `json:"customer_id" db:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service" validate:"required,max=255"`
	ShardKey          string    `json:"shardkey" db:"shardkey" validate:"required,max=5"`
	SmId              int       `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" db:"date_created" format:"2006-01-02T06:22:19Z" validate:"required,lte"`
	OofShard          string    `json:"oof_shard" db:"oof_shard" validate:"required,max=5"`
}
