package web

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/raz0red/radpresence/internal/buildinfo"
	"github.com/raz0red/radpresence/internal/presence"
)

const maxLogLines = 500

// RingBuffer is a thread-safe, fixed-capacity log line buffer implementing io.Writer.
type RingBuffer struct {
	mu    sync.Mutex
	lines []string
}

// Write implements io.Writer. Each call may contain multiple newline-separated lines.
func (r *RingBuffer) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s := strings.TrimRight(string(p), "\n")
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimRight(line, "\r"); line != "" {
			r.lines = append(r.lines, line)
		}
	}
	if len(r.lines) > maxLogLines {
		r.lines = r.lines[len(r.lines)-maxLogLines:]
	}
	return len(p), nil
}

// Lines returns a snapshot copy of all buffered lines.
func (r *RingBuffer) Lines() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, len(r.lines))
	copy(out, r.lines)
	return out
}

// Clear empties the ring buffer.
func (r *RingBuffer) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lines = r.lines[:0]
}

// Hub is shared state between the presence Worker and the web Server.
type Hub struct {
	mu           sync.RWMutex
	status       presence.StatusSnapshot
	lastPollTime time.Time
	Log          *RingBuffer
	shutdownOnce sync.Once
	// PortChange receives a new port number when the user saves a new web_port
	// via the settings page. The Server drains this channel and restarts its
	// HTTP listener on the new port without requiring a full service restart.
	// Closing this channel signals the Server to shut down entirely.
	PortChange chan int
}

// NewHub creates a Hub with an empty ring buffer.
func NewHub() *Hub {
	return &Hub{
		Log:        &RingBuffer{},
		PortChange: make(chan int, 1),
	}
}

// UpdateStatus implements presence.StatusUpdater.
func (h *Hub) UpdateStatus(s presence.StatusSnapshot) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.status = s
	h.lastPollTime = time.Now()
}

// Shutdown gracefully stops the web server by closing the PortChange channel.
// Safe to call more than once.
func (h *Hub) Shutdown() {
	h.shutdownOnce.Do(func() { close(h.PortChange) })
}

// statusResponse is the JSON shape returned by GET /api/status.
type statusResponse struct {
	IsPlaying           bool      `json:"is_playing"`
	GameTitle           string    `json:"game_title"`
	ConsoleName         string    `json:"console_name"`
	CoverArtURL         string    `json:"cover_art_url"`
	RichPresenceMsg     string    `json:"rich_presence_msg"`
	Achievements        string    `json:"achievements"`
	GameStartTime       int64     `json:"game_start_time"`
	DiscordConnected    bool      `json:"discord_connected"`
	RAConnected         bool      `json:"ra_connected"`
	LastPollTime        time.Time `json:"last_poll_time"`
	Username            string    `json:"username"`
	UserProfileURL      string    `json:"user_profile_url"`
	GamePageURL         string    `json:"game_page_url"`
	AvatarURL           string    `json:"avatar_url"`
	TotalPoints         int       `json:"total_points"`
	TotalSoftcorePoints int       `json:"total_softcore_points"`
	TotalTruePoints     int       `json:"total_true_points"`
	Rank                int       `json:"rank"`
	TotalRanked         int       `json:"total_ranked"`
	LastError           string    `json:"last_error"`
	Version             string    `json:"version"`
}

func (h *Hub) getStatusResponse() statusResponse {
	h.mu.RLock()
	defer h.mu.RUnlock()
	resp := statusResponse{
		IsPlaying:           h.status.IsPlaying,
		GameTitle:           h.status.GameTitle,
		ConsoleName:         h.status.ConsoleName,
		CoverArtURL:         h.status.CoverArtURL,
		RichPresenceMsg:     h.status.RichPresenceMsg,
		Achievements:        h.status.Achievements,
		GameStartTime:       h.status.GameStartTime,
		DiscordConnected:    h.status.DiscordConnected,
		RAConnected:         h.status.RAConnected,
		LastPollTime:        h.lastPollTime,
		Username:            h.status.Username,
		AvatarURL:           h.status.AvatarURL,
		TotalPoints:         h.status.TotalPoints,
		TotalSoftcorePoints: h.status.TotalSoftcorePoints,
		TotalTruePoints:     h.status.TotalTruePoints,
		Rank:                h.status.Rank,
		TotalRanked:         h.status.TotalRanked,
		LastError:           h.status.LastError,
	}
	if h.status.Username != "" {
		resp.UserProfileURL = fmt.Sprintf("https://retroachievements.org/user/%s", h.status.Username)
	}
	if h.status.GameID != 0 {
		resp.GamePageURL = fmt.Sprintf("https://retroachievements.org/game/%d", h.status.GameID)
	}
	resp.Version = buildinfo.Version
	return resp
}
