package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"

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
	// Optional limit/offset, default to fetch everything
	limit := 1000 // or any large number
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

	orderByStr := r.URL.Query().Get("order_by")
	orderBy := []string{}
	if orderByStr != "" {
		for _, part := range strings.Split(orderByStr, ",") {
			orderBy = append(orderBy, strings.TrimSpace(part))
		}
	}

	users, err := h.usecase.GetAllUsers(limit, offset, orderBy)
	if err != nil {
		log.Println("GetAllUsers error:", err)
		http.Error(w, `{"error":"failed to get users"}`, http.StatusInternalServerError)
		return
	}

	resp := struct {
		Data []modules.User `json:"data"`
	}{
		Data: users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ------------------ GET /users/{id} ------------------
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")

	id, err := uuid.Parse(idStr)
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

	json.NewEncoder(w).Encode(map[string]string{
		"id": id.String(),
	})
}

// ------------------ PATCH /users/{id} ------------------
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")

	id, err := uuid.Parse(idStr)
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

	id, err := uuid.Parse(idStr)
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

// ------------------ GET /users/paginated ------------------
// Fetch users with pagination, filters, and sorting
func (h *UserHandler) GetUsersPaginated(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 10

	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	if ps := r.URL.Query().Get("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	// ------------------ Filtering ------------------
	filters := map[string]interface{}{}
	for _, field := range []string{"id", "name", "email", "gender", "birthday"} {
		if val := r.URL.Query().Get(field); val != "" {
			if field == "id" {
				if id, err := uuid.Parse(val); err == nil {
					filters[field] = id
				}
			} else {
				filters[field] = val
			}
		}
	}

	// ------------------ Sorting ------------------
	orderByStr := r.URL.Query().Get("order_by")
	orderBy := []string{}
	if orderByStr != "" {
		for _, part := range strings.Split(orderByStr, ",") {
			orderBy = append(orderBy, strings.TrimSpace(part))
		}
	}

	resp, err := h.usecase.GetPaginatedUsers(page, pageSize, filters, orderBy)
	if err != nil {
		log.Println("GetUsersPaginated error:", err)
		http.Error(w, `{"error":"failed to get users"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func (h *UserHandler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	// ------------------ Query Params ------------------
	user1Str := r.URL.Query().Get("user1")
	user2Str := r.URL.Query().Get("user2")

	if user1Str == "" || user2Str == "" {
		http.Error(w, `{"error":"user1 and user2 parameters required"}`, http.StatusBadRequest)
		return
	}

	user1, err := uuid.Parse(user1Str)
	if err != nil {
		http.Error(w, `{"error":"invalid user1 UUID"}`, http.StatusBadRequest)
		return
	}

	user2, err := uuid.Parse(user2Str)
	if err != nil {
		http.Error(w, `{"error":"invalid user2 UUID"}`, http.StatusBadRequest)
		return
	}

	// ------------------ Call usecase ------------------
	friends, err := h.usecase.GetCommonFriends(r.Context(), user1, user2)
	if err != nil {
		log.Println("GetCommonFriends error:", err)
		http.Error(w, `{"error":"failed to fetch common friends"}`, http.StatusInternalServerError)
		return
	}

	// ------------------ Response ------------------
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}
