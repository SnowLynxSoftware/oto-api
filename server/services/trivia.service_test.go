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

// MockTriviaRepository is a mock implementation of ITriviaRepository
type MockTriviaRepository struct {
	mock.Mock
}

func (m *MockTriviaRepository) GetTriviaQuestionByText(question string) (*repositories.TriviaQuestionEntity, error) {
	args := m.Called(question)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.TriviaQuestionEntity), args.Error(1)
}

func (m *MockTriviaRepository) ImportTriviaQuestions(data []models.TriviaQuestionImportData) (*models.TriviaQuestionImportResults, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TriviaQuestionImportResults), args.Error(1)
}

func (m *MockTriviaRepository) GetWrongAnswerByText(answer string) (*repositories.WrongAnswerPoolEntity, error) {
	args := m.Called(answer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.WrongAnswerPoolEntity), args.Error(1)
}

func (m *MockTriviaRepository) ImportWrongAnswers(data []models.TriviaWrongAnswerImportData) (*models.TriviaWrongAnswerImportResults, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TriviaWrongAnswerImportResults), args.Error(1)
}

func (m *MockTriviaRepository) GetTriviaDeckById(deckId int64) (*repositories.TriviaDeckEntity, error) {
	args := m.Called(deckId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.TriviaDeckEntity), args.Error(1)
}

func (m *MockTriviaRepository) CreateNewTriviaDeck(name string, description string, isSystemDeck bool) (*repositories.TriviaDeckEntity, error) {
	args := m.Called(name, description, isSystemDeck)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.TriviaDeckEntity), args.Error(1)
}

func (m *MockTriviaRepository) UpdateTriviaDeckMetadata(deckId int64, name string, description string) (*repositories.TriviaDeckEntity, error) {
	args := m.Called(deckId, name, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.TriviaDeckEntity), args.Error(1)
}

func (m *MockTriviaRepository) UpdateTriviaDeckApprovalStatus(deckId int64, isApproved bool) (*repositories.TriviaDeckEntity, error) {
	args := m.Called(deckId, isApproved)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.TriviaDeckEntity), args.Error(1)
}

func (m *MockTriviaRepository) UpdateTriviaDeckArchivalStatus(deckId int64, isArchived bool) (*repositories.TriviaDeckEntity, error) {
	args := m.Called(deckId, isArchived)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.TriviaDeckEntity), args.Error(1)
}

// Test ImportTriviaQuestions - Success
func TestTriviaService_ImportTriviaQuestions_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	importData := []models.TriviaQuestionImportData{
		{
			Question:      "What is the capital of France?",
			CorrectAnswer: "Paris",
			Tags:          []string{"geography", "europe"},
		},
		{
			Question:      "What is 2 + 2?",
			CorrectAnswer: "4",
			Tags:          []string{"math", "basic"},
		},
	}

	expectedResults := &models.TriviaQuestionImportResults{
		TotalQuestionsProcessed: 2,
		QuestionsAdded:          2,
	}

	mockTriviaRepo.On("ImportTriviaQuestions", importData).Return(expectedResults, nil)

	// Act
	result, err := triviaService.ImportTriviaQuestions(importData)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResults.TotalQuestionsProcessed, result.TotalQuestionsProcessed)
	assert.Equal(t, expectedResults.QuestionsAdded, result.QuestionsAdded)
	mockTriviaRepo.AssertExpectations(t)
}

// Test ImportTriviaQuestions - Repository Error
func TestTriviaService_ImportTriviaQuestions_RepositoryError(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	importData := []models.TriviaQuestionImportData{
		{
			Question:      "Test question?",
			CorrectAnswer: "Test answer",
			Tags:          []string{"test"},
		},
	}

	mockTriviaRepo.On("ImportTriviaQuestions", importData).Return(nil, errors.New("database error"))

	// Act
	result, err := triviaService.ImportTriviaQuestions(importData)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	mockTriviaRepo.AssertExpectations(t)
}

// Test ImportWrongAnswers - Success
func TestTriviaService_ImportWrongAnswers_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	importData := []models.TriviaWrongAnswerImportData{
		{
			Question:   "What is the capital of France?",
			AnswerText: "London",
			Tags:       []string{"geography"},
		},
		{
			Question:   "What is the capital of France?",
			AnswerText: "Berlin",
			Tags:       []string{"geography"},
		},
	}

	expectedResults := &models.TriviaWrongAnswerImportResults{
		TotalAnswersProcessed: 2,
		AnswersAdded:          2,
	}

	mockTriviaRepo.On("ImportWrongAnswers", importData).Return(expectedResults, nil)

	// Act
	result, err := triviaService.ImportWrongAnswers(importData)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResults.TotalAnswersProcessed, result.TotalAnswersProcessed)
	assert.Equal(t, expectedResults.AnswersAdded, result.AnswersAdded)
	mockTriviaRepo.AssertExpectations(t)
}

// Test ImportWrongAnswers - Repository Error
func TestTriviaService_ImportWrongAnswers_RepositoryError(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	importData := []models.TriviaWrongAnswerImportData{
		{
			Question:   "Test question?",
			AnswerText: "Wrong answer",
			Tags:       []string{"test"},
		},
	}

	mockTriviaRepo.On("ImportWrongAnswers", importData).Return(nil, errors.New("import failed"))

	// Act
	result, err := triviaService.ImportWrongAnswers(importData)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "import failed")
	mockTriviaRepo.AssertExpectations(t)
}

// Test CreateNewTriviaDeck - Success
func TestTriviaService_CreateNewTriviaDeck_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	name := "General Knowledge"
	description := "A deck of general knowledge questions"
	isSystemDeck := true

	expectedDeck := &repositories.TriviaDeckEntity{
		ID:           123,
		Name:         name,
		Description:  &description,
		IsSystemDeck: isSystemDeck,
		IsApproved:   false,
		IsArchived:   false,
		CreatedAt:    time.Now(),
	}

	mockTriviaRepo.On("CreateNewTriviaDeck", name, description, isSystemDeck).Return(expectedDeck, nil)

	// Act
	result, err := triviaService.CreateNewTriviaDeck(name, description, isSystemDeck)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedDeck.Name, result.Name)
	assert.Equal(t, expectedDeck.Description, result.Description)
	assert.Equal(t, expectedDeck.IsSystemDeck, result.IsSystemDeck)
	mockTriviaRepo.AssertExpectations(t)
}

// Test CreateNewTriviaDeck - Repository Error
func TestTriviaService_CreateNewTriviaDeck_RepositoryError(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	name := "Test Deck"
	description := "Test description"
	isSystemDeck := false

	mockTriviaRepo.On("CreateNewTriviaDeck", name, description, isSystemDeck).Return(nil, errors.New("creation failed"))

	// Act
	result, err := triviaService.CreateNewTriviaDeck(name, description, isSystemDeck)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "creation failed")
	mockTriviaRepo.AssertExpectations(t)
}

// Test GetTriviaDeckById - Success
func TestTriviaService_GetTriviaDeckById_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	deckId := int64(123)
	description := "A test deck"
	expectedDeck := &repositories.TriviaDeckEntity{
		ID:          deckId,
		Name:        "Test Deck",
		Description: &description,
		CreatedAt:   time.Now(),
	}

	mockTriviaRepo.On("GetTriviaDeckById", deckId).Return(expectedDeck, nil)

	// Act
	result, err := triviaService.GetTriviaDeckById(deckId)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedDeck.ID, result.ID)
	assert.Equal(t, expectedDeck.Name, result.Name)
	mockTriviaRepo.AssertExpectations(t)
}

// Test GetTriviaDeckById - Deck Not Found
func TestTriviaService_GetTriviaDeckById_DeckNotFound(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	deckId := int64(999)
	mockTriviaRepo.On("GetTriviaDeckById", deckId).Return(nil, errors.New("deck not found"))

	// Act
	result, err := triviaService.GetTriviaDeckById(deckId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "deck not found")
	mockTriviaRepo.AssertExpectations(t)
}
