package usecase

import (
	"HeadZone/pkg/config"
	mockRepository "HeadZone/pkg/repository/Repomock"
	"HeadZone/pkg/utils/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_AddAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockRepository.NewMockUserRepository(ctrl)
	cfg := config.Config{}
	otpRepo := mockRepository.NewMockOtpRepository(ctrl)
	inventoryRepo := mockRepository.NewMockInventoryRepository(ctrl)
	helper := mockRepository.NewMockHelper(ctrl)

	userUseCase := NewUserUseCase(userRepo, cfg, otpRepo, inventoryRepo, helper)

	testData := map[string]struct {
		input   models.AddAddress
		stub    func(*mockRepository.MockUserRepository, models.AddAddress)
		wantErr error
	}{
		"success": {
			input: models.AddAddress{
				Name:      "Rahul",
				HouseName: "thakadiyil house",
				Street:    "pallippuram",
				City:      "cherthala",
				State:     "kerala",
				Pin:       "688541",
			},
			stub: func(userRepo *mockRepository.MockUserRepository, data models.AddAddress) {
				userRepo.EXPECT().AddAddress(1, data, "").Return(nil).Times(1)
			},
			wantErr: nil,
		},
		"failure": {
			input: models.AddAddress{
				Name:      "akhil",
				HouseName: "chekkiyil house",
				Street:    "pallippuram",
				City:      "cherthala",
				State:     "kerala",
				Pin:       "688541",
			},
			stub: func(userRepo *mockRepository.MockUserRepository, data models.AddAddress) {
				userRepo.EXPECT().AddAddress(1, data, "").Return(errors.New("could not add the address")).Times(1)
			},
			wantErr: errors.New("could not add the address"),
		},
	}
	for testName, test := range testData {
		t.Run(testName, func(t *testing.T) {
			test.stub(userRepo, test.input)
			err := userUseCase.AddAddress(1, test.input)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
