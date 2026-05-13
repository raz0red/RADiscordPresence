# RAD Presence

**RAD Presence** = RetroAchievements Discord Rich Presence.

A background service that mirrors your [RetroAchievements](https://retroachievements.org/) session to [Discord Rich Presence](https://discord.com/developers/docs/rich-presence/overview).

Inspired by [CheevoPresence](https://github.com/denzi-gh/CheevoPresence) — reimagined in Go as a cross-platform background service with no UI dependencies and a single self-contained binary.

![Discord Rich Presence screenshot](.github/assets/screenshot.png)

---

## Features

- Polls your RetroAchievements session every 10 seconds (configurable)
- Updates Discord with the game title, cover art, console, achievement progress, elapsed timer, and links to your RA profile and game page
- Clears presence automatically when you stop playing
- Runs as a native background service (Windows SCM / macOS launchd / Linux systemd) or in the foreground for testing
- Single binary — no runtime, no installer, no dependencies

---

## Getting Started

### 1. Get your API key

Log in to [retroachievements.org](https://retroachievements.org/), go to **Settings → Web API Key**, and copy it.

### 2. Save your credentials

```
radpresence set --username YOUR_RA_USERNAME --apikey YOUR_API_KEY
```

### 3. Make sure Discord is running

The Discord desktop app must be running on the same machine and logged in to the account you want the Rich Presence posted to. RAD Presence communicates with Discord over a local IPC socket — it does not work with the browser version of Discord.

### 4. Test in the foreground

```
radpresence run
```

You should see log output when you switch games. Press Ctrl+C to stop.

### 5. Install as a background service (optional)

Place the binary somewhere permanent **before** running `install` — the service is registered to its location at install time. If you move or delete the binary afterwards the service will fail to start. Suggested locations:

| Platform | Suggested location |
|---|---|
| Windows | `C:\Program Files\RADPresence\radpresence.exe` |
| macOS | `/usr/local/bin/radpresence` |
| Linux | `~/.local/bin/radpresence` or `/usr/local/bin/radpresence` |

To move it later, run `radpresence uninstall`, move the binary, then `radpresence install` again.

```
# Windows — run as Administrator
radpresence install
radpresence start

# macOS — no sudo needed (installs a LaunchAgent in ~/Library/LaunchAgents)
radpresence install
radpresence start

# Linux — no sudo needed (installs a systemd user service)
radpresence install
radpresence start
```

---

## All Commands

| Command | Description |
|---|---|
| `radpresence set --username X --apikey Y` | Save credentials to config |
| `radpresence set --interval 30` | Change the poll interval (seconds) |
| `radpresence set` | Show current config |
| `radpresence run` | Run in the foreground, Ctrl+C to stop |
| `radpresence run --username X --apikey Y` | Run with inline credentials (no saved config needed) |
| `radpresence install` | Register as a system service |
| `radpresence uninstall` | Remove the system service |
| `radpresence start` | Start the installed service |
| `radpresence stop` | Stop the running service |
| `radpresence status` | Show service status |
| `radpresence version` | Print version information |

---

## Building from Source

Requires [Docker](https://www.docker.com/) and [Task](https://taskfile.dev).

### Build Tasks

| Task | Description |
|---|---|
| `task build` | All platforms (Windows, Linux, macOS amd64 + arm64) |
| `task build:windows` | Windows amd64 only |
| `task build:linux` | Linux amd64 only |
| `task build:mac` | macOS amd64 + arm64 only |

Binaries are written to `dist/`.

### Dev Tasks

| Task | Description |
|---|---|
| `task fmt` | Auto-format all Go source files |
| `task fix` | Auto-format and apply golangci-lint auto-fixes |
| `task vet` | Run `go vet` |
| `task lint` | Run `golangci-lint` |
| `task validate` | Format + vet + lint (run before pushing) |
| `task clean` | Remove `dist/` |

---

## Config File Location

| Platform | Path |
|---|---|
| Windows | `%APPDATA%\RADPresence\config.json` |
| macOS | `~/Library/Application Support/RADPresence/config.json` |
| Linux | `~/.config/RADPresence/config.json` |

> **Note:** The API key is currently stored in the config file in plain text. Keyring integration (Windows Credential Manager, macOS Keychain, libsecret) is planned.

---

## Credits

Inspired by [CheevoPresence](https://github.com/denzi-gh/CheevoPresence) by [denzi_gh](https://github.com/denzi-gh).
