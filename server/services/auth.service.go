package services

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/snowlynxsoftware/oto-api/server/database/repositories"
	"github.com/snowlynxsoftware/oto-api/server/models"
	"github.com/snowlynxsoftware/oto-api/server/util"
)

type IAuthService interface {
	RegisterNewUser(dto *models.UserCreateDTO) (*repositories.UserEntity, error)
	Login(authHeaderStr *string) (*models.UserLoginResponseDTO, error)
	VerifyNewUser(verificationToken *string) (*int, error)
	SendLoginEmail(email string) (*repositories.UserEntity, error)
	LoginWithEmailLink(userId *int) (*models.UserLoginResponseDTO, error)
	UpdateUserPassword(userId *int, password string) (*int, error)
}

type AuthService struct {
	userRepository repositories.IUserRepository
	tokenService   ITokenService
	cryptoService  ICryptoService
	emailService   IEmailService
}

func NewAuthService(
	userRepository repositories.IUserRepository,
	tokenService ITokenService,
	cryptoService ICryptoService,
	emailService IEmailService,
) IAuthService {
	return &AuthService{userRepository: userRepository, tokenService: tokenService, cryptoService: cryptoService, emailService: emailService}
}

func (s *AuthService) RegisterNewUser(dto *models.UserCreateDTO) (*repositories.UserEntity, error) {

	var _, err = s.userRepository.GetUserByEmail(dto.Email)
	if err == nil {
		return nil, errors.New("a user already exists with the specified email")
	}

	hashedPassword, err := s.cryptoService.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	dto.Password = *hashedPassword

	newUser, err := s.userRepository.CreateNewUser(dto)
	if err != nil {
		return nil, err
	}

	verificationToken, err := s.tokenService.GenerateVerificationToken(int(newUser.ID))
	if err != nil {
		return nil, err
	}

	var emailOptions = &EmailSendOptions{}
	emailOptions.FromEmail = "do-not-reply@opentriviaonline.com"
	emailOptions.ToEmail = newUser.Email
	emailOptions.Subject = "Open Trivia Online - Verify Your Account"
	// TODO: Update this to use the correct URL
	emailOptions.HTMLContent = s.emailService.GetTemplates().GetNewUserEmailTemplate("http://localhost:3000", *verificationToken)
	var isEmailSuccess = s.emailService.SendEmail(emailOptions)
	if isEmailSuccess {
		return newUser, nil
	} else {
		return nil, errors.New("the user was created but the verification email failed to send")
	}
}

func (s *AuthService) SendLoginEmail(email string) (*repositories.UserEntity, error) {

	var user, err = s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user.IsBanned {
		return nil, errors.New("user is banned")
	}

	verificationToken, err := s.tokenService.GenerateLoginWithEmailToken(int(user.ID))
	if err != nil {
		return nil, err
	}

	var emailOptions = &EmailSendOptions{}
	emailOptions.FromEmail = "do-not-reply@opentriviaonline.com"
	emailOptions.ToEmail = user.Email
	emailOptions.Subject = "Open Trivia Online - Login Email"
	// TODO: Update this to use the correct URL
	emailOptions.HTMLContent = s.emailService.GetTemplates().GetLoginEmailTemplate("http://localhost:3000", *verificationToken)
	var isEmailSuccess = s.emailService.SendEmail(emailOptions)
	if isEmailSuccess {
		return user, nil
	} else {
		return nil, errors.New("the login by email failed to send")
	}
}

func (s *AuthService) LoginWithEmailLink(userId *int) (*models.UserLoginResponseDTO, error) {

	accessToken, err := s.tokenService.GenerateAccessToken(*userId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return nil, errors.New("there was an issue trying to log this user in")
	}

	// Update user's last login timestamp
	_, err = s.userRepository.UpdateUserLastLogin(userId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return nil, errors.New("there was an issue trying to log this user in")
	}

	return &models.UserLoginResponseDTO{
		AccessToken:  *accessToken,
		RefreshToken: "",
	}, nil
}

func (s *AuthService) VerifyNewUser(verificationToken *string) (*int, error) {

	var userId, err = s.tokenService.ValidateToken(verificationToken)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return nil, errors.New("the token could not be verified")
	}

	_, err = s.userRepository.MarkUserVerified(userId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return nil, err
	}

	return userId, nil
}

func (s *AuthService) UpdateUserPassword(userId *int, password string) (*int, error) {

	hashedPassword, err := s.cryptoService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	_, err = s.userRepository.UpdateUserPassword(userId, *hashedPassword)
	if err != nil {
		return nil, err
	}

	return userId, nil
}

func (s *AuthService) Login(authHeaderStr *string) (*models.UserLoginResponseDTO, error) {

	encodedCredentials := strings.TrimPrefix(*authHeaderStr, "Basic ")
	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return nil, errors.New("failed to decode authorization header")
	}

	credentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(credentials) != 2 {
		return nil, errors.New("invalid authorization header format")
	}

	email := credentials[0]
	password := credentials[1]

	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("there was an issue trying to log this user in")
	}

	isValid, err := s.cryptoService.ValidatePassword(password, *user.PasswordHash)
	if err != nil || !isValid {
		return nil, errors.New("there was an issue trying to log this user in")
	}

	accessToken, err := s.tokenService.GenerateAccessToken(int(user.ID))
	if err != nil {
		return nil, errors.New("there was an issue trying to log this user in")
	}

	// Update user's last login timestamp
	userId := int(user.ID)
	_, err = s.userRepository.UpdateUserLastLogin(&userId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return nil, errors.New("there was an issue trying to log this user in")
	}

	return &models.UserLoginResponseDTO{
		AccessToken:  *accessToken,
		RefreshToken: "",
	}, nil
}
