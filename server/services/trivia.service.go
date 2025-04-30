package services

import (
	"github.com/snowlynxsoftware/oto-api/server/database/repositories"
	"github.com/snowlynxsoftware/oto-api/server/models"
)

type ITriviaService interface {
	ImportTriviaQuestions(data []models.TriviaQuestionImportData) (*models.TriviaQuestionImportResults, error)
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
