# Google Drive File Manager CLI

## ⚠️ PROJECT FROZEN

**This project is currently frozen due to Google OAuth2 API limitations.**

Google's OAuth2 implementation requires a `client_secret` even for Desktop applications, making it impossible to create a distributable binary without embedding sensitive credentials. This affects many developers trying to create desktop applications that integrate with Google APIs.

**Error encountered**: `oauth2: "invalid_request" "client_secret is missing."`

See [CLAUDE.md](./CLAUDE.md) for detailed information about attempted solutions and current status.

---

## What This Project Was Intended To Be

A CLI tool for managing Google Drive files with OAuth2 authentication, focusing on:

- **File Listing**: Browse My Drive and specific folders
- **Custom Properties**: Manage file appProperties for tag-based organization
- **User-Friendly**: Distributable binary with browser-based OAuth flow

## Technical Implementation

- **Language**: Go
- **Architecture**: Modular design with separate auth, CLI, and API packages
- **Authentication**: OAuth2 with browser-based flow and localhost redirect
- **Configuration**: XDG Base Directory compliant (`~/.config/gdrive-props/`)
- **Security**: Token files with 600 permissions, no embedded secrets

## Current State

The project includes:
- Complete OAuth2 authentication framework
- File listing functionality
- Comprehensive test suite
- XDG-compliant configuration management
- Cross-platform browser opening

However, it cannot function due to Google's requirement for client secrets in all OAuth flows.

## Build

```bash
go build -o gdrive-props-bin
```

## Alternative Solutions

If you need Google Drive CLI functionality, consider:
- [rclone](https://rclone.org/) - Mature solution with Google Drive support
- [gdrive](https://github.com/prasmussen/gdrive) - Alternative CLI tool
- Providing your own OAuth credentials via environment variables

## License

This project serves as a reference implementation and learning resource for OAuth2 integration patterns, even though it cannot be used as intended due to API limitations.