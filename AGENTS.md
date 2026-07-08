# GoVigilante — Context for AI Agents

## Project Summary
**vigilante** is a lightweight, Linux-first log monitoring and alerting daemon written in Go. It watches directories of log files (with glob patterns), scans new lines for regex patterns, and triggers external bash scripts when matching lines are found. Alerts are accumulated over a configurable cooldown window and then reported with a summary message.

## Key Design Decisions
- **Directory watching with glob:** Monitors a directory for files matching a glob (e.g., `app-*.log`). New and rotated files are handled automatically.
- **No historical backfill:** On startup, only new content is processed; existing file content is skipped.
- **Incomplete lines discarded:** The last incomplete line of a read is ignored (simpler than buffering across reads).
- **Accumulate-then-flush alerting:** Within a cooldown window (per rule), matching lines are counted. At the end of the window, the alert script fires once with a message: `ALERT: (rule-name) (n) lines in logs for last (m) minutes with like (s)`, where `(s)` is the first 16 characters of the first matched line.
- **State persistence:** Offsets into log files are saved in a JSON file (`state.json`) periodically and on shutdown, enabling resume after restart.
- **No hot config reload:** Configuration is read at startup only; restart the daemon to pick up changes.
- **Files that disappear are forgotten:** If a file is removed, the watcher drops it immediately without tracking.
- **OR pattern matching:** Multiple patterns per rule are combined with OR logic (any match triggers).
- **Case-insensitive matching:** Patterns are compiled with `(?i)`.

## Architecture Overview
[Config (YAML)] → [File Watchers (fsnotify + polling)] → [Line Matcher (regex)] → [Alert Manager (window per rule)] → [Bash Scripts]

- **main.go** – Entry point. Loads config, initialises state & alert managers, spawns watchers, handles signals.
- **config.go** – Parses YAML config into structs. Validates rules.
- **watcher.go** – For each rule, watches a directory for file events and polls tracked files for new content. Splits content into lines, discards trailing incomplete line, matches, sends matched lines to a channel.
- **matcher.go** – Compiles regex patterns and tests lines. Provides `TruncateLogStr()` helper.
- **alerter.go** – Receives matched lines, accumulates them per rule in an alert window. On cooldown expiry, formats message and executes all action scripts with the message as first argument.
- **state.go** – Thread-safe JSON offset store. Loads on startup, saves periodically and on shutdown.

## Project File Structure

vigilante/
├── main.go
├── config.go
├── config.yaml # main configuration (tracked in git)
├── watcher.go
├── matcher.go
├── alerter.go
├── state.go
├── go.mod
├── go.sum
├── .gitignore
├── AGENTS.md # this file
├── bin/
│ ├── build.sh
│ ├── start.sh
│ ├── stop.sh
│ └── restart.sh
└── scripts/ # example alert scripts
├── email.sh
└── slack.sh


## Configuration Example (config.yaml)
```yaml
inactivity_seconds: 300
state_file: state.json

rules:
  - name: critical-errors
    log_dir: /var/log/myapp
    file_glob: "*.log"
    patterns:
      - "FATAL"
      - "CRITICAL"
      - "panic"
    actions:
      - scripts/email.sh
      - scripts/slack.sh
    cooldown_seconds: 300
```

- inactivity_seconds – Global default for considering a file idle (stops polling).

Rules can override inactivity_seconds.

- cooldown_seconds – Alert aggregation window.

## How to Use

Configure config.yaml with your log directories and patterns.

Make scripts executable: chmod +x bin/*.sh scripts/*.sh

Build: ./bin/build.sh

Start: ./bin/start.sh [config.yaml]

Stop: ./bin/stop.sh

Restart: ./bin/restart.sh

## Important Implementation Details

fsnotify is used for file events, but periodic polling (every 5s) is also done to catch new files that may have appeared while the daemon was starting or missed events.

On SIGTERM/SIGINT, pending alert windows are flushed, state is saved, and watchers gracefully stop.

Alert scripts receive the formatted message as a single command-line argument.

Offset tracking is done by full file path; on restart, if a file has grown since last offset, only new bytes are processed.

# Potential Enhancements (not yet implemented)

Windows support (fsnotify works, but inode tracking would need adaptation).

Config hot reloading via SIGHUP.

Alert script timeout control.

Better handling of copy-truncate rotation strategies (currently handles typical create/rename/append rotation).
