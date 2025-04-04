package repositories

import (
	"time"

	"github.com/snowlynxsoftware/oto-api/server/database"
	"github.com/snowlynxsoftware/oto-api/server/models"
)

type UserEntity struct {
	ID           int64      `json:"id" db:"id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt   *time.Time `json:"modified_at" db:"modified_at"`
	IsArchived   bool       `json:"is_archived" db:"is_archived"`
	Email        string     `json:"email" db:"email"`
	DisplayName  string     `json:"display_name" db:"display_name"`
	AvatarURL    *string    `json:"avatar_url" db:"avatar_url"`
	ProfileText  *string    `json:"profile_text" db:"profile_text"`
	IsVerified   bool       `json:"is_verified" db:"is_verified"`
	UserTypeKey  string     `json:"user_type_key" db:"user_type_key"`
	PasswordHash *string    `json:"-" db:"password_hash"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`
	IsBanned     bool       `json:"is_banned" db:"is_banned"`
	BanReason    *string    `json:"ban_reason" db:"ban_reason"`
}

type IUserRepository interface {
	GetUserById(id int) (*UserEntity, error)
	GetUserByEmail(email string) (*UserEntity, error)
	CreateNewUser(dto *models.UserCreateDTO) (*UserEntity, error)
	MarkUserVerified(userId *int) (bool, error)
	UpdateUserLastLogin(userId *int) (bool, error)
	UpdateUserPassword(userId *int, password string) (bool, error)
	BanUserByIdWithReason(userId *int, reason string) (bool, error)
	UnbanUserById(userId *int) (bool, error)
	SetUserTypeKey(userId *int, key string) (bool, error)
}

type UserRepository struct {
	db *database.AppDataSource
}

func NewUserRepository(db *database.AppDataSource) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserById(id int) (*UserEntity, error) {
	user := &UserEntity{}
	sql := `SELECT
		*
	FROM users
	WHERE id = $1`
	err := r.db.DB.Get(user, sql, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*UserEntity, error) {
	user := &UserEntity{}
	sql := `SELECT
		*
	FROM users
	WHERE email = $1`
	err := r.db.DB.Get(user, sql, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CreateNewUser(dto *models.UserCreateDTO) (*UserEntity, error) {
	sql := `INSERT INTO users (email, display_name, password_hash, user_type_key)
    VALUES ($1, $2, $3, $4)
    RETURNING id;`
	row := r.db.DB.QueryRow(sql, dto.Email, dto.DisplayName, dto.Password, "player")
	var insertedId int
	err := row.Scan(&insertedId)
	if err != nil {
		return nil, err
	}

	user, err := r.GetUserById(insertedId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) MarkUserVerified(userId *int) (bool, error) {
	sql := `UPDATE users SET is_verified = true WHERE id = $1;`
	_, err := r.db.DB.Exec(sql, &userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) UpdateUserLastLogin(userId *int) (bool, error) {
	sql := `UPDATE users SET last_login = NOW() WHERE id = $1;`
	_, err := r.db.DB.Exec(sql, &userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) UpdateUserPassword(userId *int, password string) (bool, error) {
	sql := `UPDATE users SET password_hash = $1 WHERE id = $2;`
	_, err := r.db.DB.Exec(sql, password, &userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) BanUserByIdWithReason(userId *int, reason string) (bool, error) {
	sql := `UPDATE users
		SET
			is_banned = true,
			ban_reason = $1
		WHERE id = $2;`
	_, err := r.db.DB.Exec(sql, reason, &userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) UnbanUserById(userId *int) (bool, error) {
	sql := `UPDATE users
		SET
			is_banned = false,
			ban_reason = ''
		WHERE id = $1;`
	_, err := r.db.DB.Exec(sql, &userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) SetUserTypeKey(userId *int, key string) (bool, error) {
	sql := `UPDATE users
		SET
			user_type_id = $1
		WHERE id = $2;`
	_, err := r.db.DB.Exec(sql, key, &userId)
	if err != nil {
		return false, err
	}

	return true, nil
}
