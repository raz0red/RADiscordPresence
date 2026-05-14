//go:build !windows

package discord

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

// openConn connects to the Discord IPC Unix socket on Linux and macOS.
// Search order: Flatpak Discord → Snap Discord → native Discord.
func openConn() (net.Conn, error) {
	for _, dir := range candidateDirs() {
		for i := 0; i < 10; i++ {
			path := filepath.Join(dir, fmt.Sprintf("discord-ipc-%d", i))
			conn, err := net.Dial("unix", path)
			if err == nil {
				return conn, nil
			}
		}
	}
	return nil, fmt.Errorf("discord-ipc socket not found in any candidate directory — is Discord running?")
}

// IsRunning reports whether Discord is running by attempting a socket connection.
// A stale socket (crash) fails with "connection refused"; a live one succeeds.
func IsRunning() bool {
	for _, dir := range candidateDirs() {
		for i := 0; i < 10; i++ {
			path := filepath.Join(dir, fmt.Sprintf("discord-ipc-%d", i))
			conn, err := net.DialTimeout("unix", path, 100*time.Millisecond)
			if err == nil {
				_ = conn.Close()
				return true
			}
		}
	}
	return false
}

// candidateDirs returns the directories to search for Discord IPC sockets,
// in priority order: Flatpak Discord, Flatpak Vesktop, Snap Discord, native Discord / Vesktop.
func candidateDirs() []string {
	var dirs []string
	if xdg := os.Getenv("XDG_RUNTIME_DIR"); xdg != "" {
		dirs = append(dirs,
			// Newer Flatpak sandbox path (bwrap xdg-run proxy)
			filepath.Join(xdg, ".flatpak", "dev.vencord.Vesktop", "xdg-run"),
			filepath.Join(xdg, ".flatpak", "com.discordapp.Discord", "xdg-run"),
			// Older Flatpak sandbox path
			filepath.Join(xdg, "app", "com.discordapp.Discord"),
			filepath.Join(xdg, "app", "dev.vencord.Vesktop"),
			filepath.Join(xdg, "snap.discord"),
			xdg,
		)
	}
	if tmp := os.Getenv("TMPDIR"); tmp != "" {
		dirs = append(dirs, tmp)
	}
	dirs = append(dirs, os.TempDir())
	return dirs
}
