# agent-works

A Git repository for managing AI-generated tools. Each tool is organized in its own directory under the project root.

## Directory Structure

- **Project Root**: Contains individual tool directories, `.github/` setup files, and project-wide configuration files only
- **Tool Directories**: Each contains a specific tool and its configuration files (including individual CLAUDE.md files)

## Available Tools

### go-tap-ton
A minimalist terminal-based tap tempo analyzer written in Go.
- Real-time BPM calculation from spacebar taps
- Ao8 averaging algorithm with outlier detection  
- 60fps frame conversion for game development
- [Documentation](./go-tap-ton/README.md)

## Important Principles

- No universal rules apply across all tools in this project
- Each tool follows its own rules defined in its directory's CLAUDE.md file
- When adding new tools, create a dedicated directory directly under the project root

## Working Guidelines

- When working on a specific tool, prioritize the CLAUDE.md file within that tool's directory
- Project root level changes should be limited to `.github/` directory and project-wide configuration files