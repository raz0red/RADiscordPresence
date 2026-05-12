//go:build windows

package discord

import (
	"fmt"
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

// openConn connects to the Discord IPC named pipe on Windows.
// Discord exposes \\.\pipe\discord-ipc-0 through discord-ipc-9.
func openConn() (net.Conn, error) {
	timeout := 2 * time.Second
	for i := 0; i < 10; i++ {
		conn, err := winio.DialPipe(fmt.Sprintf(`\\.\pipe\discord-ipc-%d`, i), &timeout)
		if err == nil {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("named pipe discord-ipc-0 through discord-ipc-9 not found — is Discord running?")
}
