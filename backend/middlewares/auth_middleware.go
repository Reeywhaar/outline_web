package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/http"
)

type authMiddleware struct {
	authEndpoint string
	PasswordHash string
	tokens       map[string]bool
}

func NewAuthMiddleware(authEndpoint string, password string) *authMiddleware {
	passwordHash := ""

	if password != "" {
		pass_bytes := sha256.Sum256([]byte(password))
		passwordHash = hex.EncodeToString(pass_bytes[:])
	}

	return &authMiddleware{
		authEndpoint: authEndpoint,
		tokens:       make(map[string]bool),
		PasswordHash: passwordHash,
	}
}

func (amw *authMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == amw.authEndpoint {
			amw.handleAuth(w, r)
			return
		}

		if amw.PasswordHash != "" {
			token, err := amw.parseToken(r)
			if err != nil {
				jsonError("Missing auth", w, http.StatusForbidden)
				return
			}

			if !amw.tokenValid(token) {
				jsonError("Forbidden", w, http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (amw *authMiddleware) handleAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jsonError("Method not allowed", w, http.StatusMethodNotAllowed)
		return
	}

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		jsonError("Invalid request", w, http.StatusBadRequest)
		return
	}

	if data["password_hash"] != amw.PasswordHash {
		jsonError("Invalid auth", w, http.StatusForbidden)
		return
	}

	authToken := NewIdGenerator().generate(64)
	amw.tokens[authToken] = true

	http.SetCookie(w, &http.Cookie{
		Name:  "outline__auth",
		Value: authToken,
		Path:  "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}

func (amw *authMiddleware) tokenValid(token string) bool {
	return amw.tokens[token]
}

func (amw *authMiddleware) parseToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("outline__auth")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func jsonError(message string, w http.ResponseWriter, status int) {
	data, _ := json.Marshal(map[string]interface{}{"error": true, "message": message})
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, string(data), status)
}

type idGenerator struct {
	runes []rune
}

func NewIdGenerator() *idGenerator {
	return &idGenerator{
		runes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
	}
}

func (tg *idGenerator) generate(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = tg.runes[rand.Intn(len(tg.runes))]
	}
	return string(b)
}
