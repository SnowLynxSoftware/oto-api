package services

import (
	"errors"

	"github.com/snowlynxsoftware/oto-api/server/database/repositories"
	"github.com/snowlynxsoftware/oto-api/server/models"
)

type ITriviaService interface {
	ImportTriviaQuestions(data []models.TriviaQuestionImportData) (*models.TriviaQuestionImportResults, error)
	ImportWrongAnswers(data []models.TriviaWrongAnswerImportData) (*models.TriviaWrongAnswerImportResults, error)
	GetTriviaDeckById(deckId int64) (*repositories.TriviaDeckEntity, error)
	CreateNewTriviaDeck(name string, description string, isSystemDeck bool) (*repositories.TriviaDeckEntity, error)
	UpdateTriviaDeckMetadata(deckId int64, name string, description string) (*repositories.TriviaDeckEntity, error)
	UpdateTriviaDeckApprovalStatus(deckId int64, isApproved bool) (*repositories.TriviaDeckEntity, error)
	UpdateTriviaDeckArchivalStatus(deckId int64, isArchived bool) (*repositories.TriviaDeckEntity, error)

	// New CRUD methods for questions
	GetQuestions(pageSize, offset int, searchString, statusFilter, tagFilter string) (*models.PaginatedResponse, error)
	GetQuestionById(id int64) (*repositories.TriviaQuestionEntity, error)
	CreateQuestion(dto *models.TriviaQuestionCreateDTO) (*repositories.TriviaQuestionEntity, error)
	UpdateQuestion(dto *models.TriviaQuestionUpdateDTO, id int64) (*repositories.TriviaQuestionEntity, error)
	ToggleQuestionArchived(id int64) error
	ToggleQuestionPublished(id int64) error

	// New CRUD methods for wrong answers
	GetWrongAnswers(pageSize, offset int, searchString, statusFilter, tagFilter string) (*models.PaginatedResponse, error)
	GetWrongAnswerById(id int64) (*repositories.WrongAnswerPoolEntity, error)
	CreateWrongAnswer(dto *models.WrongAnswerCreateDTO) (*repositories.WrongAnswerPoolEntity, error)
	UpdateWrongAnswer(dto *models.WrongAnswerUpdateDTO, id int64) (*repositories.WrongAnswerPoolEntity, error)
	ToggleWrongAnswerArchived(id int64) error
}

type TriviaService struct {
	triviaRepository repositories.ITriviaRepository
}

func NewTriviaService(triviaRepository repositories.ITriviaRepository) ITriviaService {
	return &TriviaService{
		triviaRepository: triviaRepository,
	}
}

func (s *TriviaService) ImportTriviaQuestions(data []models.TriviaQuestionImportData) (*models.TriviaQuestionImportResults, error) {
	results, err := s.triviaRepository.ImportTriviaQuestions(data)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *TriviaService) ImportWrongAnswers(data []models.TriviaWrongAnswerImportData) (*models.TriviaWrongAnswerImportResults, error) {
	results, err := s.triviaRepository.ImportWrongAnswers(data)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *TriviaService) GetTriviaDeckById(deckId int64) (*repositories.TriviaDeckEntity, error) {
	deck, err := s.triviaRepository.GetTriviaDeckById(deckId)
	if err != nil {
		return nil, err
	}
	return deck, nil
}

func (s *TriviaService) CreateNewTriviaDeck(name string, description string, isSystemDeck bool) (*repositories.TriviaDeckEntity, error) {
	deck, err := s.triviaRepository.CreateNewTriviaDeck(name, description, isSystemDeck)
	if err != nil {
		return nil, err
	}
	return deck, nil
}

func (s *TriviaService) UpdateTriviaDeckMetadata(deckId int64, name string, description string) (*repositories.TriviaDeckEntity, error) {
	deck, err := s.triviaRepository.UpdateTriviaDeckMetadata(deckId, name, description)
	if err != nil {
		return nil, err
	}
	return deck, nil
}

func (s *TriviaService) UpdateTriviaDeckApprovalStatus(deckId int64, isApproved bool) (*repositories.TriviaDeckEntity, error) {
	deck, err := s.triviaRepository.UpdateTriviaDeckApprovalStatus(deckId, isApproved)
	if err != nil {
		return nil, err
	}
	return deck, nil
}

func (s *TriviaService) UpdateTriviaDeckArchivalStatus(deckId int64, isArchived bool) (*repositories.TriviaDeckEntity, error) {
	deck, err := s.triviaRepository.UpdateTriviaDeckArchivalStatus(deckId, isArchived)
	if err != nil {
		return nil, err
	}
	return deck, nil
}

// Question CRUD methods
func (s *TriviaService) GetQuestions(pageSize, offset int, searchString, statusFilter, tagFilter string) (*models.PaginatedResponse, error) {
	questions, err := s.triviaRepository.GetQuestions(pageSize, offset, searchString, statusFilter, tagFilter)
	if err != nil {
		return nil, err
	}

	count, err := s.triviaRepository.GetQuestionsCount(searchString, statusFilter, tagFilter)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice
	results := make([]interface{}, len(questions))
	for i, question := range questions {
		results[i] = question
	}

	currentPage := (offset / pageSize) + 1

	return &models.PaginatedResponse{
		Results:  results,
		Total:    *count,
		PageSize: pageSize,
		Page:     currentPage,
	}, nil
}

func (s *TriviaService) GetQuestionById(id int64) (*repositories.TriviaQuestionEntity, error) {
	question, err := s.triviaRepository.GetQuestionById(id)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (s *TriviaService) CreateQuestion(dto *models.TriviaQuestionCreateDTO) (*repositories.TriviaQuestionEntity, error) {
	// Validate required fields
	if dto.Question == "" || dto.CorrectAnswer == "" {
		return nil, errors.New("question and correct answer are required")
	}

	question, err := s.triviaRepository.CreateQuestion(dto)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (s *TriviaService) UpdateQuestion(dto *models.TriviaQuestionUpdateDTO, id int64) (*repositories.TriviaQuestionEntity, error) {
	// Validate required fields
	if dto.Question == "" || dto.CorrectAnswer == "" {
		return nil, errors.New("question and correct answer are required")
	}

	question, err := s.triviaRepository.UpdateQuestion(dto, id)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (s *TriviaService) ToggleQuestionArchived(id int64) error {
	return s.triviaRepository.ToggleQuestionArchived(id)
}

func (s *TriviaService) ToggleQuestionPublished(id int64) error {
	return s.triviaRepository.ToggleQuestionPublished(id)
}

// Wrong Answer CRUD methods
func (s *TriviaService) GetWrongAnswers(pageSize, offset int, searchString, statusFilter, tagFilter string) (*models.PaginatedResponse, error) {
	answers, err := s.triviaRepository.GetWrongAnswers(pageSize, offset, searchString, statusFilter, tagFilter)
	if err != nil {
		return nil, err
	}

	count, err := s.triviaRepository.GetWrongAnswersCount(searchString, statusFilter, tagFilter)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice
	results := make([]interface{}, len(answers))
	for i, answer := range answers {
		results[i] = answer
	}

	currentPage := (offset / pageSize) + 1

	return &models.PaginatedResponse{
		Results:  results,
		Total:    *count,
		PageSize: pageSize,
		Page:     currentPage,
	}, nil
}

func (s *TriviaService) GetWrongAnswerById(id int64) (*repositories.WrongAnswerPoolEntity, error) {
	answer, err := s.triviaRepository.GetWrongAnswerById(id)
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (s *TriviaService) CreateWrongAnswer(dto *models.WrongAnswerCreateDTO) (*repositories.WrongAnswerPoolEntity, error) {
	// Validate required fields
	if dto.AnswerText == "" {
		return nil, errors.New("answer text is required")
	}

	answer, err := s.triviaRepository.CreateWrongAnswer(dto)
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (s *TriviaService) UpdateWrongAnswer(dto *models.WrongAnswerUpdateDTO, id int64) (*repositories.WrongAnswerPoolEntity, error) {
	// Validate required fields
	if dto.AnswerText == "" {
		return nil, errors.New("answer text is required")
	}

	answer, err := s.triviaRepository.UpdateWrongAnswer(dto, id)
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (s *TriviaService) ToggleWrongAnswerArchived(id int64) error {
	return s.triviaRepository.ToggleWrongAnswerArchived(id)
}
