package web

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/raz0red/radpresence/internal/buildinfo"
	"github.com/raz0red/radpresence/internal/config"
	"github.com/raz0red/radpresence/internal/ra"
)

//go:embed templates static
var templateFS embed.FS

// Server is the optional web UI HTTP server.
type Server struct {
	hub       *Hub
	port      int
	csrfToken string
	mu        sync.Mutex
	httpSrv   *http.Server
}

// NewServer creates a Server bound to the given Hub and port.
func NewServer(hub *Hub, port int) *Server {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return &Server{
		hub:       hub,
		port:      port,
		csrfToken: hex.EncodeToString(b),
	}
}

// Start begins listening on 127.0.0.1:port. Blocks until the hub is closed.
// Automatically restarts on a new port when hub.PortChange fires.
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/favicon.webp", s.handleFavicon)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/favicon.webp", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/api/status", s.handleAPIStatus)
	mux.HandleFunc("/api/logs", s.handleAPILogs)
	mux.HandleFunc("/api/logs/clear", s.handleAPILogsClear)
	mux.HandleFunc("/settings", s.handleSettings)
	mux.HandleFunc("/switching", s.handleSwitching)

	s.startListener(s.port, mux)

	for newPort := range s.hub.PortChange {
		log.Printf("[web] restarting on port %d", newPort)
		s.mu.Lock()
		srv := s.httpSrv
		s.mu.Unlock()
		if srv != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_ = srv.Shutdown(ctx)
			cancel()
		}
		s.mu.Lock()
		s.port = newPort
		s.mu.Unlock()
		s.startListener(newPort, mux)
	}

	// PortChange was closed — shut down the current listener and exit.
	log.Println("[web] shutting down")
	s.mu.Lock()
	srv := s.httpSrv
	s.mu.Unlock()
	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = srv.Shutdown(ctx)
		cancel()
	}
	return nil
}

// startListener launches an http.Server in a background goroutine.
func (s *Server) startListener(port int, mux http.Handler) {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	srv := &http.Server{Addr: addr, Handler: mux}
	s.mu.Lock()
	s.httpSrv = srv
	s.mu.Unlock()
	log.Printf("[web] UI available at http://%s", addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("[web] listener error: %v", err)
		}
	}()
}

type indexData struct {
	Version          string
	Username         string
	HasAPIKey        bool
	Interval         int
	HideButtons      bool
	HideAchievements bool
	WebPort          int
	CSRFToken        string
	SavedMsg         string
	ErrorMsg         string
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	cfg, _ := config.Load()
	interval := cfg.Interval
	if interval == 0 {
		interval = 10
	}
	webPort := cfg.WebPort
	if webPort == 0 {
		webPort = s.port
	}

	t, err := template.ParseFS(templateFS, "templates/index.html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		log.Printf("[web] template parse error: %v", err)
		return
	}

	data := indexData{
		Version:          buildinfo.Version,
		Username:         cfg.Username,
		HasAPIKey:        cfg.APIKey != "",
		Interval:         interval,
		HideButtons:      cfg.HideButtons,
		HideAchievements: cfg.HideAchievements,
		WebPort:          webPort,
		CSRFToken:        s.csrfToken,
		SavedMsg:         r.URL.Query().Get("saved"),
		ErrorMsg:         r.URL.Query().Get("error"),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		log.Printf("[web] template execute error: %v", err)
	}
}

func (s *Server) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.hub.getStatusResponse()); err != nil {
		log.Printf("[web] status encode error: %v", err)
	}
}

func (s *Server) handleAPILogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.hub.Log.Lines()); err != nil {
		log.Printf("[web] logs encode error: %v", err)
	}
}

func (s *Server) handleFavicon(w http.ResponseWriter, r *http.Request) {
	data, err := templateFS.ReadFile("static/favicon.webp")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/webp")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = w.Write(data)
}

func (s *Server) handleAPILogsClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	s.hub.Log.Clear()
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// CSRF check — mitigates cross-origin form submission attacks.
	if r.FormValue("csrf_token") != s.csrfToken {
		http.Error(w, "invalid request", http.StatusForbidden)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		http.Error(w, "failed to load config", http.StatusInternalServerError)
		return
	}

	existingCfg, _ := config.Load()

	if v := r.FormValue("username"); v != "" {
		cfg.Username = v
	}
	// Only update API key if user typed a new one.
	if v := r.FormValue("api_key"); v != "" {
		cfg.APIKey = v
	}
	if v := r.FormValue("interval"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 5 {
			cfg.Interval = n
		}
	}
	cfg.HideButtons = r.FormValue("hide_buttons") == "on"
	cfg.HideAchievements = r.FormValue("hide_achievements") == "on"

	var portChanged bool
	if v := r.FormValue("web_port"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 1024 && n <= 65535 {
			s.mu.Lock()
			portChanged = n != s.port
			s.mu.Unlock()
			cfg.WebPort = n
		}
	}

	// Only validate credentials against the RA API when they actually changed.
	newUsername := r.FormValue("username")
	newAPIKey := r.FormValue("api_key")
	credentialsChanged := (newUsername != "" && newUsername != existingCfg.Username) ||
		newAPIKey != "" // api_key field is blank unless user typed something
	if credentialsChanged && cfg.Username != "" && cfg.APIKey != "" {
		client := ra.New(cfg.Username, cfg.APIKey)
		if _, err := client.GetUserSummary(); err != nil {
			log.Printf("[web] credential validation failed: %v", err)
			// Revert credential fields to existing values so other changes are not lost.
			cfg.Username = existingCfg.Username
			cfg.APIKey = existingCfg.APIKey
			// Save the non-credential changes (interval, port, etc.) before redirecting.
			_ = config.Save(cfg)
			http.Redirect(w, r, "/?tab=settings&error=invalid_credentials", http.StatusSeeOther)
			return
		}
	}

	if err := config.Save(cfg); err != nil {
		http.Error(w, "failed to save config", http.StatusInternalServerError)
		return
	}

	// Signal a port restart after a short delay so the browser receives the
	// /switching redirect and loads the countdown page before the old listener
	// closes. 2 s gives plenty of headroom; the countdown page waits 5 s.
	if portChanged {
		newPort := cfg.WebPort
		go func() {
			time.Sleep(2 * time.Second)
			select {
			case s.hub.PortChange <- newPort:
			default:
			}
		}()
		http.Redirect(w, r, fmt.Sprintf("/switching?port=%d", newPort), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/?tab=settings&saved=1", http.StatusSeeOther)
}

const switchingHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>RAD Presence — Switching Port</title>
  <style>
    body { background:#1a1a2e; color:#e0e0f0; font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;
           display:flex; align-items:center; justify-content:center; min-height:100vh; margin:0; }
    .box { text-align:center; }
    h2  { color:#5f9ef2; margin-bottom:0.6rem; }
    p   { color:#6a6a8a; font-size:0.9rem; margin-bottom:1.5rem; }
    .countdown { font-size:3.5rem; font-weight:700; color:#f0c040; }
    .url  { font-size:0.82rem; color:#6a6a8a; margin-top:1rem; }
    .url a { color:#5f9ef2; }
  </style>
</head>
<body>
<div class="box">
  <h2>Switching port&hellip;</h2>
  <p>The web UI is restarting on the new port.</p>
  <div class="countdown" id="c">5</div>
  <div class="url">Redirecting to <a id="link" href=""></a></div>
</div>
<script>
  var port = %d;
  var url  = 'http://127.0.0.1:' + port + '/?tab=settings&saved=1';
  document.getElementById('link').textContent = url;
  document.getElementById('link').href = url;
  var n = 5;
  var iv = setInterval(function() {
    n--;
    document.getElementById('c').textContent = n;
    if (n <= 0) { clearInterval(iv); window.location.href = url; }
  }, 1000);
</script>
</body>
</html>`

func (s *Server) handleSwitching(w http.ResponseWriter, r *http.Request) {
	portStr := r.URL.Query().Get("port")
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1024 || port > 65535 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, switchingHTML, port)
}
