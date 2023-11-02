package repository

import (
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
	err := c.DB.Raw("INSERT INTO users (name, email, password, phone) VALUES (?, ?, ?, ?) RETURNING id, name, email, phone", user.Name, user.Email, user.Password, user.Phone).Scan(&userDetails).Error

	if err != nil {
		return models.UserDetailsResponse{}, err
	}

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
