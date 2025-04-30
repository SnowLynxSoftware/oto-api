package models

type TriviaQuestionImportData struct {
	Question      string   `json:"question"`
	CorrectAnswer string   `json:"correct_answer"`
	Tags          []string `json:"tags"`
}

type TriviaQuestionImportResults struct {
	TotalQuestionsProcessed int64 `json:"total_questions_processed"`
	QuestionsAdded          int64 `json:"questions_added"`
}
