package domain

type InventoryUpdate struct {
	Productid int `json:"product_id"`
	Stock     int `json:"stock"`
}

type Category struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Category string `json:"category"`
}

type Inventory struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"category" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	ProductName string   `json:"product_name"`
	Color       string   `json:"color" gorm:"color:4;default:'Black';Check:color IN ('Black', 'Blue', 'Red', 'Green');"`
	Stock       int      `json:"stock"`
	Price       float64  `json:"price"`
}

type Rating struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	Productid int     `json:"product_id"`
	Rating    float64 `json:"rating"`
}
