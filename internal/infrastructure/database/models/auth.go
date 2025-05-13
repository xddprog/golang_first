package models

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}


type AuthResponseModel struct {
	TokenPair
	User BaseUserModel `json:"user"`
}


type LoginUserModel struct {
	Email 	 string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"email" validate:"required"`
}


type RegisterUserModel struct {
	UserModel
	Password string `json:"password" validate:"required,min=8"`
}
