package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/snowlynxsoftware/oto-api/server/models"
	"github.com/snowlynxsoftware/oto-api/server/services"
)

type WaitlistController struct {
	waitlistService services.IWaitlistService
}

func NewWaitlistController(waitlistService services.IWaitlistService) *WaitlistController {
	return &WaitlistController{
		waitlistService: waitlistService,
	}
}

func (s *WaitlistController) MapController() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", s.createWaitlistEntry)
	return r
}

func (c *WaitlistController) createWaitlistEntry(w http.ResponseWriter, r *http.Request) {
	var waitlistCreateDTO models.WaitlistCreateDTO

	err := json.NewDecoder(r.Body).Decode(&waitlistCreateDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.waitlistService.CreateNewWaitlistEntry(&waitlistCreateDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user added to waitlist successfully"))
}
