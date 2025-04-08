package services

import (
	"errors"

	"github.com/snowlynxsoftware/oto-api/server/database/repositories"
	"github.com/snowlynxsoftware/oto-api/server/models"
	"github.com/snowlynxsoftware/oto-api/server/util"
)

type IWaitlistService interface {
	CreateNewWaitlistEntry(dto *models.WaitlistCreateDTO) (*repositories.WaitlistEntity, error)
}

type WaitlistService struct {
	waitlistRepository repositories.IWaitlistRepository
}

func NewWaitlistService(waitlistRepository repositories.IWaitlistRepository) IWaitlistService {
	return &WaitlistService{
		waitlistRepository: waitlistRepository,
	}
}

func (s *WaitlistService) CreateNewWaitlistEntry(dto *models.WaitlistCreateDTO) (*repositories.WaitlistEntity, error) {

	if len(dto.Email) == 0 {
		return nil, errors.New("email cannot be empty")
	}
	if len(dto.Email) < 5 {
		return nil, errors.New("email must be at least 5 characters long")
	}

	waitlistEntry, err := s.waitlistRepository.CreateNewWaitlistEntry(dto)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return nil, err
	}
	return waitlistEntry, nil
}
