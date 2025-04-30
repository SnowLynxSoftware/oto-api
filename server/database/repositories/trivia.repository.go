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
	ID         int64      `json:"id" db:"id"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at" db:"modified_at"`
	IsArchived bool       `json:"is_archived" db:"is_archived"`
	AnswerText string     `json:"answer_text" db:"answer_text"`
	Tags       []string   `json:"tags" db:"tags"`
}

type ITriviaRepository interface {
	GetTriviaQuestionByText(question string) (*TriviaQuestionEntity, error)
	ImportTriviaQuestions(data []models.TriviaQuestionImportData) (*models.TriviaQuestionImportResults, error)
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
