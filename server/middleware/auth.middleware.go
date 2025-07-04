package middleware

import (
	"errors"
	"net/http"

	"github.com/snowlynxsoftware/oto-api/server/database/repositories"
	"github.com/snowlynxsoftware/oto-api/server/services"
	"github.com/snowlynxsoftware/oto-api/server/util"
)

type IAuthMiddleware interface {
	Authorize(r *http.Request, requiredUserTypeKeys []string) (*AuthorizedUserContext, error)
}

type AuthMiddleware struct {
	userRepository repositories.IUserRepository
	tokenService   services.ITokenService
}

func NewAuthMiddleware(userRepository repositories.IUserRepository, tokenService services.ITokenService) IAuthMiddleware {
	return &AuthMiddleware{
		userRepository: userRepository,
		tokenService:   tokenService,
	}
}

// If a request is authorized, it will return this context to the controller
// so that information from the user can be used as an immutable object.
type AuthorizedUserContext struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	IsAdmin   bool   `json:"is_admin,omitempty"`   // Optional, only set if the user is an admin
	IsSupport bool   `json:"is_support,omitempty"` // Optional, only set if the user is a support agent
}

func (m *AuthMiddleware) Authorize(r *http.Request, requiredUserTypeKeys []string) (*AuthorizedUserContext, error) {

	cookie, err := r.Cookie("access_token")
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return nil, errors.New("access token not found in request")
	}

	userId, err := m.tokenService.ValidateToken(&cookie.Value)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return nil, err
	}

	userEntity, err := m.userRepository.GetUserById(*userId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return nil, err
	}

	if userEntity.IsArchived {
		return nil, errors.New("user is archived")
	}

	if !userEntity.IsVerified {
		return nil, errors.New("user is not verified")
	}

	if userEntity.IsBanned {
		return nil, errors.New("user is banned. Reason - " + *userEntity.BanReason)
	}

	if requiredUserTypeKeys != nil {
		isUserTypeAllowed := false
		for _, key := range requiredUserTypeKeys {
			if userEntity.UserTypeKey == key {
				isUserTypeAllowed = true
				break
			}
		}
		if !isUserTypeAllowed {
			return nil, errors.New("forbidden - user type is not allowed")
		}
	}

	return &AuthorizedUserContext{
		Id:        int(userEntity.ID),
		Email:     userEntity.Email,
		Username:  userEntity.DisplayName,
		IsAdmin:   userEntity.UserTypeKey == repositories.UserTypeAdmin,
		IsSupport: userEntity.UserTypeKey == repositories.UserTypeSupport,
	}, nil

}
