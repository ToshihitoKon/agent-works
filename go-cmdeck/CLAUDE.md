# CLAUDE.md - Go CmDeck

## Project Overview

Go CmDeck is a Rundeck-style CLI/TUI job execution management tool written in Go. It allows users to define and execute jobs with execution history tracking, variable substitution, and comprehensive status reporting.

## Key Features

- **Job Execution Management**: Define jobs with labels, descriptions, commands, and variables
- **Execution History**: Track job runs with timestamps, exit codes, success/failure status, and output
- **CLI Interface**: Command-line interface for job operations (list, run, add, remove)
- **TUI Interface**: Interactive terminal user interface with job status visualization
- **Command Execution**: Execute jobs with variable substitution and detailed logging
- **Configuration**: JSON-based configuration stored in `~/.config/go-cmdeck/config.json`

## Architecture

### Core Components

1. **config.go**: Configuration management and JSON serialization with ExecutionResult tracking
2. **executor.go**: Job execution logic with detailed output capture and history recording
3. **cli.go**: Command-line interface with run command for job execution
4. **tui.go**: Terminal user interface using Bubble Tea with job status icons and details panel
5. **main.go**: Application entry point

### Dependencies

- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/lipgloss`: Styling for TUI

## Commands

- `go-cmdeck init`: Initialize configuration with example jobs
- `go-cmdeck list`: List all jobs with execution status
- `go-cmdeck execute <name>`: Execute job and record execution history
- `go-cmdeck run <name>`: Execute job and record execution history
- `go-cmdeck add`: Add new job
- `go-cmdeck remove <name>`: Remove job
- `go-cmdeck tui`: Start TUI mode

## Configuration Structure

```json
{
  "contexts": {
    "job-name": {
      "name": "job-name",
      "label": "Display Label",
      "description": "Optional description",
      "commands": {
        "run": "command to execute"
      },
      "variables": {
        "VAR_NAME": "value"
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

## Development Guidelines

- Follow Go best practices and conventions
- Use structured error handling
- Maintain backwards compatibility for configuration format
- Add tests for core functionality when extending features
- Use the existing styling patterns for TUI components

## Development History & Important Context

### Project Evolution
- **Origin**: Started as "any-context-switcher" for context management
- **Transformation**: Evolved to "go-cmdeck" - a Rundeck-style job execution tool
- **Key Insight**: User recognized similarity to Rundeck and requested paradigm shift from state management to job execution tracking

### Architecture Decisions
- **Unified Commands**: Both `execute` and `run` commands perform job execution with history tracking
- **No Current Context**: Removed CurrentContext concept from config - tool focuses on execution rather than state management
- **Execution History**: ExecutionResult struct tracks timestamp, success/failure, exit codes, and full output
- **TUI Status Icons**: Uses ✓/✗ icons instead of checkboxes to show execution status

### Recent Refactoring (2025-06-12)
- **switch → execute**: Renamed switch command to execute for clarity (context switching remnant)
- **Removed current command**: Eliminated current subcommand as it was context switching legacy
- **Simplified config**: Removed CurrentContext field from Config struct
- **Documentation**: Updated all docs to reflect job execution paradigm

### Key Implementation Details
- **Variable Substitution**: Uses `${VAR}` syntax for environment variable expansion
- **Exit Code Handling**: Properly captures exit codes using syscall.WaitStatus
- **Output Separation**: Distinguishes STDOUT and STDERR in execution results
- **Error Resilience**: Handles command execution failures gracefully with detailed reporting

### Testing Notes
- Always run `go build -o go-cmdeck` to verify compilation
- Test both CLI and TUI interfaces after changes
- Verify variable substitution works correctly
- Check execution history persistence across restarts

### Future Considerations
- Consider adding job scheduling capabilities
- Implement job dependencies/pipelines
- Add more detailed execution metrics
- Consider adding job templates or inheritance