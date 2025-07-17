package services

import (
	"errors"
	"testing"
	"time"

	"github.com/lib/pq"
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

// New CRUD methods for questions
func (m *MockTriviaRepository) GetQuestionsCount(searchString, statusFilter, tagFilter string) (*int, error) {
	args := m.Called(searchString, statusFilter, tagFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockTriviaRepository) GetQuestions(limit, offset int, searchString, statusFilter, tagFilter string) ([]*repositories.TriviaQuestionEntity, error) {
	args := m.Called(limit, offset, searchString, statusFilter, tagFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.TriviaQuestionEntity), args.Error(1)
}

func (m *MockTriviaRepository) GetQuestionById(id int64) (*repositories.TriviaQuestionEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.TriviaQuestionEntity), args.Error(1)
}

func (m *MockTriviaRepository) CreateQuestion(dto *models.TriviaQuestionCreateDTO) (*repositories.TriviaQuestionEntity, error) {
	args := m.Called(dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.TriviaQuestionEntity), args.Error(1)
}

func (m *MockTriviaRepository) UpdateQuestion(dto *models.TriviaQuestionUpdateDTO, id int64) (*repositories.TriviaQuestionEntity, error) {
	args := m.Called(dto, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.TriviaQuestionEntity), args.Error(1)
}

func (m *MockTriviaRepository) ToggleQuestionArchived(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTriviaRepository) ToggleQuestionPublished(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// New CRUD methods for wrong answers
func (m *MockTriviaRepository) GetWrongAnswersCount(searchString, statusFilter, tagFilter string) (*int, error) {
	args := m.Called(searchString, statusFilter, tagFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockTriviaRepository) GetWrongAnswers(limit, offset int, searchString, statusFilter, tagFilter string) ([]*repositories.WrongAnswerPoolEntity, error) {
	args := m.Called(limit, offset, searchString, statusFilter, tagFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.WrongAnswerPoolEntity), args.Error(1)
}

func (m *MockTriviaRepository) GetWrongAnswerById(id int64) (*repositories.WrongAnswerPoolEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.WrongAnswerPoolEntity), args.Error(1)
}

func (m *MockTriviaRepository) CreateWrongAnswer(dto *models.WrongAnswerCreateDTO) (*repositories.WrongAnswerPoolEntity, error) {
	args := m.Called(dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.WrongAnswerPoolEntity), args.Error(1)
}

func (m *MockTriviaRepository) UpdateWrongAnswer(dto *models.WrongAnswerUpdateDTO, id int64) (*repositories.WrongAnswerPoolEntity, error) {
	args := m.Called(dto, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.WrongAnswerPoolEntity), args.Error(1)
}

func (m *MockTriviaRepository) ToggleWrongAnswerArchived(id int64) error {
	args := m.Called(id)
	return args.Error(0)
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

// ===========================================
// TRIVIA QUESTION CRUD TESTS
// ===========================================

// Test GetQuestions - Success with pagination
func TestTriviaService_GetQuestions_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	count := 150
	expectedQuestions := []*repositories.TriviaQuestionEntity{
		{
			ID:            1,
			Question:      "What is the capital of France?",
			CorrectAnswer: "Paris",
			Tags:          pq.StringArray{"geography", "europe"},
			IsPublished:   true,
			IsArchived:    false,
			CreatedAt:     time.Now(),
			ModifiedAt:    nil,
		},
		{
			ID:            2,
			Question:      "What is 2+2?",
			CorrectAnswer: "4",
			Tags:          pq.StringArray{"math", "basic"},
			IsPublished:   true,
			IsArchived:    false,
			CreatedAt:     time.Now(),
			ModifiedAt:    nil,
		},
	}

	mockTriviaRepo.On("GetQuestionsCount", "", "", "").Return(&count, nil)
	mockTriviaRepo.On("GetQuestions", 25, 0, "", "", "").Return(expectedQuestions, nil)

	// Act
	result, err := triviaService.GetQuestions(25, 0, "", "", "")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 150, result.Total)
	assert.Equal(t, 25, result.PageSize)
	assert.Equal(t, 1, result.Page)
	assert.Len(t, result.Results, 2)
	mockTriviaRepo.AssertExpectations(t)
}

// Test GetQuestions - Success with search and filters
func TestTriviaService_GetQuestions_WithFilters_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	count := 10
	expectedQuestions := []*repositories.TriviaQuestionEntity{
		{
			ID:            1,
			Question:      "What is the capital of France?",
			CorrectAnswer: "Paris",
			Tags:          pq.StringArray{"geography", "europe"},
			IsPublished:   true,
			IsArchived:    false,
			CreatedAt:     time.Now(),
			ModifiedAt:    nil,
		},
	}

	mockTriviaRepo.On("GetQuestionsCount", "capital", "published", "geography").Return(&count, nil)
	mockTriviaRepo.On("GetQuestions", 10, 0, "capital", "published", "geography").Return(expectedQuestions, nil)

	// Act
	result, err := triviaService.GetQuestions(10, 0, "capital", "published", "geography")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.Total)
	assert.Equal(t, 10, result.PageSize)
	assert.Len(t, result.Results, 1)
	mockTriviaRepo.AssertExpectations(t)
}

// Test GetQuestions - Repository error
func TestTriviaService_GetQuestions_RepositoryError(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	mockTriviaRepo.On("GetQuestions", 25, 0, "", "", "").Return(nil, errors.New("database error"))

	// Act
	result, err := triviaService.GetQuestions(25, 0, "", "", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	mockTriviaRepo.AssertExpectations(t)
}

// Test GetQuestionById - Success
func TestTriviaService_GetQuestionById_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	expectedQuestion := &repositories.TriviaQuestionEntity{
		ID:            1,
		Question:      "What is the capital of France?",
		CorrectAnswer: "Paris",
		Tags:          pq.StringArray{"geography", "europe"},
		IsPublished:   true,
		IsArchived:    false,
		CreatedAt:     time.Now(),
		ModifiedAt:    nil,
	}

	mockTriviaRepo.On("GetQuestionById", int64(1)).Return(expectedQuestion, nil)

	// Act
	result, err := triviaService.GetQuestionById(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "What is the capital of France?", result.Question)
	assert.Equal(t, "Paris", result.CorrectAnswer)
	assert.Equal(t, pq.StringArray{"geography", "europe"}, result.Tags)
	mockTriviaRepo.AssertExpectations(t)
}

// Test GetQuestionById - Not found
func TestTriviaService_GetQuestionById_NotFound(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	mockTriviaRepo.On("GetQuestionById", int64(999)).Return(nil, errors.New("question not found"))

	// Act
	result, err := triviaService.GetQuestionById(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockTriviaRepo.AssertExpectations(t)
}

// Test CreateQuestion - Success
func TestTriviaService_CreateQuestion_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	createDTO := &models.TriviaQuestionCreateDTO{
		Question:      "What is the capital of Spain?",
		CorrectAnswer: "Madrid",
		Tags:          []string{"Geography", "Europe"},
		IsPublished:   true,
	}

	expectedQuestion := &repositories.TriviaQuestionEntity{
		ID:            1,
		Question:      "What is the capital of Spain?",
		CorrectAnswer: "Madrid",
		Tags:          pq.StringArray{"geography", "europe"}, // Tags should be lowercased
		IsPublished:   true,
		IsArchived:    false,
		CreatedAt:     time.Now(),
		ModifiedAt:    nil,
	}

	// Mock duplicate check - not implemented in service yet, so skip this mock
	mockTriviaRepo.On("CreateQuestion", mock.AnythingOfType("*models.TriviaQuestionCreateDTO")).Return(expectedQuestion, nil)

	// Act
	result, err := triviaService.CreateQuestion(createDTO)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "What is the capital of Spain?", result.Question)
	assert.Equal(t, "Madrid", result.CorrectAnswer)
	assert.Equal(t, pq.StringArray{"geography", "europe"}, result.Tags)
	mockTriviaRepo.AssertExpectations(t)
}

// Test CreateQuestion - Validation error
func TestTriviaService_CreateQuestion_ValidationError(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	createDTO := &models.TriviaQuestionCreateDTO{
		Question:      "", // Empty question
		CorrectAnswer: "Madrid",
		Tags:          []string{"geography"},
		IsPublished:   true,
	}

	// Act
	result, err := triviaService.CreateQuestion(createDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "question and correct answer are required")
}

// Test CreateQuestion - Duplicate error (Skip since not implemented in service)
func TestTriviaService_CreateQuestion_DuplicateError(t *testing.T) {
	// TODO: Implement duplicate check in service layer
	t.Skip("Duplicate check not implemented in service layer yet")
}

// Test UpdateQuestion - Success
func TestTriviaService_UpdateQuestion_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	updateDTO := &models.TriviaQuestionUpdateDTO{
		Question:      "What is the capital of Italy?",
		CorrectAnswer: "Rome",
		Tags:          []string{"Geography", "Europe"},
		IsPublished:   true,
	}

	expectedQuestion := &repositories.TriviaQuestionEntity{
		ID:            1,
		Question:      "What is the capital of Italy?",
		CorrectAnswer: "Rome",
		Tags:          pq.StringArray{"geography", "europe"},
		IsPublished:   true,
		IsArchived:    false,
		CreatedAt:     time.Now(),
		ModifiedAt:    nil,
	}

	mockTriviaRepo.On("UpdateQuestion", mock.AnythingOfType("*models.TriviaQuestionUpdateDTO"), int64(1)).Return(expectedQuestion, nil)

	// Act
	result, err := triviaService.UpdateQuestion(updateDTO, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "What is the capital of Italy?", result.Question)
	assert.Equal(t, "Rome", result.CorrectAnswer)
	assert.Equal(t, pq.StringArray{"geography", "europe"}, result.Tags)
	mockTriviaRepo.AssertExpectations(t)
}

// Test UpdateQuestion - Validation error
func TestTriviaService_UpdateQuestion_ValidationError(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	updateDTO := &models.TriviaQuestionUpdateDTO{
		Question:      "What is the capital of Italy?",
		CorrectAnswer: "", // Empty answer
		Tags:          []string{"geography"},
		IsPublished:   true,
	}

	// Act
	result, err := triviaService.UpdateQuestion(updateDTO, 1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "question and correct answer are required")
}

// Test ToggleQuestionArchived - Success
func TestTriviaService_ToggleQuestionArchived_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	mockTriviaRepo.On("ToggleQuestionArchived", int64(1)).Return(nil)

	// Act
	err := triviaService.ToggleQuestionArchived(1)

	// Assert
	assert.NoError(t, err)
	mockTriviaRepo.AssertExpectations(t)
}

// Test ToggleQuestionPublished - Success
func TestTriviaService_ToggleQuestionPublished_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	mockTriviaRepo.On("ToggleQuestionPublished", int64(1)).Return(nil)

	// Act
	err := triviaService.ToggleQuestionPublished(1)

	// Assert
	assert.NoError(t, err)
	mockTriviaRepo.AssertExpectations(t)
}

// ===========================================
// WRONG ANSWER CRUD TESTS
// ===========================================

// Test GetWrongAnswers - Success with pagination
func TestTriviaService_GetWrongAnswers_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	count := 75
	expectedAnswers := []*repositories.WrongAnswerPoolEntity{
		{
			ID:         1,
			AnswerText: "London",
			Tags:       pq.StringArray{"geography", "europe"},
			IsArchived: false,
			CreatedAt:  time.Now(),
			ModifiedAt: nil,
		},
		{
			ID:         2,
			AnswerText: "Berlin",
			Tags:       pq.StringArray{"geography", "europe"},
			IsArchived: false,
			CreatedAt:  time.Now(),
			ModifiedAt: nil,
		},
	}

	mockTriviaRepo.On("GetWrongAnswersCount", "", "", "").Return(&count, nil)
	mockTriviaRepo.On("GetWrongAnswers", 25, 0, "", "", "").Return(expectedAnswers, nil)

	// Act
	result, err := triviaService.GetWrongAnswers(25, 0, "", "", "")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 75, result.Total)
	assert.Equal(t, 25, result.PageSize)
	assert.Equal(t, 1, result.Page)
	assert.Len(t, result.Results, 2)
	mockTriviaRepo.AssertExpectations(t)
}

// Test GetWrongAnswers - Success with filters
func TestTriviaService_GetWrongAnswers_WithFilters_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	count := 5
	expectedAnswers := []*repositories.WrongAnswerPoolEntity{
		{
			ID:         1,
			AnswerText: "London",
			Tags:       pq.StringArray{"geography", "europe"},
			IsArchived: false,
			CreatedAt:  time.Now(),
			ModifiedAt: nil,
		},
	}

	mockTriviaRepo.On("GetWrongAnswersCount", "london", "active", "geography").Return(&count, nil)
	mockTriviaRepo.On("GetWrongAnswers", 10, 0, "london", "active", "geography").Return(expectedAnswers, nil)

	// Act
	result, err := triviaService.GetWrongAnswers(10, 0, "london", "active", "geography")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 5, result.Total)
	assert.Len(t, result.Results, 1)
	mockTriviaRepo.AssertExpectations(t)
}

// Test GetWrongAnswerById - Success
func TestTriviaService_GetWrongAnswerById_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	expectedAnswer := &repositories.WrongAnswerPoolEntity{
		ID:         1,
		AnswerText: "London",
		Tags:       pq.StringArray{"geography", "europe"},
		IsArchived: false,
		CreatedAt:  time.Now(),
		ModifiedAt: nil,
	}

	mockTriviaRepo.On("GetWrongAnswerById", int64(1)).Return(expectedAnswer, nil)

	// Act
	result, err := triviaService.GetWrongAnswerById(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "London", result.AnswerText)
	assert.Equal(t, pq.StringArray{"geography", "europe"}, result.Tags)
	mockTriviaRepo.AssertExpectations(t)
}

// Test GetWrongAnswerById - Not found
func TestTriviaService_GetWrongAnswerById_NotFound(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	mockTriviaRepo.On("GetWrongAnswerById", int64(999)).Return(nil, errors.New("wrong answer not found"))

	// Act
	result, err := triviaService.GetWrongAnswerById(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockTriviaRepo.AssertExpectations(t)
}

// Test CreateWrongAnswer - Success
func TestTriviaService_CreateWrongAnswer_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	createDTO := &models.WrongAnswerCreateDTO{
		AnswerText: "Berlin",
		Tags:       []string{"Geography", "Europe"},
	}

	expectedAnswer := &repositories.WrongAnswerPoolEntity{
		ID:         1,
		AnswerText: "Berlin",
		Tags:       pq.StringArray{"geography", "europe"}, // Tags should be lowercased
		IsArchived: false,
		CreatedAt:  time.Now(),
		ModifiedAt: nil,
	}

	// Mock duplicate check - not implemented in service yet, so skip this mock
	mockTriviaRepo.On("CreateWrongAnswer", mock.AnythingOfType("*models.WrongAnswerCreateDTO")).Return(expectedAnswer, nil)

	// Act
	result, err := triviaService.CreateWrongAnswer(createDTO)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Berlin", result.AnswerText)
	assert.Equal(t, pq.StringArray{"geography", "europe"}, result.Tags)
	mockTriviaRepo.AssertExpectations(t)
}

// Test CreateWrongAnswer - Validation error
func TestTriviaService_CreateWrongAnswer_ValidationError(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	createDTO := &models.WrongAnswerCreateDTO{
		AnswerText: "", // Empty answer text
		Tags:       []string{"geography"},
	}

	// Act
	result, err := triviaService.CreateWrongAnswer(createDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "answer text is required")
}

// Test CreateWrongAnswer - Duplicate error (Skip since not implemented in service)
func TestTriviaService_CreateWrongAnswer_DuplicateError(t *testing.T) {
	// TODO: Implement duplicate check in service layer
	t.Skip("Duplicate check not implemented in service layer yet")
}

// Test UpdateWrongAnswer - Success
func TestTriviaService_UpdateWrongAnswer_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	updateDTO := &models.WrongAnswerUpdateDTO{
		AnswerText: "Munich",
		Tags:       []string{"Geography", "Europe"},
	}

	expectedAnswer := &repositories.WrongAnswerPoolEntity{
		ID:         1,
		AnswerText: "Munich",
		Tags:       pq.StringArray{"geography", "europe"},
		IsArchived: false,
		CreatedAt:  time.Now(),
		ModifiedAt: nil,
	}

	mockTriviaRepo.On("UpdateWrongAnswer", mock.AnythingOfType("*models.WrongAnswerUpdateDTO"), int64(1)).Return(expectedAnswer, nil)

	// Act
	result, err := triviaService.UpdateWrongAnswer(updateDTO, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Munich", result.AnswerText)
	assert.Equal(t, pq.StringArray{"geography", "europe"}, result.Tags)
	mockTriviaRepo.AssertExpectations(t)
}

// Test UpdateWrongAnswer - Validation error
func TestTriviaService_UpdateWrongAnswer_ValidationError(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	updateDTO := &models.WrongAnswerUpdateDTO{
		AnswerText: "", // Empty answer text
		Tags:       []string{"geography"},
	}

	// Act
	result, err := triviaService.UpdateWrongAnswer(updateDTO, 1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "answer text is required")
}

// Test ToggleWrongAnswerArchived - Success
func TestTriviaService_ToggleWrongAnswerArchived_Success(t *testing.T) {
	// Arrange
	mockTriviaRepo := new(MockTriviaRepository)
	triviaService := NewTriviaService(mockTriviaRepo)

	mockTriviaRepo.On("ToggleWrongAnswerArchived", int64(1)).Return(nil)

	// Act
	err := triviaService.ToggleWrongAnswerArchived(1)

	// Assert
	assert.NoError(t, err)
	mockTriviaRepo.AssertExpectations(t)
}
