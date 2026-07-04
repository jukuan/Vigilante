# Go Vigilante

A lightweight, zero-dependency log monitoring and alerting daemon written in Go.  
Watch directories of log files, match lines against regex patterns, and trigger external bash scripts — with built‑in accumulation to avoid alert storms.

**Why "Vigilante"?**  

A vigilante is someone who takes matters into their own hands when official systems fall short. That's exactly what this tool does. When your primary monitoring stack is too heavy, too slow, or simply not an option — Vigilante steps in. It watches your logs, spots trouble, and alerts you. No agents, no dashboards, no infrastructure. Just a small binary and a config file, doing the job that bigger tools couldn't.

Think of it as your silent partner — always watching, never sleeping, and ready to act when the official channels can't.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Directory watching** – Monitor a folder for new or rotated log files (supports globs like `*.log`).
- **Regex / literal matching** – Case‑insensitive patterns; special regex chars are automatically escaped if the pattern isn't a valid regex.
- **Accumulate-then-flush** – Collect matching lines over a configurable cooldown window, then fire one alert with a summary.
- **External actions** – Alert scripts are plain Bash files – use email, Slack, PagerDuty, or anything you can script.
- **State persistence** – Offsets are saved in a JSON file, so a restart picks up where it left off.
- **Simple configuration** – Single YAML file for all rules.
- **Production‑ready** – Graceful shutdown flushes pending alerts and saves state. Built‑in scripts to start/stop/restart.
- **Tested** – Unit tests for core logic (matcher, config, state, alerter).

## How It Works
[log files] → [file watcher] → [line matcher] → [alert window] → [bash scripts]

- New lines in matching files are checked against your patterns.
- Matched lines accumulate per rule.
- After `cooldown_seconds`, the alert fires with:  
  `ALERT: 12 lines in logs for last 5 minutes with like FATAL: disk full`
- The alert message is passed as the first argument to your script(s).

## Getting Started

### Prerequisites

- Go 1.21 or later
- Linux (primary target; Windows works with some limitations)

### Installation

```bash
git clone https://github.com/jukuan/vigilante.git
cd vigilante
./bin/dev_setup.sh   # installs dependencies and makes scripts executable
./bin/build.sh       # compiles the binary
```

## Quick Test

Create a log file to watch:
```bash
mkdir -p /tmp/vigilante-demo
echo "INFO: all good" > /tmp/vigilante-demo/app.log
```

Edit config.yaml (or copy and adjust):
```yaml

    inactivity_seconds: 300
    state_file: state.json

    rules:
      - name: demo-errors
        log_dir: /tmp/vigilante-demo
        file_glob: "*.log"
        patterns:
          - "ERROR"
          - "FATAL"
        actions:
          - scripts/dummy-alert.sh
        cooldown_seconds: 30
```

###    Start the daemon:

```bash
    ./bin/start.sh
```

Append a matching line:
```bash

echo "ERROR: something broke" >> /tmp/vigilante-demo/app.log

Wait up to 30 seconds and check alerts.log:
bash

cat alerts.log

Stop the daemon:
bash

./bin/stop.sh
```

### Configuration

All rules live in config.yaml. See the included file for a full example.
```yaml

inactivity_seconds: 300          # global idle timeout for files
state_file: state.json           # where to store file offsets

rules:
  - name: critical-errors        # unique rule name
    log_dir: /var/log/myapp      # directory to watch
    file_glob: "*.log"           # file pattern (glob)
    patterns:                    # lines matching ANY of these trigger
      - "FATAL"
      - "CRITICAL"
      - "panic"
    actions:                     # scripts to run (passed the alert message)
      - scripts/email.sh
      - scripts/slack.sh
    cooldown_seconds: 300        # aggregation window
    inactivity_seconds: 600      # optional per‑rule override
```
Important: patterns are compiled as case‑insensitive regexes. If a pattern is not a valid regex, it is escaped and treated as a literal substring.

## Included Scripts

bin/build.sh – compile the project

bin/start.sh [config.yaml] – start as a daemon

bin/stop.sh – gracefully stop

bin/restart.sh – stop then start

bin/dev_setup.sh – install deps and set permissions

bin/test.sh – run unit tests

Example alert scripts in scripts/:

dummy-alert.sh – append alert to a local file

email.sh – send an email (customize it!)

slack.sh – post to a Slack webhook (customize it!)

### Testing

Run the unit tests:
```bash
./bin/test.sh
```

### Roadmap & Contributing

Ideas for the future:

- Windows support (already mostly compatible via fsnotify)
- Config hot‑reload on SIGHUP
- Alert script timeout handling
- Better support for copy‑truncate log rotation

Pull requests are welcome! Please open an issue first to discuss what you'd like to change.

## License

MIT © 2026 @Jukuan

Vigilante – watch your logs, stay notified, sleep peacefully.
