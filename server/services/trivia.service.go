package services

import (
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
