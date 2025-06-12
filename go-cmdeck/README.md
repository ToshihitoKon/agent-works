# Go CmDeck

A Rundeck-style CLI/TUI job execution management tool written in Go. Go CmDeck allows you to define, execute, and track jobs with comprehensive execution history, making it easy to manage repetitive tasks and commands.

## Features

- **Job Execution Management**: Define jobs with labels, descriptions, commands, and variables
- **Execution History**: Track job runs with timestamps, exit codes, success/failure status, and detailed output
- **CLI Interface**: Command-line interface for job operations (list, run, add, remove)
- **TUI Interface**: Interactive terminal user interface with job status visualization using ✓/✗ icons
- **Variable Substitution**: Execute jobs with environment variable expansion using `${VAR}` syntax
- **Detailed Logging**: Comprehensive execution reporting with STDOUT/STDERR separation
- **Configurable Themes**: Customizable color themes for the TUI interface
- **Responsive Layout**: Adaptive terminal layout with overflow handling

## Installation

```bash
# Clone the repository
git clone https://github.com/ToshihitoKon/agent-works.git
cd agent-works/go-cmdeck

# Build the binary
go build -o go-cmdeck

# Make it executable and optionally move to PATH
chmod +x go-cmdeck
# sudo mv go-cmdeck /usr/local/bin/  # Optional: install globally
```

## Quick Start

```bash
# Initialize with example jobs
./go-cmdeck init

# List all jobs with execution status
./go-cmdeck list

# Execute a specific job
./go-cmdeck run monitoring

# Start interactive TUI mode
./go-cmdeck tui

# Add a new job
./go-cmdeck add -name "backup" -label "Database Backup" -description "Daily backup job"
```

## Commands

| Command | Description |
|---------|-------------|
| `init` | Initialize configuration with example jobs |
| `list`, `ls` | List all jobs with execution status |
| `execute`, `exec <name>` | Execute job and record execution history |
| `run <name>` | Execute job and record execution history |
| `add` | Add new job (interactive) |
| `remove`, `rm <name>` | Remove job |
| `tui` | Start TUI mode |
| `help` | Show help |

## TUI Interface

The TUI (Terminal User Interface) provides an interactive way to manage and execute jobs:

- **Job List**: Shows all jobs with status icons (✓ for success, ✗ for failure)
- **Navigation**: Use arrow keys or j/k to navigate
- **Job Execution**: Press space to execute the selected job
- **Job Details**: Bottom panel shows detailed information about the selected job
- **Real-time Updates**: Execution status updates in real-time

### TUI Controls

- `↑/↓` or `j/k`: Navigate through jobs
- `Space`: Execute selected job
- `q` or `Ctrl+C`: Quit

## Configuration

Configuration is stored in `~/.config/go-cmdeck/config.json`:

```json
{
  "contexts": {
    "monitoring": {
      "name": "monitoring",
      "label": "System Monitoring",
      "description": "Enable system monitoring tools",
      "commands": {
        "run": "echo 'Monitoring enabled' && ps aux | head -5"
      },
      "variables": {
        "LOG_PATH": "/var/log/monitoring",
        "INTERVAL": "5"
      },
      "last_result": {
        "timestamp": "2025-06-11T22:56:44.500268+09:00",
        "success": true,
        "exit_code": 0,
        "output": "Command execution output..."
      }
    }
  },
  "theme": {
    "title": "205",
    "selected": "199", 
    "border": "168",
    "output_title": "212"
  }
}
```

### Job Structure

Each job consists of:

- **name**: Unique identifier for the job
- **label**: Human-readable display name
- **description**: Optional description of what the job does
- **commands.run**: The command to execute
- **variables**: Key-value pairs for variable substitution
- **last_result**: Execution history (automatically managed)

### Variable Substitution

Use `${VARIABLE_NAME}` in commands to substitute variables:

```json
{
  "commands": {
    "run": "echo 'Connecting to ${HOST}:${PORT}'"
  },
  "variables": {
    "HOST": "localhost",
    "PORT": "8080"
  }
}
```

## Examples

### Creating a Backup Job

```bash
./go-cmdeck add -name "backup" -label "Database Backup" -description "Backup PostgreSQL database"
```

Then edit the configuration to add the command:

```json
{
  "name": "backup",
  "label": "Database Backup", 
  "description": "Backup PostgreSQL database",
  "commands": {
    "run": "pg_dump -h ${DB_HOST} -U ${DB_USER} ${DB_NAME} > backup_$(date +%Y%m%d).sql"
  },
  "variables": {
    "DB_HOST": "localhost",
    "DB_USER": "postgres", 
    "DB_NAME": "myapp"
  }
}
```

### Monitoring Job

```json
{
  "name": "monitoring",
  "label": "System Health Check",
  "description": "Check system resources and services",
  "commands": {
    "run": "echo 'System Status:' && uptime && echo 'Disk Usage:' && df -h / && echo 'Memory:' && free -h"
  },
  "variables": {
    "ALERT_EMAIL": "admin@company.com"
  }
}
```

## Development

### Prerequisites

- Go 1.21 or later
- Terminal with color support

### Dependencies

- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/lipgloss`: Styling for TUI

### Building from Source

```bash
git clone https://github.com/ToshihitoKon/agent-works.git
cd agent-works/go-cmdeck
go mod tidy
go build -o go-cmdeck
```

### Testing

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

Please refer to [CLAUDE.md](./CLAUDE.md) for development guidelines and architectural decisions.

## License

This project is part of the Agent Works collection. Please check the parent repository for license information.

## Similar Projects

- [Rundeck](https://www.rundeck.com/): Enterprise job scheduler and runbook automation
- [Ansible](https://www.ansible.com/): IT automation platform
- [Jenkins](https://www.jenkins.io/): CI/CD automation server

Go CmDeck is designed as a lightweight, terminal-based alternative for personal and small team use cases.