// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/repository/interfaces/otp.go
//
// Generated by this command:
//
//	mockgen -source pkg/repository/interfaces/otp.go -destination pkg/usecase/mock/otp_mock.go -package mock
//

// Package mock is a generated GoMock package.
package mock

import (
	models "HeadZone/pkg/utils/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockOtpRepository is a mock of OtpRepository interface.
type MockOtpRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOtpRepositoryMockRecorder
}

// MockOtpRepositoryMockRecorder is the mock recorder for MockOtpRepository.
type MockOtpRepositoryMockRecorder struct {
	mock *MockOtpRepository
}

// NewMockOtpRepository creates a new mock instance.
func NewMockOtpRepository(ctrl *gomock.Controller) *MockOtpRepository {
	mock := &MockOtpRepository{ctrl: ctrl}
	mock.recorder = &MockOtpRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOtpRepository) EXPECT() *MockOtpRepositoryMockRecorder {
	return m.recorder
}

// FindUserByMobileNumber mocks base method.
func (m *MockOtpRepository) FindUserByMobileNumber(phone string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByMobileNumber", phone)
	ret0, _ := ret[0].(bool)
	return ret0
}

// FindUserByMobileNumber indicates an expected call of FindUserByMobileNumber.
func (mr *MockOtpRepositoryMockRecorder) FindUserByMobileNumber(phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByMobileNumber", reflect.TypeOf((*MockOtpRepository)(nil).FindUserByMobileNumber), phone)
}

// UserDetailsUsingPhone mocks base method.
func (m *MockOtpRepository) UserDetailsUsingPhone(phone string) (models.UserDetailsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserDetailsUsingPhone", phone)
	ret0, _ := ret[0].(models.UserDetailsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserDetailsUsingPhone indicates an expected call of UserDetailsUsingPhone.
func (mr *MockOtpRepositoryMockRecorder) UserDetailsUsingPhone(phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserDetailsUsingPhone", reflect.TypeOf((*MockOtpRepository)(nil).UserDetailsUsingPhone), phone)
}
