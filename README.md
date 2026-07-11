# Vibe-Shield

**The ultimate zero-config CLI error proxy for Vibe-Coding. Turn raw terminal crashes into surgical AI prompts in 1-click.**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go&logoColor=white)](go.mod)

---

<!-- ANIMATED DEMO PLACEHOLDER -->

> **Demo GIF (coming soon):** A developer runs `vibe-shield npm run dev`. The app crashes with a stack trace. Vibe-Shield detects the error, extracts the broken file with a `>>` marker on the failing line, and copies a clean markdown prompt to the clipboard. The developer pastes it into Claude (or any AI assistant) and gets a targeted fix in one click — no manual log scrubbing required.

---

## Why Vibe-Shield?

- **For beginners:** No more panic over walls of terminal output. When something breaks, you get a ready-to-paste prompt instead of staring at `node_modules` noise.
- **For pros:** Save time and tokens. **Context Slimming** strips system paths and dependency frames so your AI sees only what matters.
- **For vibe-coders:** Stay in flow. Run your command, crash, paste, fix — no config files, no API keys, no signup.

---

## Features

### Zero-Config & Local-First

Works right out of the box — no API keys, internet connection, or registration required. Vibe-Shield runs entirely on your machine and writes the prompt directly to your system clipboard.

### Context Slimming

Automatically filters system noise from stack traces: `node_modules`, `node:internal`, `internal/modules`, `site-packages`, `webpack-internal`, and other framework internals are stripped before the prompt is built.

### Smart Snippet Extraction

Finds the broken file from the stack trace and extracts **±15 lines** of code around the error, marking the failing line with `>>`:

```
42 | const config = loadConfig()
43 | const port = config.port
44 >> app.listen(port, () => {   // <-- crash here
45 |   console.log(`Listening on ${port}`)
```

---

## Installation

### One-line install (macOS / Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/neracu/vibe-shield/main/install.sh | bash
```

### Go install

```bash
go install github.com/neracu/vibe-shield/cmd/vibe-shield@latest
```

### Download a binary

Grab a prebuilt release for your platform from [GitHub Releases](https://github.com/neracu/vibe-shield/releases):

| Platform | Binary |
|----------|--------|
| macOS (Intel) | `vibe-shield-darwin-amd64` |
| macOS (Apple Silicon) | `vibe-shield-darwin-arm64` |
| Linux (amd64) | `vibe-shield-linux-amd64` |
| Windows (amd64) | `vibe-shield-windows-amd64.exe` |

Place the binary on your `PATH` and run `vibe-shield` from any terminal.

### Build from source

```bash
git clone https://github.com/neracu/vibe-shield.git
cd vibe-shield
go build -o vibe-shield ./cmd/vibe-shield
```

---

## Usage

Prefix any command with `vibe-shield`. On crash, the surgical prompt is copied to your clipboard automatically.

```bash
vibe-shield <your-command> [args...]
```

### Node.js

```bash
vibe-shield npm run dev
vibe-shield node server.js
```

### Python

```bash
vibe-shield python main.py
vibe-shield python -m pytest
```

When a crash is detected, you'll see:

```
🛡️ [Vibe-Shield] Shielding your code session (running: npm run dev)...
🚨 [Vibe-Shield] Crash detected in app.js:44!
📋 [Vibe-Shield] Surgical prompt successfully copied to clipboard. Paste it into your AI!
```

Paste into Claude, ChatGPT, Cursor, or any AI assistant — the prompt includes the error, slimmed stack trace, code snippet, and a focused fix instruction.

---

## How It Works

```
Your Command  →  Process Proxy  →  RegEx Detection  →  Snippet Extraction  →  System Clipboard
```

1. **Process Proxy** — Vibe-Shield wraps your command via `os/exec`, streaming stdout and stderr through capture buffers while preserving stdin and exit codes.
2. **RegEx Detection** — On non-zero exit, compiled patterns parse Python tracebacks and Node.js error stacks, skipping system paths to locate your code.
3. **Snippet Extraction** — The analyzer opens the failing file and pulls ±15 lines around the error line, marking it with `>>`.
4. **System Clipboard** — A structured markdown prompt (error, stack trace, snippet, last logs, fix instruction) is written to the clipboard — ready to paste.

All processing is local. Nothing leaves your machine unless you paste the prompt yourself.

---

## Development

```bash
# Run tests
go test ./...

# Build locally
go build -o vibe-shield ./cmd/vibe-shield

# Cross-compile all platforms
make build-all

# Smoke test with example crash fixtures
./vibe-shield python examples/fake_crash.py
./vibe-shield node examples/crash.js
```

---

## License

MIT — see [LICENSE](LICENSE).
