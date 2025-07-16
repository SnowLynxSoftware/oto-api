package models

type TriviaQuestionImportData struct {
	Question      string   `json:"question"`
	CorrectAnswer string   `json:"correct_answer"`
	Tags          []string `json:"tags"`
}

type TriviaWrongAnswerImportData struct {
	Question   string   `json:"question"`
	AnswerText string   `json:"answer_text"`
	Tags       []string `json:"tags"`
}

type TriviaQuestionImportResults struct {
	TotalQuestionsProcessed int64 `json:"total_questions_processed"`
	QuestionsAdded          int64 `json:"questions_added"`
}

type TriviaWrongAnswerImportResults struct {
	TotalAnswersProcessed int64 `json:"total_answers_processed"`
	AnswersAdded          int64 `json:"answers_added"`
}

type TriviaDeckMetadataUpdateRequest struct {
	DeckID      int64  `json:"deck_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DTOs for CRUD operations
type TriviaQuestionCreateDTO struct {
	Question      string   `json:"question"`
	CorrectAnswer string   `json:"correct_answer"`
	Tags          []string `json:"tags"`
	IsPublished   bool     `json:"is_published"`
}

type TriviaQuestionUpdateDTO struct {
	Question      string   `json:"question"`
	CorrectAnswer string   `json:"correct_answer"`
	Tags          []string `json:"tags"`
	IsPublished   bool     `json:"is_published"`
}

type WrongAnswerCreateDTO struct {
	AnswerText string   `json:"answer_text"`
	Tags       []string `json:"tags"`
}

type WrongAnswerUpdateDTO struct {
	AnswerText string   `json:"answer_text"`
	Tags       []string `json:"tags"`
}
