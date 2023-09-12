package models

// TODO добавить валидацию

type Item struct {
	ChrtId      int    `json:"chrt_id"  db:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" db:"track_number" validate:"required,min=10,max=255"`
	Price       uint   `json:"price" db:"price" validate:"required,gt=0"`
	Rid         string `json:"rid" db:"rid" validate:"required,max=25"`
	Name        string `json:"name" db:"name" validate:"required,max=255"`
	Sale        uint8  `json:"sale" db:"sale" validate:"max=100"`
	Size        string `json:"size" db:"size" validate:"required,max=5"`
	TotalPrice  uint   `json:"total_price" db:"total_price" validate:"required,gt=0"`
	NmId        int    `json:"nm_id" db:"nm_id" validate:"required"`
	Brand       string `json:"brand" db:"brand" validate:"required,max=255"`
	Status      int    `json:"status" db:"status" validate:"required"`
}
