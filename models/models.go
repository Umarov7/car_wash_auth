package models

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}

type RegisterResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDetails struct {
	Id       string `json:"id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
}
