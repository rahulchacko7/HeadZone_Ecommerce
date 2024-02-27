package handler

import (
	"HeadZone/pkg/utils/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_interfaces "HeadZone/pkg/usecase/mock"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_UserSignUp(t *testing.T) {
	testcase := map[string]struct {
		input         models.UserDetails
		buildstub     func(useCaseMock *mock_interfaces.MockUserUseCase, sighupData models.UserDetails)
		checkResponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		"valid Signup": {
			input: models.UserDetails{
				Name:            "rahul",
				Email:           "rahul@gmail.com",
				Phone:           "7012483965",
				Password:        "54321",
				ConfirmPassword: "54321",
			},
			buildstub: func(useCaseMock *mock_interfaces.MockUserUseCase, signupData models.UserDetails) {
				err := validator.New().Struct(signupData)
				if err != nil {
					fmt.Println("validation failed")
				}
				useCaseMock.EXPECT().UserSignUp(signupData).Times(1).Return(models.TokenUsers{
					Users: models.UserDetailsResponse{
						Id:    1,
						Name:  "rahul",
						Email: "rahul@gmail.com",
						Phone: "7012483965",
					},
					Token: "aksjgnal.fiugliagbldfgbldf.gdbladfjnb",
				}, nil)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, responseRecorder.Code)
			},
		},
		"user could not sign up": {
			input: models.UserDetails{
				Name:            "rahul",
				Email:           "rahul@gmail.com",
				Phone:           "7012483965",
				Password:        "54321",
				ConfirmPassword: "54321",
			},
			buildstub: func(useCaseMock *mock_interfaces.MockUserUseCase, signupData models.UserDetails) {
				err := validator.New().Struct(signupData)
				if err != nil {
					fmt.Println("vallidation failed")
				}
				useCaseMock.EXPECT().UserSignUp(signupData).Times(1).Return(models.TokenUsers{}, errors.New("cannot signup"))
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
	}
	for testName, test := range testcase {
		test := test
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockUseCase := mock_interfaces.NewMockUserUseCase(ctrl)
			test.buildstub(mockUseCase, test.input)
			UserHandler := NewUserHandler(mockUseCase)

			server := gin.Default()
			server.POST("/signup", UserHandler.UserSignUp)

			jsonData, err := json.Marshal(test.input)
			assert.NoError(t, err)
			body := bytes.NewBuffer(jsonData)

			mockRequest, err := http.NewRequest(http.MethodPost, "/signup", body)
			assert.NoError(t, err)
			responseRecorder := httptest.NewRecorder()
			server.ServeHTTP(responseRecorder, mockRequest)
			test.checkResponse(t, responseRecorder)
		})
	}
}

func Test_LoginHandler(t *testing.T) {
	testCase := map[string]struct {
		input         models.UserLogin
		buildstub     func(userCaseMock *mock_interfaces.MockUserUseCase, login models.UserLogin)
		checkResponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		"sucess": {
			input: models.UserLogin{
				Email:    "rahulchacko@gmail.com",
				Password: "54321",
			},
			buildstub: func(userCaseMock *mock_interfaces.MockUserUseCase, login models.UserLogin) {
				err := validator.New().Struct(login)
				if err != nil {
					fmt.Println("validation failed")
				}
				userCaseMock.EXPECT().LoginHandler(login).Times(1).Return(models.TokenUsers{
					Users: models.UserDetailsResponse{
						Id:    1,
						Name:  "rahul",
						Email: "rahulchacko7@gmail.com",
						Phone: "7012394965",
					},
					Token: "hdiuwhedui8wqey89.dikosaios",
				}, nil)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, responseRecorder.Code)
			},
		},
		"user couldn't login": {
			input: models.UserLogin{
				Email:    "rahulchacko7@gmail.com",
				Password: "54321",
			},
			buildstub: func(userCaseMock *mock_interfaces.MockUserUseCase, login models.UserLogin) {
				err := validator.New().Struct(login)
				if err != nil {
					fmt.Println("validation failed")
				}
				userCaseMock.EXPECT().LoginHandler(login).Times(1).Return(models.TokenUsers{}, errors.New("cannot login up"))

			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
	}

	for testName, test := range testCase {
		test := test
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockUseCase := mock_interfaces.NewMockUserUseCase(ctrl)
			test.buildstub(mockUseCase, test.input)
			UserHandler := NewUserHandler(mockUseCase)
			server := gin.Default()
			server.POST("/login", UserHandler.LoginHandler)
			jsonData, err := json.Marshal(test.input)
			assert.NoError(t, err)
			body := bytes.NewBuffer(jsonData)
			mockRequest, err := http.NewRequest(http.MethodPost, "/login", body)
			assert.NoError(t, err)
			responseRecorder := httptest.NewRecorder()
			server.ServeHTTP(responseRecorder, mockRequest)
			test.checkResponse(t, responseRecorder)
		})
	}
}
