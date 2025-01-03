package inbound

import (
	"encoding/json"
	"net/http"
	"post-service/internal/app/service"
	"post-service/internal/model"
	"github.com/gorilla/mux"
)

type StoryHandler struct {
	storyService service.StoryService
}

func NewStoryHandler(storyService service.StoryService) *StoryHandler {
	return &StoryHandler{storyService: storyService}
}

func (handler *StoryHandler) CreateStory(w http.ResponseWriter, r *http.Request) {
	var story model.Story
	err := json.NewDecoder(r.Body).Decode(&story)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.storyService.CreateStory(&story)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *StoryHandler) GetStory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	story, err := handler.storyService.GetStory(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if story == nil {
		http.Error(w, "Story not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(story)
}

func (handler *StoryHandler) UpdateStory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var story model.Story
	err := json.NewDecoder(r.Body).Decode(&story)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	story.ID = id
	err = handler.storyService.UpdateStory(&story)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *StoryHandler) DeleteStory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := handler.storyService.DeleteStory(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
