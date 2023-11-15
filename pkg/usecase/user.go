package usecase

import (
	"HeadZone/pkg/config"
	domain "HeadZone/pkg/domain"
	helper "HeadZone/pkg/helper/interfaces"

	"HeadZone/pkg/usecase/interfaces"

	repo "HeadZone/pkg/repository/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
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
	// Check whether the user already exist. If yes, show the error message, since this is signUp
	userExist := u.userRepo.CheckUserAvailability(user.Email)
	if userExist {
		return models.TokenUsers{}, errors.New("user already exist, sign in")
	}
	if user.Password != user.ConfirmPassword {
		return models.TokenUsers{}, errors.New("password does not match")
	}

	// Hash password since details are validated

	hashedPassword, err := u.helper.PasswordHashing(user.Password)
	if err != nil {
		return models.TokenUsers{}, errors.New(ErrorHashingPassword)
	}

	user.Password = hashedPassword

	// add user details to the database
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

func (i *userUseCase) AddAddress(id int, address models.AddAddress) error {

	rslt := i.userRepo.CheckIfFirstAddress(id)
	var result bool

	if !rslt {
		result = true
	} else {
		result = false
	}

	err := i.userRepo.AddAddress(id, address, result)
	if err != nil {
		return errors.New("error in adding address")
	}

	return nil

}

func (i *userUseCase) GetAddresses(id int) ([]domain.Address, error) {

	addresses, err := i.userRepo.GetAddresses(id)
	if err != nil {
		return []domain.Address{}, errors.New("error in getting addresses")
	}

	return addresses, nil

}

func (i *userUseCase) EditName(id int, name string) error {

	err := i.userRepo.EditName(id, name)
	if err != nil {
		return errors.New("could not change")
	}

	return nil

}

func (i *userUseCase) EditEmail(id int, email string) error {

	err := i.userRepo.EditEmail(id, email)
	if err != nil {
		return errors.New("could not change")
	}

	return nil

}

func (i *userUseCase) EditPhone(id int, phone string) error {

	err := i.userRepo.EditPhone(id, phone)
	if err != nil {
		return errors.New("could not change")
	}

	return nil

}

func (i *userUseCase) ChangePassword(id int, old string, password string, repassword string) error {

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
