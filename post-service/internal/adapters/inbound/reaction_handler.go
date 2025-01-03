package inbound

import (
	"encoding/json"
	"net/http"
	"post-service/internal/app/service"
	"post-service/internal/model"
	"github.com/gorilla/mux"
)

type ReactionHandler struct {
	reactionService service.ReactionService
}

func NewReactionHandler(reactionService service.ReactionService) *ReactionHandler {
	return &ReactionHandler{reactionService: reactionService}
}

func (handler *ReactionHandler) CreateReaction(w http.ResponseWriter, r *http.Request) {
	var reaction model.Reaction
	err := json.NewDecoder(r.Body).Decode(&reaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.reactionService.CreateReaction(&reaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *ReactionHandler) GetReaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	reaction, err := handler.reactionService.GetReaction(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if reaction == nil {
		http.Error(w, "Reaction not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(reaction)
}

func (handler *ReactionHandler) UpdateReaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var reaction model.Reaction
	err := json.NewDecoder(r.Body).Decode(&reaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reaction.ID = id

	err = handler.reactionService.UpdateReaction(&reaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ReactionHandler) DeleteReaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := handler.reactionService.DeleteReaction(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
