package interfaces

type CartUseCase interface {
	AddToCart(user_id, inventory_id int) error
}
