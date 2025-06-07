package repositories

import (
	"time"

	"github.com/lib/pq"
	"github.com/snowlynxsoftware/oto-api/server/database"
	"github.com/snowlynxsoftware/oto-api/server/models"
)

type TriviaQuestionEntity struct {
	ID            int64          `json:"id" db:"id"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	ModifiedAt    *time.Time     `json:"modified_at" db:"modified_at"`
	IsArchived    bool           `json:"is_archived" db:"is_archived"`
	IsPublished   bool           `json:"is_published" db:"is_published"`
	Question      string         `json:"question" db:"question"`
	CorrectAnswer string         `json:"correct_answer" db:"correct_answer"`
	Tags          pq.StringArray `json:"tags" db:"tags"`
}

type WrongAnswerPoolEntity struct {
	ID         int64          `json:"id" db:"id"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time     `json:"modified_at" db:"modified_at"`
	IsArchived bool           `json:"is_archived" db:"is_archived"`
	AnswerText string         `json:"answer_text" db:"answer_text"`
	Tags       pq.StringArray `json:"tags" db:"tags"`
}

type TriviaDeckEntity struct {
	ID           int64      `json:"id" db:"id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt   *time.Time `json:"modified_at" db:"modified_at"`
	IsArchived   bool       `json:"is_archived" db:"is_archived"`
	Name         string     `json:"name" db:"name"`
	Description  *string    `json:"description" db:"description"`
	IsApproved   bool       `json:"is_approved" db:"is_approved"`
	IsSystemDeck bool       `json:"is_system_deck" db:"is_system_deck"`
}

type ITriviaRepository interface {
	GetTriviaQuestionByText(question string) (*TriviaQuestionEntity, error)
	ImportTriviaQuestions(data []models.TriviaQuestionImportData) (*models.TriviaQuestionImportResults, error)
	GetWrongAnswerByText(answer string) (*WrongAnswerPoolEntity, error)
	ImportWrongAnswers(data []models.TriviaWrongAnswerImportData) (*models.TriviaWrongAnswerImportResults, error)
	GetTriviaDeckById(deckId int64) (*TriviaDeckEntity, error)
	CreateNewTriviaDeck(name string, description string, isSystemDeck bool) (*TriviaDeckEntity, error)
	UpdateTriviaDeckMetadata(deckId int64, name string, description string) (*TriviaDeckEntity, error)
	UpdateTriviaDeckApprovalStatus(deckId int64, isApproved bool) (*TriviaDeckEntity, error)
	UpdateTriviaDeckArchivalStatus(deckId int64, isArchived bool) (*TriviaDeckEntity, error)
}

type TriviaRepository struct {
	db *database.AppDataSource
}

func NewTriviaRepository(db *database.AppDataSource) ITriviaRepository {
	return &TriviaRepository{
		db: db,
	}
}

func (r *TriviaRepository) GetTriviaQuestionByText(question string) (*TriviaQuestionEntity, error) {
	questionEntity := TriviaQuestionEntity{}
	sql := `SELECT
		id, created_at, modified_at, is_archived, is_published, question, correct_answer, tags
	FROM trivia_questions
	WHERE question = $1`
	err := r.db.DB.Get(&questionEntity, sql, question)
	if err != nil {
		return nil, err
	}
	return &questionEntity, nil
}

func (r *TriviaRepository) ImportTriviaQuestions(data []models.TriviaQuestionImportData) (*models.TriviaQuestionImportResults, error) {
	results := &models.TriviaQuestionImportResults{
		TotalQuestionsProcessed: int64(len(data)),
		QuestionsAdded:          0,
	}

	for _, questionData := range data {
		question := &TriviaQuestionEntity{
			CreatedAt:     time.Now(),
			IsArchived:    false,
			IsPublished:   true,
			Question:      questionData.Question,
			CorrectAnswer: questionData.CorrectAnswer,
			Tags:          questionData.Tags,
		}

		existingQuestion, _ := r.GetTriviaQuestionByText(question.Question)

		if existingQuestion == nil || existingQuestion.ID == 0 {

			sql := `INSERT INTO trivia_questions (question, correct_answer, tags)
				VALUES ($1, $2, $3) RETURNING id`
			err := r.db.DB.QueryRow(sql, question.Question, question.CorrectAnswer, pq.Array(question.Tags)).Scan(&question.ID)
			if err != nil {
				return nil, err
			}

			results.QuestionsAdded++
		}
	}

	return results, nil
}

func (r *TriviaRepository) GetWrongAnswerByText(answer string) (*WrongAnswerPoolEntity, error) {
	entity := WrongAnswerPoolEntity{}
	sql := `SELECT
		id, created_at, modified_at, is_archived, answer_text, tags
	FROM wrong_answer_pool
	WHERE answer_text = $1`
	err := r.db.DB.Get(&entity, sql, answer)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *TriviaRepository) ImportWrongAnswers(data []models.TriviaWrongAnswerImportData) (*models.TriviaWrongAnswerImportResults, error) {
	results := &models.TriviaWrongAnswerImportResults{
		TotalAnswersProcessed: int64(len(data)),
		AnswersAdded:          0,
	}

	for _, answerData := range data {
		answer := &WrongAnswerPoolEntity{
			CreatedAt:  time.Now(),
			IsArchived: false,
			AnswerText: answerData.AnswerText,
			Tags:       answerData.Tags,
		}

		existingQuestion, _ := r.GetTriviaQuestionByText(answer.AnswerText)

		if existingQuestion == nil || existingQuestion.ID == 0 {

			sql := `INSERT INTO wrong_answer_pool (answer_text, tags)
				VALUES ($1, $2) RETURNING id`
			err := r.db.DB.QueryRow(sql, answer.AnswerText, pq.Array(answer.Tags)).Scan(&answer.ID)
			if err != nil {
				return nil, err
			}

			results.AnswersAdded++
		}
	}

	return results, nil
}

func (r *TriviaRepository) GetTriviaDeckById(deckId int64) (*TriviaDeckEntity, error) {
	deckEntity := TriviaDeckEntity{}
	sql := `SELECT
		id, created_at, modified_at, is_archived, name, description, is_approved, is_system_deck
	FROM trivia_decks
	WHERE id = $1`
	err := r.db.DB.Get(&deckEntity, sql, deckId)
	if err != nil {
		return nil, err
	}
	return &deckEntity, nil
}

func (r *TriviaRepository) CreateNewTriviaDeck(name string, description string, isSystemDeck bool) (*TriviaDeckEntity, error) {
	deck := &TriviaDeckEntity{
		CreatedAt:    time.Now(),
		IsArchived:   false,
		Name:         name,
		Description:  &description,
		IsApproved:   false,
		IsSystemDeck: isSystemDeck,
	}

	sql := `INSERT INTO trivia_decks (name, description, is_system_deck)
			VALUES ($1, $2, $3) RETURNING id`
	err := r.db.DB.QueryRow(sql, deck.Name, deck.Description, deck.IsSystemDeck).Scan(&deck.ID)
	if err != nil {
		return nil, err
	}

	return deck, nil
}

func (r *TriviaRepository) UpdateTriviaDeckMetadata(deckId int64, name string, description string) (*TriviaDeckEntity, error) {
	deck, err := r.GetTriviaDeckById(deckId)
	if err != nil {
		return nil, err
	}

	deck.Name = name
	deck.Description = &description
	deck.ModifiedAt = &time.Time{}

	sql := `UPDATE trivia_decks SET name = $1, description = $2, modified_at = NOW()
			WHERE id = $3 RETURNING modified_at`
	err = r.db.DB.QueryRow(sql, deck.Name, deck.Description, deckId).Scan(&deck.ModifiedAt)
	if err != nil {
		return nil, err
	}

	return deck, nil
}

func (r *TriviaRepository) UpdateTriviaDeckApprovalStatus(deckId int64, isApproved bool) (*TriviaDeckEntity, error) {
	deck, err := r.GetTriviaDeckById(deckId)
	if err != nil {
		return nil, err
	}

	deck.IsApproved = isApproved
	deck.ModifiedAt = &time.Time{}

	sql := `UPDATE trivia_decks SET is_approved = $1, modified_at = NOW()
			WHERE id = $2 RETURNING modified_at`
	err = r.db.DB.QueryRow(sql, deck.IsApproved, deckId).Scan(&deck.ModifiedAt)
	if err != nil {
		return nil, err
	}

	return deck, nil
}

func (r *TriviaRepository) UpdateTriviaDeckArchivalStatus(deckId int64, isArchived bool) (*TriviaDeckEntity, error) {
	deck, err := r.GetTriviaDeckById(deckId)
	if err != nil {
		return nil, err
	}

	deck.IsArchived = isArchived
	deck.ModifiedAt = &time.Time{}

	sql := `UPDATE trivia_decks SET is_archived = $1, modified_at = NOW()
			WHERE id = $2 RETURNING modified_at`
	err = r.db.DB.QueryRow(sql, deck.IsArchived, deckId).Scan(&deck.ModifiedAt)
	if err != nil {
		return nil, err
	}

	return deck, nil
}
