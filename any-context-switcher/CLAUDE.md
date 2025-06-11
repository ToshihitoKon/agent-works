# CLAUDE.md - Any Context Switcher

## Project Overview

Any Context Switcher is a CLI/TUI tool written in Go that allows users to manage and switch between different contexts (environments, configurations, etc.) with associated commands and variables.

## Key Features

- **Context Management**: Define contexts with labels, descriptions, status, commands, and variables
- **CLI Interface**: Command-line interface for context operations (list, switch, add, remove)
- **TUI Interface**: Interactive terminal user interface for context switching
- **Command Execution**: Execute commands when switching contexts with variable substitution
- **Configuration**: JSON-based configuration stored in `~/.config/any-context-switcher/config.json`

## Architecture

### Core Components

1. **config.go**: Configuration management and JSON serialization
2. **executor.go**: Context switching logic and command execution
3. **cli.go**: Command-line interface implementation
4. **tui.go**: Terminal user interface using Bubble Tea
5. **main.go**: Application entry point

### Dependencies

- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/lipgloss`: Styling for TUI

## Commands

- `any-context-switcher list`: List all contexts
- `any-context-switcher current`: Show current context
- `any-context-switcher switch <name>`: Switch to context
- `any-context-switcher add`: Add new context
- `any-context-switcher remove <name>`: Remove context
- `any-context-switcher tui`: Start TUI mode

## Configuration Structure

```json
{
  "current_context": "context-name",
  "contexts": {
    "context-name": {
      "name": "context-name",
      "label": "Display Label",
      "description": "Optional description",
      "status": "active|inactive",
      "commands": {
        "activate": "command to run when switching"
      },
      "variables": {
        "VAR_NAME": "value"
      }
    }
  }
}
```

## Development Guidelines

- Follow Go best practices and conventions
- Use structured error handling
- Maintain backwards compatibility for configuration format
- Add tests for core functionality when extending features
- Use the existing styling patterns for TUI components