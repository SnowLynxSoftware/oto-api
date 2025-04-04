package models

type UserCreateDTO struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

type UserLoginResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserUpdatePasswordDTO struct {
	Password string `json:"password"`
}
