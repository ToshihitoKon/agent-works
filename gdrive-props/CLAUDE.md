# CLAUDE.md - Google Drive File Manager

## ⚠️ PROJECT FROZEN

**This project is currently frozen due to Google OAuth2 API limitations.**

**Issue**: Google's OAuth2 implementation requires a `client_secret` even for Desktop applications, making it impossible to create a distributable binary without embedding sensitive credentials. This is a known limitation affecting many developers in the community.

**Attempted Solutions**:
- Public client configuration (failed - still requires client_secret)
- Browser-based OAuth flow with localhost redirect (failed - still requires client_secret)
- PKCE implementation (failed - still requires client_secret)
- Manual copy/paste authorization flow (failed - still requires client_secret)

**Error Message**: `oauth2: "invalid_request" "client_secret is missing."`

**Current Status**: Development suspended until Google provides a viable solution for distributable desktop applications or community finds a workaround.

**Alternative Approaches** (not implemented):
- Require users to provide their own OAuth credentials via environment variables
- Use Service Account authentication (requires different use case)
- Switch to a different cloud storage provider with better OAuth support

---

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
- Uses Google's official OAuth2 quickstart pattern
- Public client configuration (no client secret required)
- Simple copy/paste authorization code flow
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
- Google's official OAuth2 quickstart implementation
- Automatic OAuth2 flow initiation on first run
- Token refresh handling
- Comprehensive test coverage (6 test cases)
- Simple and reliable authentication flow
- Proper error handling

## Testing

Run tests with:
```bash
go test -v ./pkg/auth/
```

Test coverage includes:
- OAuth2 public client configuration validation
- Token save/load operations  
- File permissions verification
- Configuration directory handling
- XDG Base Directory compliance
- Error conditions

## Setup Instructions

1. Create a Google Cloud Project and enable the Google Drive API
2. Create OAuth2 credentials for a "Desktop application"
3. Replace the placeholder ClientID in `pkg/auth/auth.go` with your actual client ID
4. Build and run `./gdrive-props-bin list` to initiate OAuth flow
5. Copy the authorization URL to your browser and grant permissions
6. Copy the authorization code back to the CLI when prompted
7. Token will be saved automatically for future use

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

- Uses public client OAuth2 flow (no client secret required)
- Follows Google's official security recommendations
- Token files have 600 permissions (user read/write only)
- No sensitive data logged or exposed
- Automatic token refresh handling
- Simple and secure copy/paste authorization flow

## Future Enhancements

- Custom properties (appProperties) management
- Tag-based file organization
- Search functionality
- File upload/download capabilities
- Batch operations