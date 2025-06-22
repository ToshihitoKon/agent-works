package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	// Token file to store access and refresh tokens
	TokenFile = "token.json"
	// Configuration directory
	ConfigDir = "gdrive-props"
)

// Auth handles OAuth2 authentication for Google Drive API
type Auth struct {
	config    *oauth2.Config
	tokenFile string
}

// NewAuth creates a new Auth instance using Google's default application credentials
// This follows Google's official quickstart pattern
func NewAuth() (*Auth, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("unable to create config directory: %v", err)
	}

	tokenPath := filepath.Join(configDir, TokenFile)

	// Find available port for localhost redirect
	port, err := findAvailablePort()
	if err != nil {
		return nil, fmt.Errorf("unable to find available port: %v", err)
	}

	// Google's OAuth2 config for browser-based flow
	config := &oauth2.Config{
		ClientID:    "334618002502-m87faubnbb1ekmehr6k85sl18dh2g90q.apps.googleusercontent.com",
		RedirectURL: fmt.Sprintf("http://localhost:%d/callback", port),
		Scopes:      []string{drive.DriveScope},
		Endpoint:    google.Endpoint,
	}

	return &Auth{
		config:    config,
		tokenFile: tokenPath,
	}, nil
}

// getConfigDir returns the configuration directory path
// Uses XDG Base Directory specification: ~/.config/gdrive-props
func getConfigDir() (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to get user home directory: %v", err)
		}
		configDir = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configDir, ConfigDir), nil
}

// findAvailablePort finds an available localhost port
func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func (a *Auth) GetClient(ctx context.Context) (*http.Client, error) {
	tok, err := a.tokenFromFile()
	if err != nil {
		tok, err = a.getTokenFromWeb()
		if err != nil {
			return nil, err
		}
		a.saveToken(tok)
	}
	return a.config.Client(ctx, tok), nil
}

// getTokenFromWeb requests a token from the web using browser-based flow
func (a *Auth) getTokenFromWeb() (*oauth2.Token, error) {
	// Extract port from redirect URL
	redirectURL, err := url.Parse(a.config.RedirectURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redirect URL: %v", err)
	}
	port := redirectURL.Port()

	// Channel to receive authorization code
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	// Create new mux to avoid global handler conflicts
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no authorization code received")
			return
		}

		// Send success page
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
		<html>
		<head><title>Authorization Successful</title></head>
		<body>
		<h1>✅ Authorization Successful!</h1>
		<p>You can now close this browser window and return to the CLI application.</p>
		</body>
		</html>`)

		codeCh <- code
	})

	// Start local HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- fmt.Errorf("failed to start server: %v", err)
		}
	}()

	// Generate authorization URL
	authURL := a.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	
	fmt.Printf("Opening browser for authentication...\n")
	fmt.Printf("If browser doesn't open automatically, visit: %s\n", authURL)

	// Try to open browser automatically
	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Failed to open browser automatically: %v\n", err)
		fmt.Printf("Please open the URL manually in your browser.\n")
	}

	// Wait for authorization code or error
	var authCode string
	select {
	case authCode = <-codeCh:
		// Success
	case err := <-errCh:
		server.Shutdown(context.Background())
		return nil, err
	}

	// Shutdown server
	server.Shutdown(context.Background())

	// Exchange authorization code for token
	tok, err := a.config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return tok, nil
}

// tokenFromFile retrieves a token from a local file.
func (a *Auth) tokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(a.tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken saves a token to a file path.
func (a *Auth) saveToken(token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", a.tokenFile)
	f, err := os.OpenFile(a.tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// GetDriveService creates and returns a Google Drive service client
func (a *Auth) GetDriveService(ctx context.Context) (*drive.Service, error) {
	client, err := a.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	return srv, nil
}
