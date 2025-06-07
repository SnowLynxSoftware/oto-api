package controllers

import (
	"encoding/json"
	"net/http"

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
	r.Post("/import-questions", c.importTriviaQuestions)
	r.Post("/import-wrong-answers", c.importWrongAnswers)
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
