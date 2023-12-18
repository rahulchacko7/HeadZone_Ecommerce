package usecase

import (
	"HeadZone/pkg/config"
	domain "HeadZone/pkg/domain"
	helper "HeadZone/pkg/helper/interfaces"

	"HeadZone/pkg/usecase/interfaces"

	repo "HeadZone/pkg/repository/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo            repo.UserRepository
	cfg                 config.Config
	otpRepository       repo.OtpRepository
	inventoryRepository repo.InventoryRepository
	helper              helper.Helper
}

func NewUserUseCase(repo repo.UserRepository, cfg config.Config, otp repo.OtpRepository, inv repo.InventoryRepository, h helper.Helper) interfaces.UserUseCase {
	return &userUseCase{
		userRepo:            repo,
		cfg:                 cfg,
		otpRepository:       otp,
		inventoryRepository: inv,
		helper:              h,
	}
}

var InternalError = "Internal Server Error"
var ErrorHashingPassword = "Error In Hashing Password"

func (u *userUseCase) UserSignUp(user models.UserDetails) (models.TokenUsers, error) {
	if user.Name == "" {
		return models.TokenUsers{}, errors.New("username cannot be empty")
	}
	namevalidate, err := u.helper.ValidateDatatype(user.Name, "string")
	if err != nil {
		return models.TokenUsers{}, errors.New("invalid format for name")
	}
	if !namevalidate {
		return models.TokenUsers{}, errors.New("not a string")
	}

	userExist := u.userRepo.CheckUserAvailability(user.Email)
	if userExist {
		return models.TokenUsers{}, errors.New("user already exist, sign in")
	}

	phonenumber := u.helper.ValidatePhoneNumber(user.Phone)
	if !phonenumber {
		return models.TokenUsers{}, errors.New("invalid phone")
	}

	if user.Password != user.ConfirmPassword {
		return models.TokenUsers{}, errors.New("password does not match")
	}
	if user.Password == "" {
		return models.TokenUsers{}, errors.New("password cannot be empty")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return models.TokenUsers{}, errors.New("internal server error")
	}
	user.Password = string(hashedPassword)

	userData, err := u.userRepo.UserSignUp(user)
	if err != nil {
		return models.TokenUsers{}, errors.New("could not add the user")
	}

	// crete a JWT token string for the user
	tokenString, err := u.helper.GenerateTokenClients(userData)
	if err != nil {
		return models.TokenUsers{}, errors.New("could not create token due to some internal error")
	}
	return models.TokenUsers{
		Users: userData,
		Token: tokenString,
	}, nil
}

func (u *userUseCase) LoginHandler(user models.UserLogin) (models.TokenUsers, error) {

	// checking if a username exist with this email address
	ok := u.userRepo.CheckUserAvailability(user.Email)
	if !ok {
		return models.TokenUsers{}, errors.New("the user does not exist")
	}

	isBlocked, err := u.userRepo.UserBlockStatus(user.Email)
	if err != nil {
		return models.TokenUsers{}, errors.New(InternalError)
	}

	if isBlocked {
		return models.TokenUsers{}, errors.New("user is blocked by admin")
	}

	// Get the user details in order to check the password, in this case ( The same function can be reused in future )
	user_details, err := u.userRepo.FindUserByEmail(user)
	if err != nil {
		return models.TokenUsers{}, errors.New(InternalError)
	}

	err = u.helper.CompareHashAndPassword(user_details.Password, user.Password)
	if err != nil {
		return models.TokenUsers{}, errors.New("password incorrect")
	}

	var userDetails models.UserDetailsResponse

	userDetails.Id = int(user_details.Id)
	userDetails.Name = user_details.Name
	userDetails.Email = user_details.Email
	userDetails.Phone = user_details.Phone

	tokenString, err := u.helper.GenerateTokenClients(userDetails)
	if err != nil {
		return models.TokenUsers{}, errors.New("could not create token")
	}

	return models.TokenUsers{
		Users: userDetails,
		Token: tokenString,
	}, nil

}

func (i *userUseCase) GetUserDetails(id int) (models.UserDetailsResponse, error) {

	details, err := i.userRepo.GetUserDetails(id)
	if err != nil {
		return models.UserDetailsResponse{}, errors.New("error in getting details")
	}

	return details, nil

}

func (u *userUseCase) AddAddress(id int, address models.AddAddress) error {

	if address.Name == "" || address.HouseName == "" || address.Street == "" || address.City == "" || address.State == "" || address.Phone == "" || address.Pin == "" {
		return errors.New("field cannot be empty")
	}
	ok, err := u.helper.ValidateAlphabets(address.Name)
	if err != nil {
		return errors.New("invalid format for name")
	}
	if !ok {
		return errors.New("invalid format for name")
	}
	phonenumber := u.helper.ValidatePhoneNumber(address.Phone)
	if !phonenumber {
		return errors.New("invalid phone")
	}
	pinnumber := u.helper.ValidatePin(address.Pin)
	if !pinnumber {
		return errors.New("invalid pin number")
	}

	rslt := u.userRepo.CheckIfFirstAddress(id)
	var result bool

	if !rslt {
		result = true
	} else {
		result = false
	}

	err = u.userRepo.AddAddress(id, address, result)
	if err != nil {
		return errors.New("error in adding address")
	}
	return nil

}

func (i *userUseCase) GetAddresses(id int) ([]domain.Address, error) {

	if id <= 0 {
		return []domain.Address{}, errors.New("invalid id")
	}

	addresses, err := i.userRepo.GetAddresses(id)
	if err != nil {
		return []domain.Address{}, errors.New("error in getting addresses")
	}

	return addresses, nil

}

func (i *userUseCase) EditDetails(id int, user models.EditDetailsResponse) (models.EditDetailsResponse, error) {

	if id <= 0 {
		return models.EditDetailsResponse{}, errors.New("invalid id")
	}

	body, err := i.userRepo.EditDetails(id, user)
	if err != nil {
		return models.EditDetailsResponse{}, err
	}

	return body, nil

}

func (i *userUseCase) ChangePassword(id int, old string, password string, repassword string) error {

	if id <= 0 {
		return errors.New("invalid id")
	}

	userPassword, err := i.userRepo.GetPassword(id)
	if err != nil {
		return errors.New(InternalError)
	}

	err = i.helper.CompareHashAndPassword(userPassword, old)
	if err != nil {
		return errors.New("password incorrect")
	}

	if password != repassword {
		return errors.New("passwords does not match")
	}

	newpassword, err := i.helper.PasswordHashing(password)
	if err != nil {
		return errors.New("error in hashing password")
	}

	return i.userRepo.ChangePassword(id, string(newpassword))

}

func (u *userUseCase) GetCart(id int) (models.GetCartResponse, error) {

	if id <= 0 {
		return models.GetCartResponse{}, errors.New("invalid id")
	}

	//find cart id
	cart_id, err := u.userRepo.GetCartID(id)
	if err != nil {
		return models.GetCartResponse{}, errors.New(InternalError)
	}
	//find products inide cart
	products, err := u.userRepo.GetProductsInCart(cart_id)
	if err != nil {
		return models.GetCartResponse{}, errors.New(InternalError)
	}
	//find product names
	var product_names []string
	for i := range products {
		product_name, err := u.userRepo.FindProductNames(products[i])
		if err != nil {
			return models.GetCartResponse{}, errors.New(InternalError)
		}
		product_names = append(product_names, product_name)
	}

	//find quantity
	var quantity []int
	for i := range products {
		q, err := u.userRepo.FindCartQuantity(cart_id, products[i])
		if err != nil {
			return models.GetCartResponse{}, errors.New(InternalError)
		}
		quantity = append(quantity, q)
	}

	var price []float64
	for i := range products {
		q, err := u.userRepo.FindPrice(products[i])
		if err != nil {
			return models.GetCartResponse{}, errors.New(InternalError)
		}
		price = append(price, q)
	}

	var categories []int
	for i := range products {
		c, err := u.userRepo.FindCategory(products[i])
		if err != nil {
			return models.GetCartResponse{}, errors.New(InternalError)
		}
		categories = append(categories, c)
	}

	var getcart []models.GetCart
	for i := range product_names {
		var get models.GetCart
		get.ID = products[i]
		get.ProductName = product_names[i]
		get.Category_id = categories[i]
		get.Quantity = quantity[i]
		get.Price = int(price[i])
		get.Total = (price[i]) * float64(quantity[i])

		getcart = append(getcart, get)
	}

	var response models.GetCartResponse
	response.ID = cart_id
	response.Data = getcart
	//then return in appropriate format

	return response, nil

}

func (i *userUseCase) RemoveFromCart(cart, inventory int) error {

	if cart <= 0 || inventory <= 0 {
		return errors.New("enter a valid number")
	}

	err := i.userRepo.RemoveFromCart(cart, inventory)
	if err != nil {
		return err
	}

	return nil

}

func (i *userUseCase) UpdateQuantity(id, inv, qty int) error {

	if id <= 0 || inv <= 0 || qty <= 0 {
		return errors.New("enter a valid number")
	}
	stock, err := i.inventoryRepository.CheckStock(inv)
	if err != nil {
		return err
	}

	if qty > stock {
		return errors.New("out of stock")
	}

	err = i.userRepo.UpdateQuantity(id, inv, qty)
	return err
}
