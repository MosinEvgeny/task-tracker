package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MosinEvgeny/task-tracker/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// LabelHandler обрабатывает HTTP-запросы для работы с метками.
type LabelHandler struct {
	labelService service.LabelService
}

// NewLabelHandler создает новый экземпляр LabelHandler.
func NewLabelHandler(labelService service.LabelService) *LabelHandler {
	return &LabelHandler{labelService: labelService}
}

func (h *LabelHandler) CreateLabel(w http.ResponseWriter, r *http.Request) {
	var labelData struct {
		Name   string    `json:"name"`
		Color  string    `json:"color"`
		UserID uuid.UUID `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&labelData); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	createdLabel, err := h.labelService.CreateLabel(r.Context(), labelData.Name, labelData.Color, labelData.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdLabel)
}

func (h *LabelHandler) GetLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID метки", http.StatusBadRequest)
		return
	}

	label, err := h.labelService.GetLabelByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Метка не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(label)
}

func (h *LabelHandler) UpdateLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID метки", http.StatusBadRequest)
		return
	}

	var labelData struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&labelData); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	updatedLabel, err := h.labelService.UpdateLabel(r.Context(), id, labelData.Name, labelData.Color)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedLabel)
}

func (h *LabelHandler) DeleteLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID метки", http.StatusBadRequest)
		return
	}

	err = h.labelService.DeleteLabel(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
