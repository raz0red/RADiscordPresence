// Package presence implements the RetroAchievements → Discord Rich Presence poll loop.
package presence

import (
	"fmt"
	"log"
	"time"

	"github.com/raz0red/radpresence/internal/discord"
	"github.com/raz0red/radpresence/internal/ra"
)

// Worker polls RetroAchievements and keeps Discord Rich Presence up to date.
type Worker struct {
	ra            *ra.Client
	username      string
	interval      time.Duration
	stop          chan struct{}
	currentGameID int
	gameStartTime int64
}

// New creates a Worker with the given credentials and poll interval.
func New(username, apiKey string, intervalSecs int) *Worker {
	if intervalSecs <= 0 {
		intervalSecs = 10
	}
	return &Worker{
		ra:       ra.New(username, apiKey),
		username: username,
		interval: time.Duration(intervalSecs) * time.Second,
		stop:     make(chan struct{}),
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

// Stop signals the Run loop to exit. Safe to call more than once.
func (w *Worker) Stop() {
	select {
	case <-w.stop:
	default:
		close(w.stop)
	}
}

func (w *Worker) tick(rpc **discord.Client) error {
	summary, err := w.ra.GetUserSummary()
	if err != nil {
		return fmt.Errorf("GetUserSummary: %w", err)
	}

	if !summary.IsActive() {
		if *rpc != nil {
			if err := (*rpc).SetActivity(nil); err != nil {
				log.Printf("[warn] clearing presence failed (%v) — reconnecting next cycle", err)
				(*rpc).Close()
				*rpc = nil
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

	// Reconnect to Discord if needed.
	if *rpc == nil {
		*rpc, err = discord.Connect(discord.AppID)
		if err != nil {
			return fmt.Errorf("discord connect: %w", err)
		}
		log.Println("Discord IPC connected")
	}

	activity := buildActivity(w.username, summary, game, progress, w.gameStartTime)
	if err := (*rpc).SetActivity(activity); err != nil {
		log.Printf("[warn] SetActivity failed (%v) — reconnecting next cycle", err)
		(*rpc).Close()
		*rpc = nil
		return nil
	}

	// Only log when something meaningful changed.
	if gameChanged {
		log.Printf("Now playing: %s (%s)", game.Title, game.ConsoleName)
	}
	return nil
}

func buildActivity(username string, s ra.UserSummary, g ra.Game, p ra.UserProgress, startTime int64) *discord.Activity {
	// name overrides the Discord app name ("RAPresence") shown in the presence card.
	name := g.Title

	details := s.RichPresenceMsg
	if len(details) > 128 {
		details = details[:128]
	}

	var state string
	largeText := g.Title
	if p.NumPossibleAchievements > 0 {
		mode := "Softcore"
		if p.NumAchievedHardcore > 0 && p.NumAchievedHardcore == p.NumAchieved {
			mode = "Hardcore"
		}
		state = fmt.Sprintf("%d/%d achievements (%s)", p.NumAchieved, p.NumPossibleAchievements, mode)
		largeText = fmt.Sprintf("%d/%d achievements", p.NumAchieved, p.NumPossibleAchievements)
	}

	return &discord.Activity{
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
		Buttons: []discord.Button{
			{Label: "RA Profile", URL: fmt.Sprintf("https://retroachievements.org/user/%s", username)},
			{Label: "Game Page", URL: fmt.Sprintf("https://retroachievements.org/game/%d", s.LastGameID)},
		},
	}
}
