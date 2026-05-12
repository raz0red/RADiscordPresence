// Package svc wraps github.com/kardianos/service to register and manage
// RADPresence as a background service on Windows, macOS, and Linux.
package svc

import (
	"fmt"
	"log"

	ksvc "github.com/kardianos/service"

	"github.com/raz0red/radpresence/internal/config"
	"github.com/raz0red/radpresence/internal/presence"
)

var svcConfig = &ksvc.Config{
	Name:        "RADPresence",
	DisplayName: "RAD Presence",
	Description: "Mirrors your RetroAchievements session to Discord Rich Presence.",
	Option: ksvc.KeyValue{
		// Install as a user service (runs under the logged-in user's account).
		// This means no admin rights required, and %APPDATA% / config paths
		// resolve correctly to the current user's profile.
		"UserService": true,
	},
}

type program struct {
	worker *presence.Worker
}

func (p *program) Start(_ ksvc.Service) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	if cfg.Username == "" || cfg.APIKey == "" {
		return fmt.Errorf("username and api_key are not configured — run: radpresence set --username X --apikey Y")
	}
	p.worker = presence.New(cfg.Username, cfg.APIKey, cfg.Interval)
	go p.worker.Run()
	return nil
}

func (p *program) Stop(_ ksvc.Service) error {
	if p.worker != nil {
		p.worker.Stop()
	}
	return nil
}

func newService() (ksvc.Service, error) {
	return ksvc.New(&program{}, svcConfig)
}

// Interactive returns true when the binary is running in a terminal (not as a service).
func Interactive() bool {
	return ksvc.Interactive()
}

// Run hands control to the service manager. Called when not interactive.
func Run() error {
	s, err := newService()
	if err != nil {
		return err
	}
	logger, err := s.Logger(nil)
	if err == nil {
		log.SetOutput(&svcLogWriter{logger})
	}
	return s.Run()
}

// Install registers the binary as a system service.
//
// The resolved config directory is embedded as a --config-dir argument in the
// service registration so the service process (which may run as LocalSystem
// without user environment variables) can locate the correct config file.
func Install() error {
	dir, err := config.Dir()
	if err != nil {
		return fmt.Errorf("resolving config dir: %w", err)
	}
	svcConfig.Arguments = []string{"--config-dir", dir}
	s, err := newService()
	if err != nil {
		return err
	}
	return s.Install()
}

// Uninstall removes the system service registration.
func Uninstall() error {
	s, err := newService()
	if err != nil {
		return err
	}
	return s.Uninstall()
}

// Start starts the installed service.
func Start() error {
	s, err := newService()
	if err != nil {
		return err
	}
	return s.Start()
}

// Stop stops the running service.
func Stop() error {
	s, err := newService()
	if err != nil {
		return err
	}
	return s.Stop()
}

// Status returns a human-readable service status string.
func Status() (string, error) {
	s, err := newService()
	if err != nil {
		return "", err
	}
	st, err := s.Status()
	if err != nil {
		return "", err
	}
	switch st {
	case ksvc.StatusRunning:
		return "running", nil
	case ksvc.StatusStopped:
		return "stopped", nil
	default:
		return "unknown", nil
	}
}

// svcLogWriter bridges kardianos/service logger to the standard log package.
type svcLogWriter struct{ l ksvc.Logger }

func (w *svcLogWriter) Write(p []byte) (int, error) {
	_ = w.l.Info(string(p))
	return len(p), nil
}
