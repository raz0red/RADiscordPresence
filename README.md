# RADiscordPresence

A background service that mirrors your [RetroAchievements](https://retroachievements.org/) session to [Discord Rich Presence](https://discord.com/developers/docs/rich-presence/overview).

Inspired by [CheevoPresence](https://github.com/denzi-gh/CheevoPresence) — reimagined in Go as a cross-platform background service with no UI dependencies and a single self-contained binary.

---

## Features

- Polls your RetroAchievements session every 10 seconds (configurable)
- Updates Discord with the game title, console, achievement progress, and links to your RA profile and game page
- Clears presence automatically when you stop playing
- Runs as a native background service (Windows SCM / macOS launchd / Linux systemd) or in the foreground for testing
- Single binary — no runtime, no installer, no dependencies

---

## Getting Started

### 1. Get your API key

Log in to [retroachievements.org](https://retroachievements.org/), go to **Settings → Web API Key**, and copy it.

### 2. Save your credentials

```
radiscordpresence set --username YOUR_RA_USERNAME --apikey YOUR_API_KEY
```

### 3. Test in the foreground (Discord must be running)

```
radiscordpresence run
```

You should see log output every 10 seconds. Press Ctrl+C to stop.

### 4. Install as a background service (optional)

```
# Windows — run as Administrator
radiscordpresence install
radiscordpresence start

# macOS / Linux — run with sudo
sudo radiscordpresence install
sudo radiscordpresence start
```

---

## All Commands

| Command | Description |
|---|---|
| `set --username X --apikey Y` | Save credentials to config |
| `set --interval 30` | Change the poll interval (seconds) |
| `run` | Run in the foreground, Ctrl+C to stop |
| `run --username X --apikey Y` | Run with inline credentials (no saved config needed) |
| `install` | Register as a system service |
| `uninstall` | Remove the system service |
| `start` | Start the installed service |
| `stop` | Stop the running service |
| `status` | Show service status |

---

## Building

Requires [Docker](https://www.docker.com/).

**Windows (PowerShell) — Windows binary only (fastest for local testing):**

```powershell
.\scripts\build.ps1 -WindowsOnly
```

**All platforms (Windows, Linux, macOS amd64 + arm64):**

```powershell
.\scripts\build.ps1
```

On subsequent runs, skip the image rebuild:

```powershell
.\scripts\build.ps1 -SkipImageBuild
```

Binaries are written to `dist/`.

---

## Config File Location

| Platform | Path |
|---|---|
| Windows | `%APPDATA%\RADiscordPresence\config.json` |
| macOS | `~/Library/Application Support/RADiscordPresence/config.json` |
| Linux | `~/.config/RADiscordPresence/config.json` |

> **Note:** The API key is currently stored in the config file in plain text. Keyring integration (Windows Credential Manager, macOS Keychain, libsecret) is planned.

---

## Credits

Inspired by [CheevoPresence](https://github.com/denzi-gh/CheevoPresence) by [denzi_gh](https://github.com/denzi-gh).
