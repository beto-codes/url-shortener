package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/beto-codes/url-shortener/internal/shortener"
	"github.com/beto-codes/url-shortener/internal/storage"
	"github.com/beto-codes/url-shortener/internal/utils/constants"
)

type Handler struct {
	service *shortener.Service
	baseURL string
}

func NewHandler(service *shortener.Service, baseURL string) *Handler {
	return &Handler{
		service: service,
		baseURL: baseURL,
	}
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.jsonError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.jsonError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if req.URL == constants.EmptyString {
		h.jsonError(w, "URL is required", http.StatusBadRequest)
		return
	}

	err = validateURL(&req.URL)
	if err != nil {
		h.jsonError(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	shortCode, err := h.service.Shorten(req.URL)
	if err != nil {
		h.jsonError(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	resp := ShortenResponse{
		ShortCode: shortCode,
		ShortURL:  h.baseURL + "/" + shortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.jsonError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == constants.EmptyString || shortCode == "health" || shortCode == "shorten" {
		http.NotFound(w, r)
		return
	}

	if !isValidShortCode(shortCode) {
		h.jsonError(w, "Invalid Short code format", http.StatusBadRequest)
		return
	}

	longUrl, err := h.service.Resolve(shortCode)
	if err != nil {
		if err == storage.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		h.jsonError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, longUrl, http.StatusMovedPermanently)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (h *Handler) jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func validateURL(urlStr *string) error {
	if !strings.HasPrefix(*urlStr, "http://") && !strings.HasPrefix(*urlStr, "https://") {
		*urlStr = "http://" + *urlStr
	}

	parsedURL, err := url.Parse(*urlStr)
	if err != nil {
		return err
	}

	if parsedURL.Scheme == constants.EmptyString || parsedURL.Host == constants.EmptyString {
		return storage.ErrEmpty
	}

	return nil
}

func isValidShortCode(code string) bool {
	if len(code) == 0 {
		return false
	}

	for _, ch := range code {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '_') {
			return false
		}
	}

	return true
}
