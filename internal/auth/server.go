package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	controld "github.com/baptistecdr/controld-go"

	"github.com/salmonumbrella/controld-cli/internal/secrets"
)

var validAccountName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// clientLimit tracks attempts for a specific client
type clientLimit struct {
	count   int
	resetAt time.Time
}

// rateLimiter tracks attempts per client IP and endpoint to prevent brute-force
type rateLimiter struct {
	mu          sync.Mutex
	attempts    map[string]*clientLimit
	maxAttempts int
	window      time.Duration
}

func newRateLimiter(maxAttempts int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		attempts:    make(map[string]*clientLimit),
		maxAttempts: maxAttempts,
		window:      window,
	}
}

func (rl *rateLimiter) check(clientIP, endpoint string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	key := clientIP + ":" + endpoint
	now := time.Now()

	if limit, exists := rl.attempts[key]; exists && now.After(limit.resetAt) {
		delete(rl.attempts, key)
	}

	if rl.attempts[key] == nil {
		rl.attempts[key] = &clientLimit{
			count:   1,
			resetAt: now.Add(rl.window),
		}
		return nil
	}

	rl.attempts[key].count++
	if rl.attempts[key].count > rl.maxAttempts {
		return fmt.Errorf("too many attempts, please try again later")
	}
	return nil
}

func (rl *rateLimiter) startCleanup(interval time.Duration, stop <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				rl.cleanup()
			case <-stop:
				return
			}
		}
	}()
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, limit := range rl.attempts {
		if now.After(limit.resetAt) {
			delete(rl.attempts, key)
		}
	}
}

func getClientIP(r *http.Request) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

// ValidateAccountName validates an account name
func ValidateAccountName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("account name cannot be empty")
	}
	if len(name) > 64 {
		return fmt.Errorf("account name too long (max 64 characters)")
	}
	if !validAccountName.MatchString(name) {
		return fmt.Errorf("account name contains invalid characters (use only letters, numbers, dash, underscore)")
	}
	return nil
}

// ValidateAPIToken validates an API token
func ValidateAPIToken(token string) error {
	if len(token) == 0 {
		return fmt.Errorf("API token cannot be empty")
	}
	if len(token) > 256 {
		return fmt.Errorf("API token too long (max 256 characters)")
	}
	if !strings.HasPrefix(token, "api.") {
		return fmt.Errorf("API token must start with 'api.'")
	}
	return nil
}

// SetupResult contains the result of a browser-based setup
type SetupResult struct {
	AccountName string
	Error       error
}

// SetupServer handles the browser-based authentication flow
type SetupServer struct {
	result        chan SetupResult
	shutdown      chan struct{}
	stopCleanup   chan struct{}
	pendingResult *SetupResult
	pendingMu     sync.Mutex
	csrfToken     string
	store         secrets.Store
	limiter       *rateLimiter
}

// NewSetupServer creates a new setup server
func NewSetupServer(store secrets.Store) (*SetupServer, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate CSRF token: %w", err)
	}

	stopCleanup := make(chan struct{})
	limiter := newRateLimiter(10, 15*time.Minute)
	limiter.startCleanup(5*time.Minute, stopCleanup)

	return &SetupServer{
		result:      make(chan SetupResult, 1),
		shutdown:    make(chan struct{}),
		stopCleanup: stopCleanup,
		csrfToken:   hex.EncodeToString(tokenBytes),
		store:       store,
		limiter:     limiter,
	}, nil
}

// Start starts the setup server and opens the browser
func (s *SetupServer) Start(ctx context.Context) (*SetupResult, error) {
	defer close(s.stopCleanup)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to start server: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleSetup)
	mux.HandleFunc("/validate", s.handleValidate)
	mux.HandleFunc("/submit", s.handleSubmit)
	mux.HandleFunc("/success", s.handleSuccess)
	mux.HandleFunc("/complete", s.handleComplete)

	server := &http.Server{
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		_ = server.Serve(listener)
	}()

	go func() {
		if err := openBrowser(baseURL); err != nil {
			slog.Info("failed to open browser, user can navigate manually", "url", baseURL)
		}
	}()

	select {
	case result := <-s.result:
		_ = server.Shutdown(context.Background())
		return &result, nil
	case <-ctx.Done():
		_ = server.Shutdown(context.Background())
		return nil, ctx.Err()
	case <-s.shutdown:
		_ = server.Shutdown(context.Background())
		s.pendingMu.Lock()
		defer s.pendingMu.Unlock()
		if s.pendingResult != nil {
			return s.pendingResult, nil
		}
		return nil, fmt.Errorf("setup cancelled")
	}
}

func (s *SetupServer) handleSetup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.New("setup").Parse(setupTemplate)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"CSRFToken": s.csrfToken,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; font-src https://fonts.gstatic.com; connect-src 'self' https://fonts.googleapis.com https://fonts.gstatic.com")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")

	if err := tmpl.Execute(w, data); err != nil {
		slog.Error("setup template execution failed", "error", err)
	}
}

func (s *SetupServer) validateCredentials(ctx context.Context, apiToken string) error {
	if apiToken == "" {
		return fmt.Errorf("API token is required")
	}

	client, err := controld.New(apiToken)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	// Test the connection by listing profiles
	_, err = client.ListProfiles(ctx)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	return nil
}

func (s *SetupServer) handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providedToken := r.Header.Get("X-CSRF-Token")
	if subtle.ConstantTimeCompare([]byte(providedToken), []byte(s.csrfToken)) != 1 {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	clientIP := getClientIP(r)
	if err := s.limiter.check(clientIP, "/validate"); err != nil {
		writeJSON(w, http.StatusTooManyRequests, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var req struct {
		AccountName string `json:"account_name"`
		APIToken    string `json:"api_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	req.AccountName = strings.TrimSpace(req.AccountName)
	req.APIToken = strings.TrimSpace(req.APIToken)

	if err := ValidateAccountName(req.AccountName); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if err := ValidateAPIToken(req.APIToken); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := s.validateCredentials(r.Context(), req.APIToken); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Connection successful!",
	})
}

func (s *SetupServer) handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providedToken := r.Header.Get("X-CSRF-Token")
	if subtle.ConstantTimeCompare([]byte(providedToken), []byte(s.csrfToken)) != 1 {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	clientIP := getClientIP(r)
	if err := s.limiter.check(clientIP, "/submit"); err != nil {
		writeJSON(w, http.StatusTooManyRequests, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var req struct {
		AccountName string `json:"account_name"`
		APIToken    string `json:"api_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	req.AccountName = strings.TrimSpace(req.AccountName)
	req.APIToken = strings.TrimSpace(req.APIToken)

	if err := ValidateAccountName(req.AccountName); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if err := ValidateAPIToken(req.APIToken); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := s.validateCredentials(r.Context(), req.APIToken); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Save to keychain
	err := s.store.Set(req.AccountName, req.APIToken)
	if err != nil {
		slog.Error("failed to save credentials", "error", err)
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   fmt.Sprintf("Failed to save credentials: %v", err),
		})
		return
	}

	s.pendingMu.Lock()
	s.pendingResult = &SetupResult{
		AccountName: req.AccountName,
	}
	s.pendingMu.Unlock()

	writeJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"account_name": req.AccountName,
	})
}

func (s *SetupServer) handleSuccess(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("success").Parse(successTemplate)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	s.pendingMu.Lock()
	accountName := ""
	if s.pendingResult != nil {
		accountName = s.pendingResult.AccountName
	}
	s.pendingMu.Unlock()

	data := map[string]string{
		"AccountName": accountName,
		"CSRFToken":   s.csrfToken,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; font-src https://fonts.gstatic.com; connect-src 'self' https://fonts.googleapis.com https://fonts.gstatic.com")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")

	if err := tmpl.Execute(w, data); err != nil {
		slog.Error("success template execution failed", "error", err)
	}
}

func (s *SetupServer) handleComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providedToken := r.Header.Get("X-CSRF-Token")
	if subtle.ConstantTimeCompare([]byte(providedToken), []byte(s.csrfToken)) != 1 {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	s.pendingMu.Lock()
	if s.pendingResult != nil {
		s.result <- *s.pendingResult
	}
	s.pendingMu.Unlock()
	close(s.shutdown)
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("JSON encoding failed", "error", err)
	}
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
