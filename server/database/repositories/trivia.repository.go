package repositories

import (
	"errors"
	"fmt"
	"strings"
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

	// New CRUD methods for questions
	GetQuestionsCount(searchString, statusFilter, tagFilter string) (*int, error)
	GetQuestions(pageSize, offset int, searchString, statusFilter, tagFilter string) ([]*TriviaQuestionEntity, error)
	GetQuestionById(id int64) (*TriviaQuestionEntity, error)
	CreateQuestion(dto *models.TriviaQuestionCreateDTO) (*TriviaQuestionEntity, error)
	UpdateQuestion(dto *models.TriviaQuestionUpdateDTO, id int64) (*TriviaQuestionEntity, error)
	ToggleQuestionArchived(id int64) error
	ToggleQuestionPublished(id int64) error

	// New CRUD methods for wrong answers
	GetWrongAnswersCount(searchString, statusFilter, tagFilter string) (*int, error)
	GetWrongAnswers(pageSize, offset int, searchString, statusFilter, tagFilter string) ([]*WrongAnswerPoolEntity, error)
	GetWrongAnswerById(id int64) (*WrongAnswerPoolEntity, error)
	CreateWrongAnswer(dto *models.WrongAnswerCreateDTO) (*WrongAnswerPoolEntity, error)
	UpdateWrongAnswer(dto *models.WrongAnswerUpdateDTO, id int64) (*WrongAnswerPoolEntity, error)
	ToggleWrongAnswerArchived(id int64) error
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

// Question CRUD methods
func (r *TriviaRepository) GetQuestionsCount(searchString, statusFilter, tagFilter string) (*int, error) {
	count := new(int)
	sql := `SELECT COUNT(*) as count FROM trivia_questions WHERE 1=1`

	// Build dynamic WHERE clause
	args := []interface{}{}
	argIndex := 1

	// Search filter
	if searchString != "" {
		sql += ` AND (question ILIKE '%' || $` + fmt.Sprintf("%d", argIndex) + ` || '%' OR correct_answer ILIKE '%' || $` + fmt.Sprintf("%d", argIndex) + ` || '%' OR $` + fmt.Sprintf("%d", argIndex) + ` = ANY(tags))`
		args = append(args, searchString)
		argIndex++
	}

	// Status filter
	if statusFilter != "" {
		switch statusFilter {
		case "active":
			sql += ` AND is_archived = false AND is_published = true`
		case "archived":
			sql += ` AND is_archived = true`
		case "published":
			sql += ` AND is_published = true`
		case "unpublished":
			sql += ` AND is_published = false`
		}
	}

	// Tag filter
	if tagFilter != "" {
		sql += ` AND tags && $` + fmt.Sprintf("%d", argIndex)
		// Split comma-separated tags and lowercase them
		tags := strings.Split(tagFilter, ",")
		for i, tag := range tags {
			tags[i] = strings.ToLower(strings.TrimSpace(tag))
		}
		args = append(args, pq.Array(tags))
		argIndex++
	}

	err := r.db.DB.Get(&count, sql, args...)
	if err != nil {
		return nil, err
	}
	return count, nil
}

func (r *TriviaRepository) GetQuestions(pageSize, offset int, searchString, statusFilter, tagFilter string) ([]*TriviaQuestionEntity, error) {
	questions := []*TriviaQuestionEntity{}
	sql := `SELECT id, created_at, modified_at, is_archived, is_published, question, correct_answer, tags FROM trivia_questions WHERE 1=1`

	// Build dynamic WHERE clause
	args := []interface{}{pageSize, offset}
	argIndex := 3

	// Search filter
	if searchString != "" {
		sql += ` AND (question ILIKE '%' || $` + fmt.Sprintf("%d", argIndex) + ` || '%' OR correct_answer ILIKE '%' || $` + fmt.Sprintf("%d", argIndex) + ` || '%' OR $` + fmt.Sprintf("%d", argIndex) + ` = ANY(tags))`
		args = append(args, searchString)
		argIndex++
	}

	// Status filter
	if statusFilter != "" {
		switch statusFilter {
		case "active":
			sql += ` AND is_archived = false AND is_published = true`
		case "archived":
			sql += ` AND is_archived = true`
		case "published":
			sql += ` AND is_published = true`
		case "unpublished":
			sql += ` AND is_published = false`
		}
	}

	// Tag filter
	if tagFilter != "" {
		sql += ` AND tags && $` + fmt.Sprintf("%d", argIndex)
		// Split comma-separated tags and lowercase them
		tags := strings.Split(tagFilter, ",")
		for i, tag := range tags {
			tags[i] = strings.ToLower(strings.TrimSpace(tag))
		}
		args = append(args, pq.Array(tags))
		argIndex++
	}

	sql += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	err := r.db.DB.Select(&questions, sql, args...)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *TriviaRepository) GetQuestionById(id int64) (*TriviaQuestionEntity, error) {
	question := &TriviaQuestionEntity{}
	sql := `SELECT id, created_at, modified_at, is_archived, is_published, question, correct_answer, tags FROM trivia_questions WHERE id = $1`
	err := r.db.DB.Get(question, sql, id)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (r *TriviaRepository) CreateQuestion(dto *models.TriviaQuestionCreateDTO) (*TriviaQuestionEntity, error) {
	// Check for duplicate
	existingQuestion, _ := r.GetTriviaQuestionByText(dto.Question)
	if existingQuestion != nil && existingQuestion.ID != 0 {
		return nil, errors.New("question already exists")
	}

	// Lowercase all tags
	tags := make([]string, len(dto.Tags))
	for i, tag := range dto.Tags {
		tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}

	question := &TriviaQuestionEntity{
		CreatedAt:     time.Now(),
		IsArchived:    false,
		IsPublished:   dto.IsPublished,
		Question:      dto.Question,
		CorrectAnswer: dto.CorrectAnswer,
		Tags:          tags,
	}

	sql := `INSERT INTO trivia_questions (question, correct_answer, tags, is_published) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.DB.QueryRow(sql, question.Question, question.CorrectAnswer, pq.Array(question.Tags), question.IsPublished).Scan(&question.ID)
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (r *TriviaRepository) UpdateQuestion(dto *models.TriviaQuestionUpdateDTO, id int64) (*TriviaQuestionEntity, error) {
	// Check if question exists
	existingQuestion, err := r.GetQuestionById(id)
	if err != nil {
		return nil, err
	}

	// Check for duplicate if question text changed
	if existingQuestion.Question != dto.Question {
		duplicate, _ := r.GetTriviaQuestionByText(dto.Question)
		if duplicate != nil && duplicate.ID != id {
			return nil, errors.New("question already exists")
		}
	}

	// Lowercase all tags
	tags := make([]string, len(dto.Tags))
	for i, tag := range dto.Tags {
		tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}

	sql := `UPDATE trivia_questions SET question = $1, correct_answer = $2, tags = $3, is_published = $4, modified_at = NOW() WHERE id = $5`
	_, err = r.db.DB.Exec(sql, dto.Question, dto.CorrectAnswer, pq.Array(tags), dto.IsPublished, id)
	if err != nil {
		return nil, err
	}

	return r.GetQuestionById(id)
}

func (r *TriviaRepository) ToggleQuestionArchived(id int64) error {
	// Get current question to check archive status
	question, err := r.GetQuestionById(id)
	if err != nil {
		return err
	}

	// If archiving, also unpublish
	if !question.IsArchived {
		sql := `UPDATE trivia_questions SET is_archived = NOT is_archived, is_published = false, modified_at = NOW() WHERE id = $1`
		_, err = r.db.DB.Exec(sql, id)
	} else {
		// If unarchiving, keep unpublished (manual publish required)
		sql := `UPDATE trivia_questions SET is_archived = NOT is_archived, modified_at = NOW() WHERE id = $1`
		_, err = r.db.DB.Exec(sql, id)
	}

	return err
}

func (r *TriviaRepository) ToggleQuestionPublished(id int64) error {
	// Get current question to check if archived
	question, err := r.GetQuestionById(id)
	if err != nil {
		return err
	}

	// Cannot publish/unpublish archived questions
	if question.IsArchived {
		return errors.New("cannot change published status of archived question")
	}

	sql := `UPDATE trivia_questions SET is_published = NOT is_published, modified_at = NOW() WHERE id = $1`
	_, err = r.db.DB.Exec(sql, id)
	return err
}

// Wrong Answer CRUD methods
func (r *TriviaRepository) GetWrongAnswersCount(searchString, statusFilter, tagFilter string) (*int, error) {
	count := new(int)
	sql := `SELECT COUNT(*) as count FROM wrong_answer_pool WHERE 1=1`

	// Build dynamic WHERE clause
	args := []interface{}{}
	argIndex := 1

	// Search filter
	if searchString != "" {
		sql += ` AND (answer_text ILIKE '%' || $` + fmt.Sprintf("%d", argIndex) + ` || '%' OR $` + fmt.Sprintf("%d", argIndex) + ` = ANY(tags))`
		args = append(args, searchString)
		argIndex++
	}

	// Status filter
	if statusFilter != "" {
		switch statusFilter {
		case "active":
			sql += ` AND is_archived = false`
		case "archived":
			sql += ` AND is_archived = true`
		}
	}

	// Tag filter
	if tagFilter != "" {
		sql += ` AND tags && $` + fmt.Sprintf("%d", argIndex)
		// Split comma-separated tags and lowercase them
		tags := strings.Split(tagFilter, ",")
		for i, tag := range tags {
			tags[i] = strings.ToLower(strings.TrimSpace(tag))
		}
		args = append(args, pq.Array(tags))
		argIndex++
	}

	err := r.db.DB.Get(&count, sql, args...)
	if err != nil {
		return nil, err
	}
	return count, nil
}

func (r *TriviaRepository) GetWrongAnswers(pageSize, offset int, searchString, statusFilter, tagFilter string) ([]*WrongAnswerPoolEntity, error) {
	answers := []*WrongAnswerPoolEntity{}
	sql := `SELECT id, created_at, modified_at, is_archived, answer_text, tags FROM wrong_answer_pool WHERE 1=1`

	// Build dynamic WHERE clause
	args := []interface{}{pageSize, offset}
	argIndex := 3

	// Search filter
	if searchString != "" {
		sql += ` AND (answer_text ILIKE '%' || $` + fmt.Sprintf("%d", argIndex) + ` || '%' OR $` + fmt.Sprintf("%d", argIndex) + ` = ANY(tags))`
		args = append(args, searchString)
		argIndex++
	}

	// Status filter
	if statusFilter != "" {
		switch statusFilter {
		case "active":
			sql += ` AND is_archived = false`
		case "archived":
			sql += ` AND is_archived = true`
		}
	}

	// Tag filter
	if tagFilter != "" {
		sql += ` AND tags && $` + fmt.Sprintf("%d", argIndex)
		// Split comma-separated tags and lowercase them
		tags := strings.Split(tagFilter, ",")
		for i, tag := range tags {
			tags[i] = strings.ToLower(strings.TrimSpace(tag))
		}
		args = append(args, pq.Array(tags))
		argIndex++
	}

	sql += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	err := r.db.DB.Select(&answers, sql, args...)
	if err != nil {
		return nil, err
	}
	return answers, nil
}

func (r *TriviaRepository) GetWrongAnswerById(id int64) (*WrongAnswerPoolEntity, error) {
	answer := &WrongAnswerPoolEntity{}
	sql := `SELECT id, created_at, modified_at, is_archived, answer_text, tags FROM wrong_answer_pool WHERE id = $1`
	err := r.db.DB.Get(answer, sql, id)
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (r *TriviaRepository) CreateWrongAnswer(dto *models.WrongAnswerCreateDTO) (*WrongAnswerPoolEntity, error) {
	// Check for duplicate
	existingAnswer, _ := r.GetWrongAnswerByText(dto.AnswerText)
	if existingAnswer != nil && existingAnswer.ID != 0 {
		return nil, errors.New("wrong answer already exists")
	}

	// Lowercase all tags
	tags := make([]string, len(dto.Tags))
	for i, tag := range dto.Tags {
		tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}

	answer := &WrongAnswerPoolEntity{
		CreatedAt:  time.Now(),
		IsArchived: false,
		AnswerText: dto.AnswerText,
		Tags:       tags,
	}

	sql := `INSERT INTO wrong_answer_pool (answer_text, tags) VALUES ($1, $2) RETURNING id`
	err := r.db.DB.QueryRow(sql, answer.AnswerText, pq.Array(answer.Tags)).Scan(&answer.ID)
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func (r *TriviaRepository) UpdateWrongAnswer(dto *models.WrongAnswerUpdateDTO, id int64) (*WrongAnswerPoolEntity, error) {
	// Check if answer exists
	existingAnswer, err := r.GetWrongAnswerById(id)
	if err != nil {
		return nil, err
	}

	// Check for duplicate if answer text changed
	if existingAnswer.AnswerText != dto.AnswerText {
		duplicate, _ := r.GetWrongAnswerByText(dto.AnswerText)
		if duplicate != nil && duplicate.ID != id {
			return nil, errors.New("wrong answer already exists")
		}
	}

	// Lowercase all tags
	tags := make([]string, len(dto.Tags))
	for i, tag := range dto.Tags {
		tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}

	sql := `UPDATE wrong_answer_pool SET answer_text = $1, tags = $2, modified_at = NOW() WHERE id = $3`
	_, err = r.db.DB.Exec(sql, dto.AnswerText, pq.Array(tags), id)
	if err != nil {
		return nil, err
	}

	return r.GetWrongAnswerById(id)
}

func (r *TriviaRepository) ToggleWrongAnswerArchived(id int64) error {
	sql := `UPDATE wrong_answer_pool SET is_archived = NOT is_archived, modified_at = NOW() WHERE id = $1`
	_, err := r.db.DB.Exec(sql, id)
	return err
}
