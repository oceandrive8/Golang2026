package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

const maxTitleLength = 100

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TaskHandler struct {
	mu     sync.Mutex
	tasks  map[int]*Task
	nextID int
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	case http.MethodPatch:
		h.updateTask(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	if idStr == "" {
		list := make([]*Task, 0, len(h.tasks))
		for _, task := range h.tasks {
			list = append(list, task)
		}
		writeJSON(w, http.StatusOK, list)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "id must be a valid integer",
		})
		return
	}

	task, ok := h.tasks[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "task not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title cannot be empty",
		})
		return
	}

	if len(req.Title) > maxTitleLength {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title exceeds maximum length of 100 characters",
		})
		return
	}

	h.mu.Lock()
	task := &Task{
		ID:    h.nextID,
		Title: req.Title,
		Done:  false,
	}
	h.tasks[h.nextID] = task
	h.nextID++
	h.mu.Unlock()

	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "id query parameter is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "id must be a valid integer",
		})
		return
	}

	var req struct {
		Done *bool `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Done == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "done must be a boolean value",
		})
		return
	}

	h.mu.Lock()
	task, ok := h.tasks[id]
	if !ok {
		h.mu.Unlock()
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "task not found",
		})
		return
	}

	task.Done = *req.Done
	h.mu.Unlock()

	writeJSON(w, http.StatusOK, task)
}
