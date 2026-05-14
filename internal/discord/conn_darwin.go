//go:build darwin

package discord

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
)

// openConn connects to the Discord IPC Unix socket on macOS.
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

// IsRunning reports whether Discord is running by checking the process list
// via /usr/bin/pgrep (always present on macOS regardless of PATH).
func IsRunning() bool {
	for _, name := range []string{"Discord", "DiscordPTB", "DiscordCanary"} {
		if exec.Command("/usr/bin/pgrep", "-x", name).Run() == nil {
			return true
		}
	}
	return false
}

func candidateDirs() []string {
	var dirs []string
	if tmp := os.Getenv("TMPDIR"); tmp != "" {
		dirs = append(dirs, tmp)
	}
	dirs = append(dirs, os.TempDir())
	return dirs
}
