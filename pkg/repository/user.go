package repository

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/repository/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB}
}

func (c *userDatabase) CheckUserAvailability(email string) bool {

	var count int
	query := fmt.Sprintf("select count(*) from users where email='%s'", email)
	if err := c.DB.Raw(query).Scan(&count).Error; err != nil {
		return false
	}
	// if count is greater than 0 that means the user already exist
	return count > 0

}

func (c *userDatabase) UserSignUp(user models.UserDetails) (models.UserDetailsResponse, error) {
	var userDetails models.UserDetailsResponse

	// Insert user details into the users table
	err := c.DB.Raw("INSERT INTO users (name, email, password, phone) VALUES (?, ?, ?, ?) RETURNING id, name, email, phone", user.Name, user.Email, user.Password, user.Phone).Scan(&userDetails).Error
	if err != nil {
		return models.UserDetailsResponse{}, err
	}

	// Get the ID of the newly created user
	newUserID := userDetails.Id

	fmt.Println("userId at signup", newUserID)

	// Create a wallet entry for the new user with an initial amount of 0
	// err = c.DB.Exec("INSERT INTO wallets (user_id, amount) VALUES (?, ?)", newUserID, 0).Error
	// if err != nil {
	// 	return models.UserDetailsResponse{}, err
	// }

	return userDetails, nil
}

func (cr *userDatabase) UserBlockStatus(email string) (bool, error) {
	fmt.Println(email)
	var isBlocked bool
	err := cr.DB.Raw("select blocked from users where email = ?", email).Scan(&isBlocked).Error
	if err != nil {
		return false, err
	}
	fmt.Println(isBlocked)
	return isBlocked, nil
}

func (c *userDatabase) FindUserByEmail(user models.UserLogin) (models.UserSignInResponse, error) {

	var user_details models.UserSignInResponse

	err := c.DB.Raw(`
		SELECT *
		FROM users where email = ? and blocked = false
		`, user.Email).Scan(&user_details).Error

	if err != nil {
		return models.UserSignInResponse{}, errors.New("error checking user details")
	}

	return user_details, nil

}

func (ad *userDatabase) GetUserDetails(id int) (models.UserDetailsResponse, error) {

	var details models.UserDetailsResponse

	if err := ad.DB.Raw("select id,name,email,phone from users where id=?", id).Scan(&details).Error; err != nil {
		return models.UserDetailsResponse{}, errors.New("could not get user details")
	}

	return details, nil

}

func (c *userDatabase) CheckIfFirstAddress(id int) bool {

	var count int
	// query := fmt.Sprintf("select count(*) from addresses where user_id='%s'", id)
	if err := c.DB.Raw("select count(*) from addresses where user_id=$1", id).Scan(&count).Error; err != nil {
		return false
	}
	// if count is greater than 0 that means the user already exist
	return count > 0

}

func (i *userDatabase) AddAddress(id int, address models.AddAddress, result bool) error {
	err := i.DB.Exec(`
		INSERT INTO addresses (user_id, name, house_name, street, city, state, phone, pin,"default")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9 )`,
		id, address.Name, address.HouseName, address.Street, address.City, address.State, address.Phone, address.Pin, result).Error
	if err != nil {
		return errors.New("could not add address")
	}

	return nil
}

func (ad *userDatabase) GetAddresses(id int) ([]domain.Address, error) {

	var addresses []domain.Address

	if err := ad.DB.Raw("select * from addresses where user_id=?", id).Scan(&addresses).Error; err != nil {
		return []domain.Address{}, errors.New("error in getting addresses")
	}

	return addresses, nil

}

func (i *userDatabase) EditDetails(id int, user models.EditDetailsResponse) (models.EditDetailsResponse, error) {

	var body models.EditDetailsResponse

	args := []interface{}{}
	query := "update users set"

	if user.Email != "" {
		query += " email = $1,"

		args = append(args, user.Email)
	}

	if user.Name != "" {
		query += " name = $2,"
		args = append(args, user.Name)
	}

	if user.Phone != "" {
		query += " phone = $3,"

		args = append(args, user.Phone)
	}

	query = query[:len(query)-1] + " where id = $4"

	args = append(args, id)
	// fmt.Println(query, args)
	err := i.DB.Exec(query, args...).Error
	if err != nil {
		return models.EditDetailsResponse{}, err
	}
	query2 := "select * from users where id = ?"
	if err := i.DB.Raw(query2, id).Scan(&body).Error; err != nil {
		return models.EditDetailsResponse{}, err
	}

	return body, nil

}

func (i *userDatabase) ChangePassword(id int, password string) error {

	err := i.DB.Exec("UPDATE users SET password=$1 WHERE id=$2", password, id).Error
	if err != nil {
		return err
	}

	return nil

}

func (i *userDatabase) GetPassword(id int) (string, error) {

	var userPassword string
	err := i.DB.Raw("select password from users where id = ?", id).Scan(&userPassword).Error
	if err != nil {
		return "", err
	}
	return userPassword, nil

}

func (ad *userDatabase) GetCartID(id int) (int, error) {

	var cart_id int

	if err := ad.DB.Raw("select id from carts where user_id=?", id).Scan(&cart_id).Error; err != nil {
		return 0, err
	}

	return cart_id, nil

}

func (ad *userDatabase) GetProductsInCart(cart_id int) ([]int, error) {

	var cart_products []int

	if err := ad.DB.Raw("select inventory_id from line_items where cart_id=?", cart_id).Scan(&cart_products).Error; err != nil {
		return []int{}, err
	}

	return cart_products, nil

}

func (ad *userDatabase) FindProductNames(inventory_id int) (string, error) {

	var product_name string

	if err := ad.DB.Raw("select product_name from inventories where id=?", inventory_id).Scan(&product_name).Error; err != nil {
		return "", err
	}

	return product_name, nil

}

func (ad *userDatabase) FindCartQuantity(cart_id, inventory_id int) (int, error) {

	var quantity int

	if err := ad.DB.Raw("select quantity from line_items where cart_id=$1 and inventory_id=$2", cart_id, inventory_id).Scan(&quantity).Error; err != nil {
		return 0, err
	}

	return quantity, nil

}

func (ad *userDatabase) FindPrice(inventory_id int) (float64, error) {

	var price float64

	if err := ad.DB.Raw("select price from inventories where id=?", inventory_id).Scan(&price).Error; err != nil {
		return 0, err
	}

	return price, nil

}

func (ad *userDatabase) FindCategory(inventoryID int) (int, error) {

	var categoryID int

	if err := ad.DB.Raw("SELECT category_id FROM inventories WHERE id = ?", inventoryID).Scan(&categoryID).Error; err != nil {
		return 0, err
	}

	return categoryID, nil
}

func (i *userDatabase) FindStock(id int) (int, error) {
	var stock int
	err := i.DB.Raw("SELECT stock FROM inventories WHERE id = ?", id).Scan(&stock).Error
	if err != nil {
		return 0, err
	}

	return stock, nil
}

func (ad *userDatabase) RemoveFromCart(cart, inventory int) error {

	if err := ad.DB.Exec(`DELETE FROM line_items WHERE cart_id = $1 AND inventory_id = $2`, cart, inventory).Error; err != nil {
		return err
	}

	return nil

}
func (ad *userDatabase) UpdateQuantity(id, invID, qty int) error {
	if id <= 0 || invID <= 0 || qty <= 0 {
		return errors.New("negative or zero values are not allowed")
	}

	if qty >= 25 {
		return errors.New("choose number of items below 25")
	}

	if qty >= 0 {
		query := `
        UPDATE line_items
        SET quantity = $1
        WHERE cart_id = $2 AND inventory_id = $3
        `

		result := ad.DB.Exec(query, qty, id, invID)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
