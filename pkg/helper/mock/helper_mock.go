// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/helper/interfaces/helper.go
//
// Generated by this command:
//
//	mockgen -source pkg/helper/interfaces/helper.go -destination pkg/helper/mock/helper_mock.go -package mock
//

// Package mock is a generated GoMock package.
package mock

import (
	models "HeadZone/pkg/utils/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockHelper is a mock of Helper interface.
type MockHelper struct {
	ctrl     *gomock.Controller
	recorder *MockHelperMockRecorder
}

// MockHelperMockRecorder is the mock recorder for MockHelper.
type MockHelperMockRecorder struct {
	mock *MockHelper
}

// NewMockHelper creates a new mock instance.
func NewMockHelper(ctrl *gomock.Controller) *MockHelper {
	mock := &MockHelper{ctrl: ctrl}
	mock.recorder = &MockHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHelper) EXPECT() *MockHelperMockRecorder {
	return m.recorder
}

// CompareHashAndPassword mocks base method.
func (m *MockHelper) CompareHashAndPassword(a, b string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompareHashAndPassword", a, b)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompareHashAndPassword indicates an expected call of CompareHashAndPassword.
func (mr *MockHelperMockRecorder) CompareHashAndPassword(a, b any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompareHashAndPassword", reflect.TypeOf((*MockHelper)(nil).CompareHashAndPassword), a, b)
}

// Copy mocks base method.
func (m *MockHelper) Copy(a *models.UserDetailsResponse, b *models.UserSignInResponse) (models.UserDetailsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Copy", a, b)
	ret0, _ := ret[0].(models.UserDetailsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Copy indicates an expected call of Copy.
func (mr *MockHelperMockRecorder) Copy(a, b any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*MockHelper)(nil).Copy), a, b)
}

// GenerateRefferalCode mocks base method.
func (m *MockHelper) GenerateRefferalCode() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRefferalCode")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateRefferalCode indicates an expected call of GenerateRefferalCode.
func (mr *MockHelperMockRecorder) GenerateRefferalCode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRefferalCode", reflect.TypeOf((*MockHelper)(nil).GenerateRefferalCode))
}

// GenerateTokenAdmin mocks base method.
func (m *MockHelper) GenerateTokenAdmin(admin models.AdminDetailsResponse) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTokenAdmin", admin)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateTokenAdmin indicates an expected call of GenerateTokenAdmin.
func (mr *MockHelperMockRecorder) GenerateTokenAdmin(admin any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTokenAdmin", reflect.TypeOf((*MockHelper)(nil).GenerateTokenAdmin), admin)
}

// GenerateTokenClients mocks base method.
func (m *MockHelper) GenerateTokenClients(user models.UserDetailsResponse) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTokenClients", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateTokenClients indicates an expected call of GenerateTokenClients.
func (mr *MockHelperMockRecorder) GenerateTokenClients(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTokenClients", reflect.TypeOf((*MockHelper)(nil).GenerateTokenClients), user)
}

// PasswordHashing mocks base method.
func (m *MockHelper) PasswordHashing(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PasswordHashing", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PasswordHashing indicates an expected call of PasswordHashing.
func (mr *MockHelperMockRecorder) PasswordHashing(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PasswordHashing", reflect.TypeOf((*MockHelper)(nil).PasswordHashing), arg0)
}

// TwilioSendOTP mocks base method.
func (m *MockHelper) TwilioSendOTP(phone, serviceID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TwilioSendOTP", phone, serviceID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TwilioSendOTP indicates an expected call of TwilioSendOTP.
func (mr *MockHelperMockRecorder) TwilioSendOTP(phone, serviceID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TwilioSendOTP", reflect.TypeOf((*MockHelper)(nil).TwilioSendOTP), phone, serviceID)
}

// TwilioSetup mocks base method.
func (m *MockHelper) TwilioSetup(username, password string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "TwilioSetup", username, password)
}

// TwilioSetup indicates an expected call of TwilioSetup.
func (mr *MockHelperMockRecorder) TwilioSetup(username, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TwilioSetup", reflect.TypeOf((*MockHelper)(nil).TwilioSetup), username, password)
}

// TwilioVerifyOTP mocks base method.
func (m *MockHelper) TwilioVerifyOTP(serviceID, code, phone string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TwilioVerifyOTP", serviceID, code, phone)
	ret0, _ := ret[0].(error)
	return ret0
}

// TwilioVerifyOTP indicates an expected call of TwilioVerifyOTP.
func (mr *MockHelperMockRecorder) TwilioVerifyOTP(serviceID, code, phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TwilioVerifyOTP", reflect.TypeOf((*MockHelper)(nil).TwilioVerifyOTP), serviceID, code, phone)
}

// ValidateAlphabets mocks base method.
func (m *MockHelper) ValidateAlphabets(data string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAlphabets", data)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateAlphabets indicates an expected call of ValidateAlphabets.
func (mr *MockHelperMockRecorder) ValidateAlphabets(data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAlphabets", reflect.TypeOf((*MockHelper)(nil).ValidateAlphabets), data)
}

// ValidateDatatype mocks base method.
func (m *MockHelper) ValidateDatatype(data, intOrString string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateDatatype", data, intOrString)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateDatatype indicates an expected call of ValidateDatatype.
func (mr *MockHelperMockRecorder) ValidateDatatype(data, intOrString any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateDatatype", reflect.TypeOf((*MockHelper)(nil).ValidateDatatype), data, intOrString)
}

// ValidatePhoneNumber mocks base method.
func (m *MockHelper) ValidatePhoneNumber(phone string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatePhoneNumber", phone)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ValidatePhoneNumber indicates an expected call of ValidatePhoneNumber.
func (mr *MockHelperMockRecorder) ValidatePhoneNumber(phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatePhoneNumber", reflect.TypeOf((*MockHelper)(nil).ValidatePhoneNumber), phone)
}

// ValidatePin mocks base method.
func (m *MockHelper) ValidatePin(pin string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatePin", pin)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ValidatePin indicates an expected call of ValidatePin.
func (mr *MockHelperMockRecorder) ValidatePin(pin any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatePin", reflect.TypeOf((*MockHelper)(nil).ValidatePin), pin)
}
