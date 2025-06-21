# CLAUDE.md - Google Drive File Manager

## Project Overview

This is a CLI tool for managing Google Drive files with OAuth2 authentication. Currently focuses on file listing functionality with plans for custom properties management.

## Project Structure

```
gdrive-props/
├── main.go                 # Entry point
├── cmd/                    # CLI commands
│   ├── root.go            # Root command and CLI setup
│   └── list.go            # List files command
├── pkg/
│   └── auth/              # OAuth2 authentication
│       ├── auth.go        # Authentication logic
│       └── auth_test.go   # Authentication tests
└── gdrive-props-bin       # Compiled binary
```

## Key Features

- **OAuth2 Authentication**: Secure authentication with Google Drive API
- **File Listing**: List files in My Drive or specified folders
- **Embedded Credentials**: OAuth2 client credentials embedded in binary
- **XDG Compliance**: Configuration stored in `~/.config/gdrive-props/`

## Technical Details

### Authentication
- Uses OAuth2 with Google Drive API scope
- OAuth2 client credentials embedded in binary (no external files needed)
- Access tokens cached in `~/.config/gdrive-props/token.json`
- Follows XDG Base Directory Specification
- Token file permissions set to 600 for security

### Configuration Directory
- Primary: `$XDG_CONFIG_HOME/gdrive-props/`
- Fallback: `~/.config/gdrive-props/`
- Contains only `token.json` (access/refresh tokens)

## CLI Commands

```bash
# File operations
gdrive-props list                          # List My Drive files
gdrive-props list [folder-id]              # List files in specific folder
gdrive-props list --page-size 100         # Limit number of files shown
gdrive-props list --all                    # Include trashed files
```

## Current Implementation

### List Command Features
- Shows file ID, name, type, size, modification time, and custom properties
- Formats file sizes (B/KB/MB/GB)
- Identifies folders and Google Apps (Docs, Sheets, etc.)
- Displays custom properties (appProperties) if present
- Supports pagination and trash filtering
- Clean tabular output format

### Authentication Features
- Automatic OAuth2 flow initiation on first run
- Token refresh handling
- Comprehensive test coverage (8 test cases)
- Context cancellation support
- Proper error handling

## Testing

Run tests with:
```bash
go test -v ./pkg/auth/
```

Test coverage includes:
- OAuth2 configuration validation
- Token save/load operations  
- File permissions verification
- Configuration directory handling
- Context cancellation
- Error conditions

## Setup Instructions

1. Ensure you have valid OAuth2 credentials embedded in the binary
2. Run `./gdrive-props-bin list` to initiate OAuth flow
3. Grant permissions in browser
4. Enter authorization code when prompted
5. Token will be saved automatically for future use

## Development Guidelines

### Code Organization
- Keep CLI logic in `cmd/` package
- Authentication logic in `pkg/auth/` package
- Follow Go naming conventions
- Use contexts for API calls

### Error Handling
- Return descriptive error messages
- Use fmt.Errorf for error wrapping
- Handle Google API errors gracefully

### Testing
- Write tests for all new functionality
- Maintain high test coverage
- Use table-driven tests where appropriate
- Test error conditions and edge cases

## Security Considerations

- OAuth2 client credentials are embedded in binary
- Token files have 600 permissions (user read/write only)
- No sensitive data logged or exposed
- Proper token refresh handling

## Future Enhancements

- Custom properties (appProperties) management
- Tag-based file organization
- Search functionality
- File upload/download capabilities
- Batch operations