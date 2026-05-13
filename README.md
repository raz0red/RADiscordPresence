# RAD Presence

**RAD Presence** = RetroAchievements Discord Rich Presence. A background service that mirrors your [RetroAchievements](https://retroachievements.org/) session to [Discord Rich Presence](https://discord.com/developers/docs/rich-presence/overview).

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

### 3. Test in the foreground (Discord must be running)

```
radpresence run
```

You should see log output when you switch games. Press Ctrl+C to stop.

### 4. Install as a background service (optional)

```
# Windows — run as Administrator
radpresence install
radpresence start

# macOS — run with sudo
sudo radpresence install
sudo radpresence start

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

**First-time setup — build the Docker builder image (once):**

```
task build:image
```

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
