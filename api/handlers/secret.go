package secret

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

type Secret struct {
	ID         string    `json:"id"`
	SecretName string    `json:"secret_name"`
	Key        string    `json:"key"`
	Value      string    `json:"value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// curl -X PUT http://localhost:8080/secret/secretName -H "Content-Type: application/json" -H "Authorization: token" -d '{"key1": "value1", "key2": "value2"}'
func (h *Handler) AddSecret(w http.ResponseWriter, r *http.Request) {
	var secrets map[string]string
	if err := json.NewDecoder(r.Body).Decode(&secrets); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	secretName := strings.TrimPrefix(r.URL.Path, "/api/v1/secret/")
	for key, value := range secrets {
		var id string
		err := tx.QueryRow(`SELECT id FROM secrets WHERE secret_name = ? AND key = ?`, secretName, key).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				id = ulid.MustNew(ulid.Now(), nil).String()
				_, err = tx.Exec(`INSERT INTO secrets (id, secret_name, key, value, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`, id, secretName, key, value, time.Now(), time.Now())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				continue
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec(`UPDATE secrets SET value = ?, updated_at = ? WHERE id = ?`, value, time.Now(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// curl -X GET http://localhost:8080/secrets?offset=0&limit=10 -H "Content-Type: application/json" -H "Authorization: token"
func (h *Handler) GetAllSecrets(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	secretName := r.URL.Query().Get("secret_name")

	var rows *sql.Rows
	if secretName != "" {
		rows, err = h.db.Query(`SELECT id, secret_name, key, value, created_at, updated_at FROM secrets WHERE secret_name = ? ORDER BY secret_name LIMIT ? OFFSET ?`, secretName, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		rows, err = h.db.Query(`SELECT id, secret_name, key, value, created_at, updated_at FROM secrets ORDER BY secret_name LIMIT ? OFFSET ?`, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var secrets []Secret
	for rows.Next() {
		var s Secret
		if err := rows.Scan(&s.ID, &s.SecretName, &s.Key, &s.Value, &s.CreatedAt, &s.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		secrets = append(secrets, s)
	}

	if err := json.NewEncoder(w).Encode(secrets); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// curl -X GET http://localhost:8080/secret/secretName -H "Content-Type: application/json" -H "Authorization: token"
func (h *Handler) GetSecret(w http.ResponseWriter, r *http.Request) {
	secretName := strings.TrimPrefix(r.URL.Path, "/api/v1/secret/")
	rows, err := h.db.Query(`SELECT id, secret_name, key, value, created_at, updated_at FROM secrets WHERE secret_name = ?`, secretName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var secrets []Secret
	for rows.Next() {
		var s Secret
		if err := rows.Scan(&s.ID, &s.SecretName, &s.Key, &s.Value, &s.CreatedAt, &s.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		secrets = append(secrets, s)
	}

	if err := json.NewEncoder(w).Encode(secrets); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// curl -X DELETE http://localhost:8080/secret/secretName?key=key -H "Content-Type: application/json" -H "Authorization: token"
func (h *Handler) DeleteSecretKeyAndValue(w http.ResponseWriter, r *http.Request) {
	secretName := strings.TrimPrefix(r.URL.Path, "/api/v1/secret/")
	key := r.URL.Query().Get("key")

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := tx.Exec(`DELETE FROM secrets WHERE secret_name = ? AND key = ?`, secretName, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		_ = tx.Rollback()
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		_ = tx.Rollback()
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows affected", http.StatusNotFound)
		_ = tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// curl -X POST http://localhost:8080/secret/secretName -H "Content-Type: application/json" -H "Authorization: token"
func (h *Handler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	secretName := strings.TrimPrefix(r.URL.Path, "/api/v1/secret/")

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`DELETE FROM secrets WHERE secret_name = ?`, secretName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
