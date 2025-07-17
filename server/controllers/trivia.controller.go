package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/snowlynxsoftware/oto-api/server/middleware"
	"github.com/snowlynxsoftware/oto-api/server/models"
	"github.com/snowlynxsoftware/oto-api/server/services"
	"github.com/snowlynxsoftware/oto-api/server/util"
)

type TriviaController struct {
	triviaService  services.ITriviaService
	authMiddleware middleware.IAuthMiddleware
}

func NewTriviaController(triviaService services.ITriviaService, authMiddleware middleware.IAuthMiddleware) *TriviaController {
	return &TriviaController{
		triviaService:  triviaService,
		authMiddleware: authMiddleware,
	}
}

func (c *TriviaController) MapController() *chi.Mux {
	r := chi.NewRouter()

	// Existing endpoints
	r.Post("/import-questions", c.importTriviaQuestions)
	r.Post("/import-wrong-answers", c.importWrongAnswers)

	// New CRUD endpoints for questions
	r.Get("/questions", c.getQuestions)
	r.Get("/questions/{id}", c.getQuestionById)
	r.Post("/questions", c.createQuestion)
	r.Put("/questions/{id}", c.updateQuestion)
	r.Patch("/questions/{id}/archived", c.toggleQuestionArchived)
	r.Patch("/questions/{id}/published", c.toggleQuestionPublished)

	// New CRUD endpoints for wrong answers
	r.Get("/wrong-answers", c.getWrongAnswers)
	r.Get("/wrong-answers/{id}", c.getWrongAnswerById)
	r.Post("/wrong-answers", c.createWrongAnswer)
	r.Put("/wrong-answers/{id}", c.updateWrongAnswer)
	r.Patch("/wrong-answers/{id}/archived", c.toggleWrongAnswerArchived)

	return r
}

func (c *TriviaController) importTriviaQuestions(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	var importData []models.TriviaQuestionImportData

	err = json.NewDecoder(r.Body).Decode(&importData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := c.triviaService.ImportTriviaQuestions(importData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnStr, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(returnStr)
}

func (c *TriviaController) importWrongAnswers(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	var importData []models.TriviaWrongAnswerImportData

	err = json.NewDecoder(r.Body).Decode(&importData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := c.triviaService.ImportWrongAnswers(importData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnStr, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(returnStr)
}

// Question CRUD endpoints
func (c *TriviaController) getQuestions(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	// Get query parameters
	pageSize := 25
	page := 1
	searchString := ""
	statusFilter := ""
	tagFilter := ""

	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if psInt, err := strconv.Atoi(ps); err == nil && psInt > 0 {
			pageSize = psInt
		}
	}
	if p := r.URL.Query().Get("page"); p != "" {
		if pInt, err := strconv.Atoi(p); err == nil && pInt > 0 {
			page = pInt
		}
	}
	if search := r.URL.Query().Get("search"); search != "" {
		searchString = search
	}
	if status := r.URL.Query().Get("status"); status != "" {
		statusFilter = status
	}
	if tags := r.URL.Query().Get("tags"); tags != "" {
		tagFilter = tags
	}

	offset := (page - 1) * pageSize

	results, err := c.triviaService.GetQuestions(pageSize, offset, searchString, statusFilter, tagFilter)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to retrieve questions", http.StatusInternalServerError)
		return
	}

	// Fix page number
	results.Page = page

	returnStr, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnStr)
}

func (c *TriviaController) getQuestionById(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid question ID", http.StatusBadRequest)
		return
	}

	question, err := c.triviaService.GetQuestionById(id)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to retrieve question", http.StatusInternalServerError)
		return
	}

	returnStr, err := json.Marshal(question)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnStr)
}

func (c *TriviaController) createQuestion(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	var createDTO models.TriviaQuestionCreateDTO
	err = json.NewDecoder(r.Body).Decode(&createDTO)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	question, err := c.triviaService.CreateQuestion(&createDTO)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to create question", http.StatusInternalServerError)
		return
	}

	returnStr, err := json.Marshal(question)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(returnStr)
}

func (c *TriviaController) updateQuestion(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid question ID", http.StatusBadRequest)
		return
	}

	var updateDTO models.TriviaQuestionUpdateDTO
	err = json.NewDecoder(r.Body).Decode(&updateDTO)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	question, err := c.triviaService.UpdateQuestion(&updateDTO, id)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to update question", http.StatusInternalServerError)
		return
	}

	returnStr, err := json.Marshal(question)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnStr)
}

func (c *TriviaController) toggleQuestionArchived(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid question ID", http.StatusBadRequest)
		return
	}

	err = c.triviaService.ToggleQuestionArchived(id)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to toggle question archived status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("question archived status toggled successfully"))
}

func (c *TriviaController) toggleQuestionPublished(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid question ID", http.StatusBadRequest)
		return
	}

	err = c.triviaService.ToggleQuestionPublished(id)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to toggle question published status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("question published status toggled successfully"))
}

// Wrong Answer CRUD endpoints
func (c *TriviaController) getWrongAnswers(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	// Get query parameters
	pageSize := 25
	page := 1
	searchString := ""
	statusFilter := ""
	tagFilter := ""

	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if psInt, err := strconv.Atoi(ps); err == nil && psInt > 0 {
			pageSize = psInt
		}
	}
	if p := r.URL.Query().Get("page"); p != "" {
		if pInt, err := strconv.Atoi(p); err == nil && pInt > 0 {
			page = pInt
		}
	}
	if search := r.URL.Query().Get("search"); search != "" {
		searchString = search
	}
	if status := r.URL.Query().Get("status"); status != "" {
		statusFilter = status
	}
	if tags := r.URL.Query().Get("tags"); tags != "" {
		tagFilter = tags
	}

	offset := (page - 1) * pageSize

	results, err := c.triviaService.GetWrongAnswers(pageSize, offset, searchString, statusFilter, tagFilter)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to retrieve wrong answers", http.StatusInternalServerError)
		return
	}

	// Fix page number
	results.Page = page

	returnStr, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnStr)
}

func (c *TriviaController) getWrongAnswerById(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid wrong answer ID", http.StatusBadRequest)
		return
	}

	wrongAnswer, err := c.triviaService.GetWrongAnswerById(id)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to retrieve wrong answer", http.StatusInternalServerError)
		return
	}

	returnStr, err := json.Marshal(wrongAnswer)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnStr)
}

func (c *TriviaController) createWrongAnswer(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	var createDTO models.WrongAnswerCreateDTO
	err = json.NewDecoder(r.Body).Decode(&createDTO)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	wrongAnswer, err := c.triviaService.CreateWrongAnswer(&createDTO)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to create wrong answer", http.StatusInternalServerError)
		return
	}

	returnStr, err := json.Marshal(wrongAnswer)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(returnStr)
}

func (c *TriviaController) updateWrongAnswer(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid wrong answer ID", http.StatusBadRequest)
		return
	}

	var updateDTO models.WrongAnswerUpdateDTO
	err = json.NewDecoder(r.Body).Decode(&updateDTO)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	wrongAnswer, err := c.triviaService.UpdateWrongAnswer(&updateDTO, id)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to update wrong answer", http.StatusInternalServerError)
		return
	}

	returnStr, err := json.Marshal(wrongAnswer)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnStr)
}

func (c *TriviaController) toggleWrongAnswerArchived(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r, []string{"admin"})
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid wrong answer ID", http.StatusBadRequest)
		return
	}

	err = c.triviaService.ToggleWrongAnswerArchived(id)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to toggle wrong answer archived status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("wrong answer archived status toggled successfully"))
}
