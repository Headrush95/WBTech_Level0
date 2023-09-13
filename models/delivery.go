package models

type Delivery struct {
	Name    string `json:"name" db:"name" validate:"required,min=3,max=255"`
	Phone   string `json:"phone" db:"phone" validate:"required,min=5,max=255"`
	Zip     string `json:"zip" db:"zip" validate:"required,max=10"`
	City    string `json:"city"  db:"city" validate:"required,min=3,max=255"`
	Address string `json:"address"  db:"address" validate:"required,min=10,max=255"`
	Region  string `json:"region"  db:"region"`
	Email   string `json:"email"  db:"email" validate:"required,max=255,email"`
}
