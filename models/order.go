package models

import (
	"time"
)

type Order struct {
	Uid               string    `json:"order_uid" db:"order_uid" binding:"required"`
	TrackNumber       string    `json:"track_number" db:"track_number" binding:"required"`
	Entry             string    `json:"entry" db:"entry" binding:"required"`
	Delivery          Delivery  `json:"delivery" db:"delivery" binding:"required"`
	Payment           Payment   `json:"payment" db:"payment" binding:"required"`
	Items             []Item    `json:"items" binding:"required"`
	Locale            string    `json:"locale" db:"locale" binding:"required"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerId        string    `json:"customer_id" db:"customer_id" binding:"required"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service" binding:"required"`
	ShardKey          string    `json:"shardkey" db:"shardkey" binding:"required"`
	SmId              int       `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" db:"date_created" binding:"required"`
	OofShard          string    `json:"oof_shard" db:"oof_shard" binding:"required"`
}
