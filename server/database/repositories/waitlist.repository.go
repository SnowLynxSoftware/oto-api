package repositories

import (
	"time"

	"github.com/snowlynxsoftware/oto-api/server/database"
	"github.com/snowlynxsoftware/oto-api/server/models"
)

type WaitlistEntity struct {
	ID        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Email     string    `json:"email" db:"email"`
}

type IWaitlistRepository interface {
	GetWaitlistEntryById(id int) (*WaitlistEntity, error)
	CreateNewWaitlistEntry(dto *models.WaitlistCreateDTO) (*WaitlistEntity, error)
}

type WaitlistRepository struct {
	db *database.AppDataSource
}

func NewWaitlistRepository(db *database.AppDataSource) IWaitlistRepository {
	return &WaitlistRepository{
		db: db,
	}
}

func (r *WaitlistRepository) GetWaitlistEntryById(id int) (*WaitlistEntity, error) {
	entity := &WaitlistEntity{}
	sql := `SELECT
		*
	FROM waitlist
	WHERE id = $1`
	err := r.db.DB.Get(entity, sql, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *WaitlistRepository) CreateNewWaitlistEntry(dto *models.WaitlistCreateDTO) (*WaitlistEntity, error) {
	sql := `INSERT INTO waitlist (email)
    VALUES ($1)
    RETURNING id;`
	row := r.db.DB.QueryRow(sql, dto.Email)
	var insertedId int
	err := row.Scan(&insertedId)
	if err != nil {
		return nil, err
	}

	user, err := r.GetWaitlistEntryById(insertedId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
