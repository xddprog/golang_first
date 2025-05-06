package models


type UserModel struct {
	Id int `json:"id" validate:"required"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email string `json:"email" validate:"required"`
	Password string `json:"-"`
}


type CreateUserModel struct {
	UserModel
    Password string `json:"password" validate:"required,min=8"`
}