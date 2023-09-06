package models

type Item struct {
	ChrtId      int    `json:"chrt_id"  db:"chrt_id" binding:"required"`
	TrackNumber string `json:"track_number" db:"track_number" binding:"required"`
	Price       uint   `json:"price" db:"price" binding:"required"`
	Rid         string `json:"rid" db:"rid" binding:"required"`
	Name        string `json:"name" db:"name" binding:"required"`
	Sale        uint8  `json:"sale" db:"sale"`
	Size        string `json:"size" db:"size" binding:"required"`
	TotalPrice  uint   `json:"total_price" db:"total_price" binding:"required"`
	NmId        int    `json:"nm_id" db:"nm_id" binding:"required"`
	Brand       string `json:"brand" db:"brand" binding:"required"`
	Status      int    `json:"status" db:"status" binding:"required"`
}
