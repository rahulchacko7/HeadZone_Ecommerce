package helper

import (
	"HeadZone/pkg/helper/interfaces"
	"HeadZone/pkg/utils/models"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"github.com/twilio/twilio-go"

	cfg "HeadZone/pkg/config"

	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
	"golang.org/x/crypto/bcrypt"
)

type helper struct {
	cfg cfg.Config
}

func NewHelper(config cfg.Config) interfaces.Helper {
	return &helper{
		cfg: config,
	}
}

var client *twilio.RestClient

type AuthCustomClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func (helper *helper) GenerateTokenAdmin(admin models.AdminDetailsResponse) (string, string, error) {
	accessTokenClaims := &AuthCustomClaims{
		Id:    admin.ID,
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshTokenClaims := &AuthCustomClaims{
		Id:    admin.ID,
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(helper.cfg.ACCESS_KEY_ADMIN))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(helper.cfg.REFRESH_KEY_ADMIN))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (h *helper) TwilioSetup(username string, password string) {
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: username,
		Password: password,
	})

}

func (h *helper) TwilioSendOTP(phone string, serviceID string) (string, error) {
	to := "+91" + phone
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(serviceID, params)
	if err != nil {

		return " ", err
	}

	return *resp.Sid, nil

}

func (h *helper) TwilioVerifyOTP(serviceID string, code string, phone string) error {

	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo("+91" + phone)
	params.SetCode(code)
	resp, err := client.VerifyV2.CreateVerificationCheck(serviceID, params)

	if err != nil {
		return err
	}

	if *resp.Status == "approved" {
		return nil
	}

	return errors.New("failed to validate otp")

}

func (h *helper) GenerateTokenClients(user models.UserDetailsResponse) (string, error) {
	claims := &AuthCustomClaims{
		Id:    user.Id,
		Email: user.Email,
		Role:  "client",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.cfg.ACCESS_KEY_USER))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *helper) GenerateRefferalCode() (string, error) {
	// Calculate the required number of random bytes
	byteLength := (5 * 5) / 8

	// Generate a random byte array
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to base32
	encoder := base32.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567").WithPadding(base32.NoPadding)
	encoded := encoder.EncodeToString(randomBytes)

	// Trim any additional characters to match the desired length
	encoded = encoded[:5]

	return encoded, nil
}

func (h *helper) PasswordHashing(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New("internal server error")
	}

	hash := string(hashedPassword)
	return hash, nil
}

func (h *helper) CompareHashAndPassword(a string, b string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a), []byte(b))
	if err != nil {
		return err
	}
	return nil
}

func (h *helper) Copy(a *models.UserDetailsResponse, b *models.UserSignInResponse) (models.UserDetailsResponse, error) {
	err := copier.Copy(a, b)
	if err != nil {
		return models.UserDetailsResponse{}, err
	}

	return *a, nil
}

func ConvertToExel(sales []models.OrderDetailsAdmin) (*excelize.File, error) {

	filename := "salesReport/sales_report.xlsx"
	file := excelize.NewFile()

	file.SetCellValue("Sheet1", "A1", "Item")
	file.SetCellValue("Sheet1", "B1", "Total Amount Sold")

	for i, sale := range sales {
		col1 := fmt.Sprintf("A%d", i+1)
		col2 := fmt.Sprintf("B%d", i+1)

		file.SetCellValue("Sheet1", col1, sale.ProductName)
		file.SetCellValue("Sheet1", col2, sale.TotalAmount)

	}

	if err := file.SaveAs(filename); err != nil {
		return nil, err
	}

	return file, nil
}

func (h *helper) GetTimeFromPeriod(timePeriod string) (time.Time, time.Time) {

	endDate := time.Now()

	if timePeriod == "week" {
		startDate := endDate.AddDate(0, 0, -6)
		return startDate, endDate
	}

	if timePeriod == "month" {
		startDate := endDate.AddDate(0, -1, 0)
		return startDate, endDate
	}

	if timePeriod == "year" {
		startDate := endDate.AddDate(0, -1, 0)
		return startDate, endDate
	}

	return endDate.AddDate(0, 0, -6), endDate

}

func (h *helper) ValidatePhoneNumber(phone string) bool {
	phoneNumber := phone
	pattern := `^\d{10}$`
	regex := regexp.MustCompile(pattern)
	value := regex.MatchString(phoneNumber)
	return value
}

func (h *helper) ValidatePin(pin string) bool {

	match, _ := regexp.MatchString(`^\d{4}(\d{2})?$`, pin)
	return match

}

func (h *helper) ValidateDatatype(data, intOrString string) (bool, error) {

	switch intOrString {
	case "int":
		if _, err := strconv.Atoi(data); err != nil {
			return false, errors.New("data is not an integer")
		}
		return true, nil
	case "string":
		kind := reflect.TypeOf(data).Kind()
		return kind == reflect.String, nil
	default:
		return false, errors.New("data is not" + intOrString)
	}

}

func (h *helper) ValidateAlphabets(data string) (bool, error) {
	for _, char := range data {
		if !unicode.IsLetter(char) {
			return false, errors.New("data contains non-alphabetic characters")
		}
	}
	return true, nil
}
