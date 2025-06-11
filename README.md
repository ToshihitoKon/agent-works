# Agent Works

A collection of AI-generated tools and utilities written in Go.

## Available Tools

### go-cmdeck

A Rundeck-style CLI/TUI job execution management tool that allows you to define, execute, and track jobs with comprehensive execution history.

**Key Features:**
- **Job Execution Management**: Define jobs with labels, descriptions, commands, and variables
- **Execution History**: Track job runs with timestamps, exit codes, success/failure status, and output
- **CLI Interface**: Command-line interface for job operations (list, run, add, remove)
- **TUI Interface**: Interactive terminal user interface with job status visualization
- **Variable Substitution**: Execute jobs with environment variable expansion
- **Detailed Logging**: Comprehensive execution reporting with command output capture

**Quick Start:**
```bash
cd go-cmdeck
go build -o go-cmdeck
./go-cmdeck init          # Initialize with example jobs
./go-cmdeck list          # List all jobs with execution status
./go-cmdeck run monitoring # Execute a job with history recording
./go-cmdeck tui           # Start interactive TUI mode
```

See the [go-cmdeck directory](./go-cmdeck/) for detailed documentation and usage examples.

### go-tap-ton
A minimalist terminal-based tap tempo analyzer written in Go.
- Real-time BPM calculation from spacebar taps
- Ao8 averaging algorithm with outlier detection  
- 60fps frame conversion for game development
- [Documentation](./go-tap-ton/README.md)

## Project Structure

Each tool is organized in its own directory with:
- Individual documentation (CLAUDE.md and README files)
- Self-contained Go modules
- Independent configuration and build processes

## Contributing

Each tool follows its own development guidelines. Please refer to the individual CLAUDE.md files in each tool directory for specific contribution guidelines and architectural decisions.

## License

This project contains various tools, each potentially under different licenses. Please check individual tool directories for license information.