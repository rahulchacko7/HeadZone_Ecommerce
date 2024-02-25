package usecase

import (
	"HeadZone/pkg/config"
	domain "HeadZone/pkg/domain"
	mockhelper "HeadZone/pkg/helper/mock"
	mockRepository "HeadZone/pkg/repository/Repomock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_GetAddresses(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockRepository.NewMockUserRepository(ctrl)
	cfg := config.Config{}
	otpRepo := mockRepository.NewMockOtpRepository(ctrl)
	inventoryRepo := mockRepository.NewMockInventoryRepository(ctrl)
	helper := mockhelper.NewMockHelper(ctrl)

	userUseCase := NewUserUseCase(userRepo, cfg, otpRepo, inventoryRepo, helper)

	testData := map[string]struct {
		input   int
		stub    func(*mockRepository.MockUserRepository, *mockhelper.MockHelper, int)
		want    []domain.Address
		wantErr error
	}{
		"success": {
			input: 1,
			stub: func(userrepo *mockRepository.MockUserRepository, helper *mockhelper.MockHelper, data int) {
				userrepo.EXPECT().GetAddresses(data).Times(1).Return([]domain.Address{}, nil)
			},
			want:    []domain.Address{},
			wantErr: nil,
		},
		"failed": {
			input: 1,
			stub: func(userrepo *mockRepository.MockUserRepository, helper *mockhelper.MockHelper, data int) {
				userrepo.EXPECT().GetAddresses(data).Times(1).Return([]domain.Address{}, errors.New("error"))
			},
			want:    []domain.Address{},
			wantErr: errors.New("error in getting addresses"), // Corrected error string
		},
	}
	for _, test := range testData {
		test.stub(userRepo, helper, test.input)
		result, err := userUseCase.GetAddresses(test.input)
		assert.Equal(t, test.want, result)
		assert.Equal(t, test.wantErr, err)
	}
}
