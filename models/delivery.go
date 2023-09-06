package models

type Delivery struct {
	Name    string `json:"name" db:"name" binding:"required"`
	Phone   string `json:"phone" db:"phone" binding:"required"`
	Zip     string `json:"zip" db:"zip" binding:"required"`
	City    string `json:"city"  db:"city" binding:"required"`
	Address string `json:"address"  db:"address" binding:"required"`
	Region  string `json:"region"  db:"region"`
	Email   string `json:"email"  db:"email" binding:"required"`
}
