//go:build windows

package discord

import (
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"
	"unsafe"

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

// IsRunning reports whether Discord is running by scanning the process list.
// Covers Discord stable, PTB, and Canary.
func IsRunning() bool {
	snap, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return true // assume running on error to avoid false negatives
	}
	defer syscall.CloseHandle(snap)
	var pe syscall.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	if err = syscall.Process32First(snap, &pe); err != nil {
		return false
	}
	for {
		name := strings.ToLower(syscall.UTF16ToString(pe.ExeFile[:]))
		if name == "discord.exe" || name == "discordptb.exe" || name == "discordcanary.exe" {
			return true
		}
		if err = syscall.Process32Next(snap, &pe); err != nil {
			break
		}
	}
	return false
}
