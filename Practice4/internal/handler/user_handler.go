package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"awesomeProject/internal/usecase"
	"awesomeProject/pkg/modules"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

// ------------------ GET /users ------------------
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	users, err := h.usecase.GetAllUsers(limit, offset)
	if err != nil {
		log.Println("GetAllUsers error:", err)
		http.Error(w, `{"error":"failed to get users"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("JSON encode error:", err)
		http.Error(w, `{"error":"failed to encode users"}`, http.StatusInternalServerError)
		return
	}
}

// ------------------ GET /users/{id} ------------------
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid user id"}`, http.StatusBadRequest)
		return
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ------------------ POST /users ------------------
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user modules.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error":"invalid request payload"}`, http.StatusBadRequest)
		return
	}

	id, err := h.usecase.CreateUser(&user)
	if err != nil {
		http.Error(w, `{"error":"failed to create user"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// ------------------ PATCH /users/{id} ------------------
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid user id"}`, http.StatusBadRequest)
		return
	}

	var user modules.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error":"invalid request payload"}`, http.StatusBadRequest)
		return
	}

	user.ID = id
	if err := h.usecase.UpdateUser(&user); err != nil {
		log.Println("UPDATE ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ------------------ DELETE /users/{id} ------------------
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid user id"}`, http.StatusBadRequest)
		return
	}

	rows, err := h.usecase.DeleteUser(id)
	if err != nil || rows == 0 {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
