package usecase

import (
	config "HeadZone/pkg/config"
	helper "HeadZone/pkg/helper"
	interfaces "HeadZone/pkg/repository/interfaces"
	services "HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
)

type otpUseCase struct {
	cfg           config.Config
	otpRepository interfaces.OtpRepository
}

func NewOtpUseCase(cfg config.Config, repo interfaces.OtpRepository) services.OtpUseCase {
	return &otpUseCase{
		cfg:           cfg,
		otpRepository: repo,
	}
}

func (ot *otpUseCase) SendOTP(phone string) error {

	ok := ot.otpRepository.FindUserByMobileNumber(phone)
	if !ok {
		return errors.New("the user does not exist")
	}

	helper.TwilioSetup(ot.cfg.ACCOUNTSID, ot.cfg.AUTHTOKEN)
	fmt.Println("accsid:", ot.cfg.ACCOUNTSID)
	fmt.Println("auth:", ot.cfg.AUTHTOKEN)
	_, err := helper.TwilioSendOTP(phone, ot.cfg.SERVICESID)
	if err != nil {
		return errors.New("error ocurred while generating OTP")
	}

	return nil

}

func (ot *otpUseCase) VerifyOTP(code models.VerifyData) (models.TokenUsers, error) {

	helper.TwilioSetup(ot.cfg.ACCOUNTSID, ot.cfg.AUTHTOKEN)
	err := helper.TwilioVerifyOTP(ot.cfg.SERVICESID, code.Code, code.PhoneNumber)
	if err != nil {
		//this guard clause catches the error code runs only until here
		return models.TokenUsers{}, errors.New("error while verifying")
	}

	// if user is authenticated using OTP send back user details
	userDetails, err := ot.otpRepository.UserDetailsUsingPhone(code.PhoneNumber)
	if err != nil {
		return models.TokenUsers{}, err
	}

	tokenString, err := helper.GenerateTokenClients(userDetails)
	if err != nil {
		return models.TokenUsers{}, err
	}

	var user models.UserDetailsResponse
	err = copier.Copy(&user, &userDetails)
	if err != nil {
		return models.TokenUsers{}, err
	}

	return models.TokenUsers{
		Users: user,
		Token: tokenString,
	}, nil

}
