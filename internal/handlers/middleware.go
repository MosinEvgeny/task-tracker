package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MosinEvgeny/task-tracker/internal/auth"
	"github.com/MosinEvgeny/task-tracker/internal/config"
	"github.com/MosinEvgeny/task-tracker/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Middleware описывает интерфейс middleware.
type Middleware interface {
	Authenticate(next http.Handler) http.Handler
}

type AuthMiddleware struct {
	userService service.UserService
	config      config.Config
}

func NewAuthMiddleware(userService service.UserService, config config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
		config:      config,
	}
}

// Authenticate проверяет JWT токен в заголовке Authorization и добавляет ID пользователя в контекст.
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Получение токена из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Отсутствует заголовок Authorization", http.StatusUnauthorized)
			return
		}

		// 2. Проверка формата заголовка (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Неверный формат заголовка Authorization", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// 3. Валидация токена
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем алгоритм подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неверный алгоритм подписи: %v", token.Header["alg"])
			}

			// Возвращаем секретный ключ
			return []byte(m.config.JWTSecret), nil
		})

		if err != nil {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		// 4. Извлечение ID пользователя из токена
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userIDString, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "Неверный формат ID пользователя в токене", http.StatusInternalServerError)
				return
			}

			userID, err := uuid.Parse(userIDString)
			if err != nil {
				http.Error(w, "Неверный ID пользователя в токене", http.StatusInternalServerError)
				return
			}

			// 5. Добавление ID пользователя в контекст
			ctx := auth.ContextWithUser(r.Context(), userID)

			// 6. Передача управления следующему обработчику
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
		}
	})
}

// Log запросов (middleware).
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(startTime),
		)
	})
}
