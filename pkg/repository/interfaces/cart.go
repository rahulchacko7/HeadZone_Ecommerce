package interfaces

import "HeadZone/pkg/utils/models"

type CartRepository interface {
	GetCart(id int) ([]models.GetCart, error)
	GetCartId(user_id int) (int, error)
	CreateNewCart(user_id int) (int, error)
	AddLineItems(cart_id, inventory_id int) error
	CheckIfItemIsAlreadyAdded(cart_id, inventory_id int) (bool, error)
}
