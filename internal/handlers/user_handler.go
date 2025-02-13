package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MosinEvgeny/task-tracker/internal/auth"
	"github.com/MosinEvgeny/task-tracker/internal/config"
	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/MosinEvgeny/task-tracker/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService         service.UserService
	refreshTokenService service.RefreshTokenService
	config              config.Config
}

func NewUserHandler(userService service.UserService, refreshTokenService service.RefreshTokenService, config config.Config) *UserHandler {
	return &UserHandler{userService: userService, refreshTokenService: refreshTokenService, config: config}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	createdUser, err := h.userService.CreateUser(r.Context(), user.Username, user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByEmail(r.Context(), loginData.Email)
	if err != nil {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	if err := user.ComparePassword(loginData.Password); err != nil {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.refreshTokenService.CreateRefreshToken(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Ошибка при создании refresh токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString, "refresh_token": refreshToken.Token})
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshTokenData struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&refreshTokenData); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	refreshToken, err := h.refreshTokenService.GetRefreshToken(r.Context(), refreshTokenData.RefreshToken)
	if err != nil {
		http.Error(w, "Неверный refresh токен", http.StatusUnauthorized)
		return
	}

	if refreshToken.ExpiryDate.Before(time.Now().UTC()) {
		http.Error(w, "Срок действия refresh токена истек", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": refreshToken.UserID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (h *UserHandler) RevokeAllRefreshTokens(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromRequest(r)
	if !ok {
		http.Error(w, "Не удалось получить ID пользователя из контекста", http.StatusInternalServerError)
		return
	}

	err := h.refreshTokenService.DeleteAllRefreshTokensByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Ошибка при отзыве refresh токенов", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	updatedUser, err := h.userService.UpdateUser(r.Context(), id, user.Username, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	err = h.userService.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetUserIDFromRequest(r *http.Request) (uuid.UUID, bool) {
	return auth.UserIDFromContext(r.Context())
}
