//go:build !windows

package discord

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
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

// IsRunning reports whether a Discord IPC socket is present on disk.
// It does not open a connection — it is safe to call frequently.
func IsRunning() bool {
	for _, dir := range candidateDirs() {
		for i := 0; i < 10; i++ {
			path := filepath.Join(dir, fmt.Sprintf("discord-ipc-%d", i))
			if _, err := os.Stat(path); err == nil {
				return true
			}
		}
	}
	return false
}

// candidateDirs returns the directories to search for Discord IPC sockets,
// in priority order: Flatpak Discord, Snap Discord, native Discord.
func candidateDirs() []string {
	var dirs []string
	if xdg := os.Getenv("XDG_RUNTIME_DIR"); xdg != "" {
		dirs = append(dirs,
			filepath.Join(xdg, "app", "com.discordapp.Discord"),
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
