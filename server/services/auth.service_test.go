package services

import (
	"errors"
	"testing"
	"time"

	"github.com/snowlynxsoftware/oto-api/server/database/repositories"
	"github.com/snowlynxsoftware/oto-api/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTokenService is a mock implementation of ITokenService
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateAccessToken(userID int) (*string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockTokenService) ValidateToken(token *string) (*int, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockTokenService) GenerateVerificationToken(userID int) (*string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockTokenService) GenerateLoginWithEmailToken(userID int) (*string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockTokenService) ValidateVerificationToken(token *string) (*int, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockTokenService) ValidateLoginWithEmailToken(token *string) (*int, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockTokenService) GenerateRefreshToken(userID int) (*string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

// MockCryptoService is a mock implementation of ICryptoService
type MockCryptoService struct {
	mock.Mock
}

func (m *MockCryptoService) HashPassword(password string) (*string, error) {
	args := m.Called(password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockCryptoService) ValidatePassword(password string, hashedPassword string) (bool, error) {
	args := m.Called(password, hashedPassword)
	return args.Bool(0), args.Error(1)
}

// MockEmailService is a mock implementation of IEmailService
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(options *EmailSendOptions) bool {
	args := m.Called(options)
	return args.Bool(0)
}

func (m *MockEmailService) GetTemplates() IEmailTemplates {
	args := m.Called()
	return args.Get(0).(IEmailTemplates)
}

// MockEmailTemplates is a mock implementation of IEmailTemplates
type MockEmailTemplates struct {
	mock.Mock
}

func (m *MockEmailTemplates) GetNewUserEmailTemplate(baseURL string, verificationToken string) string {
	args := m.Called(baseURL, verificationToken)
	return args.String(0)
}

func (m *MockEmailTemplates) GetLoginEmailTemplate(baseURL string, loginToken string) string {
	args := m.Called(baseURL, loginToken)
	return args.String(0)
}

// Test RegisterNewUser - Success
func TestAuthService_RegisterNewUser_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)
	mockEmailTemplates := new(MockEmailTemplates)

	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService)

	userDTO := &models.UserCreateDTO{
		Email:       "test@example.com",
		DisplayName: "Test User",
		Password:    "password123",
	}

	hashedPassword := "hashed_password_123"
	verificationToken := "verification_token_123"
	emailTemplate := "<html>Verification email</html>"

	expectedUser := &repositories.UserEntity{
		ID:          123,
		Email:       userDTO.Email,
		DisplayName: userDTO.DisplayName,
		CreatedAt:   time.Now(),
	}

	// Mock expectations
	mockUserRepo.On("GetUserByEmail", userDTO.Email).Return(nil, errors.New("user not found"))
	mockCryptoService.On("HashPassword", userDTO.Password).Return(&hashedPassword, nil)
	mockUserRepo.On("CreateNewUser", mock.MatchedBy(func(dto *models.UserCreateDTO) bool {
		return dto.Email == userDTO.Email && dto.Password == hashedPassword
	})).Return(expectedUser, nil)
	mockTokenService.On("GenerateVerificationToken", 123).Return(&verificationToken, nil)
	mockEmailService.On("GetTemplates").Return(mockEmailTemplates)
	mockEmailTemplates.On("GetNewUserEmailTemplate", "http://localhost:3000", verificationToken).Return(emailTemplate)
	mockEmailService.On("SendEmail", mock.MatchedBy(func(options *EmailSendOptions) bool {
		return options.ToEmail == userDTO.Email && options.Subject == "Open Trivia Online - Verify Your Account"
	})).Return(true)

	// Act
	result, err := authService.RegisterNewUser(userDTO)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.DisplayName, result.DisplayName)
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
	mockEmailTemplates.AssertExpectations(t)
}

// Test RegisterNewUser - User Already Exists
func TestAuthService_RegisterNewUser_UserAlreadyExists(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService)

	userDTO := &models.UserCreateDTO{
		Email:       "existing@example.com",
		DisplayName: "Test User",
		Password:    "password123",
	}

	existingUser := &repositories.UserEntity{
		ID:    123,
		Email: userDTO.Email,
	}

	mockUserRepo.On("GetUserByEmail", userDTO.Email).Return(existingUser, nil)

	// Act
	result, err := authService.RegisterNewUser(userDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "a user already exists with the specified email")
	mockUserRepo.AssertExpectations(t)
}

// Test RegisterNewUser - Password Hashing Error
func TestAuthService_RegisterNewUser_PasswordHashingError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService)

	userDTO := &models.UserCreateDTO{
		Email:       "test@example.com",
		DisplayName: "Test User",
		Password:    "password123",
	}

	mockUserRepo.On("GetUserByEmail", userDTO.Email).Return(nil, errors.New("user not found"))
	mockCryptoService.On("HashPassword", userDTO.Password).Return(nil, errors.New("hashing failed"))

	// Act
	result, err := authService.RegisterNewUser(userDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "hashing failed")
	mockUserRepo.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
}

// Test RegisterNewUser - Email Send Failure
func TestAuthService_RegisterNewUser_EmailSendFailure(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)
	mockEmailTemplates := new(MockEmailTemplates)

	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService)

	userDTO := &models.UserCreateDTO{
		Email:       "test@example.com",
		DisplayName: "Test User",
		Password:    "password123",
	}

	hashedPassword := "hashed_password_123"
	verificationToken := "verification_token_123"
	emailTemplate := "<html>Verification email</html>"

	expectedUser := &repositories.UserEntity{
		ID:          123,
		Email:       userDTO.Email,
		DisplayName: userDTO.DisplayName,
	}

	mockUserRepo.On("GetUserByEmail", userDTO.Email).Return(nil, errors.New("user not found"))
	mockCryptoService.On("HashPassword", userDTO.Password).Return(&hashedPassword, nil)
	mockUserRepo.On("CreateNewUser", mock.AnythingOfType("*models.UserCreateDTO")).Return(expectedUser, nil)
	mockTokenService.On("GenerateVerificationToken", 123).Return(&verificationToken, nil)
	mockEmailService.On("GetTemplates").Return(mockEmailTemplates)
	mockEmailTemplates.On("GetNewUserEmailTemplate", "http://localhost:3000", verificationToken).Return(emailTemplate)
	mockEmailService.On("SendEmail", mock.AnythingOfType("*services.EmailSendOptions")).Return(false)

	// Act
	result, err := authService.RegisterNewUser(userDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "the user was created but the verification email failed to send")
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
	mockEmailTemplates.AssertExpectations(t)
}

// Test SendLoginEmail - Success
func TestAuthService_SendLoginEmail_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)
	mockEmailTemplates := new(MockEmailTemplates)

	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService)

	email := "test@example.com"
	loginToken := "login_token_123"
	emailTemplate := "<html>Login email</html>"

	user := &repositories.UserEntity{
		ID:       123,
		Email:    email,
		IsBanned: false,
	}

	mockUserRepo.On("GetUserByEmail", email).Return(user, nil)
	mockTokenService.On("GenerateLoginWithEmailToken", 123).Return(&loginToken, nil)
	mockEmailService.On("GetTemplates").Return(mockEmailTemplates)
	mockEmailTemplates.On("GetLoginEmailTemplate", "http://localhost:3000", loginToken).Return(emailTemplate)
	mockEmailService.On("SendEmail", mock.MatchedBy(func(options *EmailSendOptions) bool {
		return options.ToEmail == email && options.Subject == "Open Trivia Online - Login Email"
	})).Return(true)

	// Act
	result, err := authService.SendLoginEmail(email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Email, result.Email)
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
	mockEmailTemplates.AssertExpectations(t)
}

// Test SendLoginEmail - User Banned
func TestAuthService_SendLoginEmail_UserBanned(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService)

	email := "banned@example.com"
	user := &repositories.UserEntity{
		ID:       123,
		Email:    email,
		IsBanned: true,
	}

	mockUserRepo.On("GetUserByEmail", email).Return(user, nil)

	// Act
	result, err := authService.SendLoginEmail(email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user is banned")
	mockUserRepo.AssertExpectations(t)
}

// Test LoginWithEmailLink - Success
func TestAuthService_LoginWithEmailLink_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService)

	userId := 123
	accessToken := "access_token_123"

	mockTokenService.On("GenerateAccessToken", userId).Return(&accessToken, nil)
	mockUserRepo.On("UpdateUserLastLogin", &userId).Return(true, nil)

	// Act
	result, err := authService.LoginWithEmailLink(&userId)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, accessToken, result.AccessToken)
	assert.Equal(t, "", result.RefreshToken)
	mockTokenService.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test LoginWithEmailLink - Token Generation Error
func TestAuthService_LoginWithEmailLink_TokenGenerationError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService)

	userId := 123

	mockTokenService.On("GenerateAccessToken", userId).Return(nil, errors.New("token generation failed"))

	// Act
	result, err := authService.LoginWithEmailLink(&userId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "there was an issue trying to log this user in")
	mockTokenService.AssertExpectations(t)
}

// Test LoginWithEmailLink - Last Login Update Error
func TestAuthService_LoginWithEmailLink_LastLoginUpdateError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService)

	userId := 123
	accessToken := "access_token_123"

	mockTokenService.On("GenerateAccessToken", userId).Return(&accessToken, nil)
	mockUserRepo.On("UpdateUserLastLogin", &userId).Return(false, errors.New("database error"))

	// Act
	result, err := authService.LoginWithEmailLink(&userId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "there was an issue trying to log this user in")
	mockTokenService.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
