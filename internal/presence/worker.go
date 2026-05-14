// Package presence implements the RetroAchievements → Discord Rich Presence poll loop.
package presence

import (
	"fmt"
	"log"
	"time"

	"github.com/raz0red/radpresence/internal/config"
	"github.com/raz0red/radpresence/internal/discord"
	"github.com/raz0red/radpresence/internal/ra"
)

// Worker polls RetroAchievements and keeps Discord Rich Presence up to date.
type Worker struct {
	ra               *ra.Client
	username         string
	interval         time.Duration
	stop             chan struct{}
	currentGameID    int
	gameStartTime    int64
	hideButtons      bool
	hideAchievements bool
	hub              StatusUpdater // optional; nil when web UI is disabled
	webUIEnabled     bool
	webPort          int       // current web UI port; 0 = not yet known
	webPortChangeFn  func(int) // called when web_port changes at runtime
	webShutdown      func()    // called when web_ui is disabled at runtime
	webStartFn       func(int) // called when web_ui is enabled at runtime; arg is port
}

// New creates a Worker with the given credentials and poll interval.
func New(username, apiKey string, intervalSecs int, hideButtons, hideAchievements bool) *Worker {
	if intervalSecs <= 0 {
		intervalSecs = 30
	}
	return &Worker{
		ra:               ra.New(username, apiKey),
		username:         username,
		interval:         time.Duration(intervalSecs) * time.Second,
		stop:             make(chan struct{}),
		hideButtons:      hideButtons,
		hideAchievements: hideAchievements,
	}
}

// Run starts the poll loop. It blocks until Stop is called.
func (w *Worker) Run() {
	log.Printf("RADPresence started (poll every %s)", w.interval)
	var rpc *discord.Client
	defer func() {
		if rpc != nil {
			_ = rpc.SetActivity(nil)
			rpc.Close()
		}
		log.Println("RADPresence stopped")
	}()

	for {
		if err := w.tick(&rpc); err != nil {
			log.Printf("[error] %v", err)
		}
		select {
		case <-w.stop:
			return
		case <-time.After(w.interval):
		}
	}
}

// SetHub attaches an optional StatusUpdater that receives a snapshot after every poll cycle.
func (w *Worker) SetHub(h StatusUpdater) {
	w.hub = h
	w.webUIEnabled = true
}

// SetWebShutdown stores a callback that is invoked when web_ui is set to false
// at runtime via config reload. It should call hub.Shutdown() to stop the server.
func (w *Worker) SetWebShutdown(fn func()) {
	w.webShutdown = fn
}

// SetWebStart stores a factory that starts a fresh web server on the given port.
// It is invoked when web_ui transitions false → true via config reload.
func (w *Worker) SetWebStart(fn func(port int)) {
	w.webStartFn = fn
}

// SetWebPortChange stores a callback invoked when web_port changes at runtime.
// It should forward the new port to hub.PortChange.
func (w *Worker) SetWebPortChange(fn func(int)) {
	w.webPortChangeFn = fn
}

// SetWebPort records the current listening port so tick() can detect changes.
func (w *Worker) SetWebPort(port int) {
	w.webPort = port
}

// Stop signals the Run loop to exit. Safe to call more than once.
func (w *Worker) Stop() {
	select {
	case <-w.stop:
	default:
		close(w.stop)
	}
}

func (w *Worker) tick(rpc **discord.Client) error {
	// Reload config each tick so changes saved via the web UI (or CLI) take
	// effect without restarting the service.
	if cfg, err := config.Load(); err == nil {
		if cfg.Username != "" {
			w.username = cfg.Username
			w.ra.UpdateCredentials(cfg.Username, cfg.APIKey)
		}
		if cfg.Interval > 0 {
			w.interval = time.Duration(cfg.Interval) * time.Second
		}
		w.hideButtons = cfg.HideButtons
		w.hideAchievements = cfg.HideAchievements
		// Detect web UI being disabled at runtime.
		if w.webUIEnabled && !cfg.WebUI && w.webShutdown != nil {
			log.Println("[web] web UI disabled via config — shutting down")
			w.webShutdown()
			w.hub = nil
			w.webUIEnabled = false
			w.webShutdown = nil
		}
		// Detect web UI being enabled at runtime.
		if !w.webUIEnabled && cfg.WebUI && w.webStartFn != nil {
			port := cfg.WebPort
			if port == 0 {
				port = 7842
			}
			log.Printf("[web] web UI enabled via config — starting on port %d", port)
			w.webStartFn(port)
		}
		// Detect web port change at runtime.
		if w.webUIEnabled && w.webPortChangeFn != nil {
			newPort := cfg.WebPort
			if newPort == 0 {
				newPort = 7842
			}
			if w.webPort != 0 && newPort != w.webPort {
				log.Printf("[web] port changed via config — switching to %d", newPort)
				w.webPortChangeFn(newPort)
				w.webPort = newPort
			}
		}
	}

	var snap StatusSnapshot
	snap.Username = w.username // used as fallback label; cleared below if API fails
	defer func() {
		if w.hub != nil {
			w.hub.UpdateStatus(snap)
		}
	}()

	summary, err := w.ra.GetUserSummary()
	if err != nil {
		snap.Username = "" // don't show stale username on API failure
		snap.LastError = err.Error()
		return fmt.Errorf("GetUserSummary: %w", err)
	}
	snap.RAConnected = true
	snap.Username = summary.User // confirmed from API
	snap.AvatarURL = summary.AvatarURL()
	snap.TotalPoints = summary.TotalPoints
	snap.TotalSoftcorePoints = summary.TotalSoftcorePoints
	snap.TotalTruePoints = summary.TotalTruePoints
	snap.Rank = summary.Rank
	snap.TotalRanked = summary.TotalRanked

	if !summary.IsActive() {
		snap.DiscordConnected = *rpc != nil
		// Only clear presence on the transition from playing → not playing.
		// w.currentGameID == 0 means we already cleared it a previous tick.
		if *rpc != nil && w.currentGameID != 0 {
			if err := (*rpc).SetActivity(nil); err != nil {
				log.Printf("[warn] clearing presence failed (%v) — reconnecting next cycle", err)
				(*rpc).Close()
				*rpc = nil
				snap.DiscordConnected = false
			} else {
				log.Println("Session ended — presence cleared")
			}
			w.currentGameID = 0
		}
		return nil
	}

	// Track when the game session started so the elapsed timer resets on game change.
	gameChanged := summary.LastGameID != w.currentGameID
	if gameChanged {
		w.currentGameID = summary.LastGameID
		w.gameStartTime = time.Now().Unix()
	}

	game, err := w.ra.GetGame(summary.LastGameID)
	if err != nil {
		return fmt.Errorf("GetGame(%d): %w", summary.LastGameID, err)
	}

	progress, err := w.ra.GetUserProgress(summary.LastGameID)
	if err != nil {
		return fmt.Errorf("GetUserProgress(%d): %w", summary.LastGameID, err)
	}

	// Populate snapshot for the web UI.
	snap.IsPlaying = true
	snap.GameTitle = game.Title
	snap.ConsoleName = game.ConsoleName
	snap.CoverArtURL = game.ArtURL()
	snap.RichPresenceMsg = summary.RichPresenceMsg
	snap.GameStartTime = w.gameStartTime
	snap.GameID = summary.LastGameID
	if !w.hideAchievements && progress.NumPossibleAchievements > 0 {
		mode := "Softcore"
		if progress.NumAchievedHardcore > 0 && progress.NumAchievedHardcore == progress.NumAchieved {
			mode = "Hardcore"
		}
		snap.Achievements = fmt.Sprintf("%d/%d achievements (%s)", progress.NumAchieved, progress.NumPossibleAchievements, mode)
	}

	// Reconnect to Discord if needed.
	if *rpc == nil {
		if !discord.IsRunning() {
			return nil // Discord not running — skip silently, no error logged
		}
		*rpc, err = discord.Connect(discord.AppID)
		if err != nil {
			return fmt.Errorf("discord connect: %w", err)
		}
		log.Println("Discord IPC connected")
	}

	activity := buildActivity(w.username, summary, game, progress, w.gameStartTime, w.hideButtons, w.hideAchievements)
	if err := (*rpc).SetActivity(activity); err != nil {
		log.Printf("[warn] SetActivity failed (%v) — reconnecting next cycle", err)
		(*rpc).Close()
		*rpc = nil
		return nil
	}
	snap.DiscordConnected = true

	// Only log when something meaningful changed.
	if gameChanged {
		log.Printf("Now playing: %s (%s)", game.Title, game.ConsoleName)
	}
	return nil
}

func buildActivity(username string, s ra.UserSummary, g ra.Game, p ra.UserProgress, startTime int64, hideButtons, hideAchievements bool) *discord.Activity {
	// name overrides the Discord app name ("RAPresence") shown in the presence card.
	name := g.Title

	details := s.RichPresenceMsg
	if len(details) > 128 {
		details = details[:128]
	}

	var state string
	largeText := g.Title
	if !hideAchievements && p.NumPossibleAchievements > 0 {
		mode := "Softcore"
		if p.NumAchievedHardcore > 0 && p.NumAchievedHardcore == p.NumAchieved {
			mode = "Hardcore"
		}
		state = fmt.Sprintf("%d/%d achievements (%s)", p.NumAchieved, p.NumPossibleAchievements, mode)
		largeText = fmt.Sprintf("%d/%d achievements", p.NumAchieved, p.NumPossibleAchievements)
	}

	act := &discord.Activity{
		Name:    name,
		Type:    0, // PLAYING
		Details: details,
		State:   state,
		Assets: &discord.ActivityAssets{
			LargeImage: g.ArtURL(),
			LargeText:  largeText,
			SmallText:  g.ConsoleName,
		},
		Timestamps: &discord.ActivityTimestamps{
			Start: startTime,
		},
	}

	if !hideButtons {
		act.Buttons = []discord.Button{
			{Label: "RA Profile", URL: fmt.Sprintf("https://retroachievements.org/user/%s", username)},
			{Label: "Game Page", URL: fmt.Sprintf("https://retroachievements.org/game/%d", s.LastGameID)},
		}
	}
	return act
}
