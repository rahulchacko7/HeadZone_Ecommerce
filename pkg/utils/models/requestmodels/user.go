package requestmodel

type UserDetails struct {
	Id              string `json:"id"`
	Name            string `json:"name"           validate:"required"`
	Email           string `json:"email"          validate:"email"`
	Phone           string `json:"phone"          validate:"len=10"`
	Password        string `json:"password,omitempty"       validate:"min=4"`
	ConfirmPassword string `json:"confirmpassword,omitempty" validate:"eqfield=Password"`
	ReferalCode     string `json:"referalCode,omitempty"`
}