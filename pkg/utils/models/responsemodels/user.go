package responsemodel

type SignupData struct {
	ID            string `json:"userID"`
	Name          string `json:"name,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	OTP           string `json:"otp,omitempty"`
	Token         string `json:"token,omitempty"`
	IsUserExist   string `json:"isUserExist,omitempty"`
	ReferalCode   string `json:"referalCode"`
	WalletBelance uint   `json:"walletBelance"`
}
