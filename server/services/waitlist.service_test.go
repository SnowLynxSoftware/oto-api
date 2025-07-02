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

// MockWaitlistRepository is a mock implementation of IWaitlistRepository
type MockWaitlistRepository struct {
	mock.Mock
}

func (m *MockWaitlistRepository) GetWaitlistEntryById(id int) (*repositories.WaitlistEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.WaitlistEntity), args.Error(1)
}

func (m *MockWaitlistRepository) CreateNewWaitlistEntry(dto *models.WaitlistCreateDTO) (*repositories.WaitlistEntity, error) {
	args := m.Called(dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.WaitlistEntity), args.Error(1)
}

// Test CreateNewWaitlistEntry - Success
func TestWaitlistService_CreateNewWaitlistEntry_Success(t *testing.T) {
	// Arrange
	mockWaitlistRepo := new(MockWaitlistRepository)
	waitlistService := NewWaitlistService(mockWaitlistRepo)

	waitlistDTO := &models.WaitlistCreateDTO{
		Email: "test@example.com",
	}

	expectedEntry := &repositories.WaitlistEntity{
		ID:        123,
		Email:     waitlistDTO.Email,
		CreatedAt: time.Now(),
	}

	// Mock expectations
	mockWaitlistRepo.On("CreateNewWaitlistEntry", waitlistDTO).Return(expectedEntry, nil)

	// Act
	result, err := waitlistService.CreateNewWaitlistEntry(waitlistDTO)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedEntry.Email, result.Email)
	assert.Equal(t, expectedEntry.ID, result.ID)
	mockWaitlistRepo.AssertExpectations(t)
}

// Test CreateNewWaitlistEntry - Repository Error
func TestWaitlistService_CreateNewWaitlistEntry_RepositoryError(t *testing.T) {
	// Arrange
	mockWaitlistRepo := new(MockWaitlistRepository)
	waitlistService := NewWaitlistService(mockWaitlistRepo)

	waitlistDTO := &models.WaitlistCreateDTO{
		Email: "test@example.com",
	}

	// Mock expectations
	mockWaitlistRepo.On("CreateNewWaitlistEntry", waitlistDTO).Return(nil, errors.New("database error"))

	// Act
	result, err := waitlistService.CreateNewWaitlistEntry(waitlistDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	mockWaitlistRepo.AssertExpectations(t)
}

// Test CreateNewWaitlistEntry - Empty Email
func TestWaitlistService_CreateNewWaitlistEntry_EmptyEmail(t *testing.T) {
	// Arrange
	mockWaitlistRepo := new(MockWaitlistRepository)
	waitlistService := NewWaitlistService(mockWaitlistRepo)

	waitlistDTO := &models.WaitlistCreateDTO{
		Email: "", // Empty email
	}

	// Act
	result, err := waitlistService.CreateNewWaitlistEntry(waitlistDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "email cannot be empty")
	// No repository calls should be made
	mockWaitlistRepo.AssertExpectations(t)
}

// Test CreateNewWaitlistEntry - Email Too Short
func TestWaitlistService_CreateNewWaitlistEntry_EmailTooShort(t *testing.T) {
	// Arrange
	mockWaitlistRepo := new(MockWaitlistRepository)
	waitlistService := NewWaitlistService(mockWaitlistRepo)

	waitlistDTO := &models.WaitlistCreateDTO{
		Email: "a@b", // Too short email
	}

	// Act
	result, err := waitlistService.CreateNewWaitlistEntry(waitlistDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "email must be at least 5 characters long")
	// No repository calls should be made
	mockWaitlistRepo.AssertExpectations(t)
}
