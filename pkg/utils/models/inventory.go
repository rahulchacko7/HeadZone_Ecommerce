package models

type Category struct {
	Id       uint   `json:"id" gorm:"unique; not null"`
	Category string `json:"category" gorm:"unique;not null"`
}

type SetNewName struct {
	Current string `json:"current"`
	New     string `json:"new"`
}

type InventoryResponse struct {
	ProductID   int     `json:"id"`
	CategoryID  uint    `json:"category_id"`
	ProductName string  `json:"product_name"`
	Color       string  `json:"color" `
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
}

type InventoryUpdate struct {
	Productid int `json:"product_id"`
	Stock     int `json:"stock"`
}

type Inventory struct {
	ID          uint   `json:"id"`
	CategoryID  int    `json:"category_id"`
	ProductName string `json:"productname"`
	Color       string `json:"color"`
	Stock       int    `json:"stock"`
	Price       int    `json:"price"`
	//	IfPresentAtWishlist bool    `json:"if_present_at_wishlist"`
	//	IfPresentAtCart bool    `json:"if_present_at_cart"`
	//	DiscountedPrice float64 `json:"discounted_price"`
}

type AddInventories struct {
	ID          uint    `json:"id"`
	CategoryID  int     `json:"category_id"`
	ProductName string  `json:"product_name"`
	Color       string  `json:"color"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
}

type EditInventoryDetials struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	CategoryID int     `json:"category_id"`
	Color      string  `json:"color"`
}

type InventoryUserResponse struct {
	ID          uint    `json:"id"`
	CategoryID  int     `json:"category_id"`
	ProductName string  `json:"productname"`
	Color       string  `json:"color"`
	Price       int     `json:"price"`
	Rating      float64 `json:"rating"`
}

type AddToCart struct {
	UserID      int `json:"user_id"`
	InventoryID int `json:"inventory_id"`
}

type InventoryDetails struct {
	Price    float64 `json:"price"`
	Quantity uint    `json:"quantity"`
}

type MakeOrder struct {
	UserID          int       `json:"user_id"`
	Products        []GetCart `json:"products"`
	AddressID       int       `json:"address_id"`
	PaymentMethodID int       `json:"payment_id"`
	FinalPrice      float64   `json:"final_price"`
}

type Order struct {
	UserID          int `json:"user_id"`
	AddressID       int `json:"address_id"`
	PaymentMethodID int `json:"payment_id"`
	CouponID        int `json:"coupon_id"`
}

type SearchItems struct {
	ProductName string `json:"product_name"`
}
