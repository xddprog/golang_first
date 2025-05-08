package models


type BaseUserModel struct {
	Id       int    `json:"id" validate:"required"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required"`
}


type UserModel struct {
	BaseUserModel
	Password string `json:"-"`
}


type UpdateUserModel struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required"`
}