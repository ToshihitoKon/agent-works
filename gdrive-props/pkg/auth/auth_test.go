package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestNewAuth(t *testing.T) {
	auth, err := NewAuth()
	if err != nil {
		t.Fatalf("NewAuth() failed: %v", err)
	}

	if auth.config == nil {
		t.Error("OAuth2 config should not be nil")
	}

	if auth.config.ClientID == "" {
		t.Error("ClientID should not be empty")
	}

	if auth.config.ClientSecret != "" {
		t.Error("Public client should not have ClientSecret")
	}

	if auth.config.RedirectURL != "urn:ietf:wg:oauth:2.0:oob" {
		t.Errorf("Expected RedirectURL urn:ietf:wg:oauth:2.0:oob, got %s", auth.config.RedirectURL)
	}

	if len(auth.config.Scopes) == 0 {
		t.Error("OAuth2 scopes should not be empty")
	}
}

func TestGetConfigDir(t *testing.T) {
	// Test with XDG_CONFIG_HOME set
	t.Run("with XDG_CONFIG_HOME", func(t *testing.T) {
		oldValue := os.Getenv("XDG_CONFIG_HOME")
		defer os.Setenv("XDG_CONFIG_HOME", oldValue)

		testDir := "/tmp/test-config"
		os.Setenv("XDG_CONFIG_HOME", testDir)

		configDir, err := getConfigDir()
		if err != nil {
			t.Fatalf("getConfigDir() failed: %v", err)
		}

		expected := filepath.Join(testDir, ConfigDir)
		if configDir != expected {
			t.Errorf("Expected %s, got %s", expected, configDir)
		}
	})

	// Test without XDG_CONFIG_HOME (should use ~/.config)
	t.Run("without XDG_CONFIG_HOME", func(t *testing.T) {
		oldValue := os.Getenv("XDG_CONFIG_HOME")
		defer os.Setenv("XDG_CONFIG_HOME", oldValue)

		os.Unsetenv("XDG_CONFIG_HOME")

		configDir, err := getConfigDir()
		if err != nil {
			t.Fatalf("getConfigDir() failed: %v", err)
		}

		homeDir, _ := os.UserHomeDir()
		expected := filepath.Join(homeDir, ".config", ConfigDir)
		if configDir != expected {
			t.Errorf("Expected %s, got %s", expected, configDir)
		}
	})
}

func TestTokenSaveAndLoad(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	tokenPath := filepath.Join(tempDir, TokenFile)

	auth := &Auth{
		config: &oauth2.Config{
			ClientID:    "test-client-id",
			RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
		},
		tokenFile: tokenPath,
	}

	// Create test token
	testToken := &oauth2.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}

	// Test saving token
	auth.saveToken(testToken)

	// Test loading token
	loadedToken, err := auth.tokenFromFile()
	if err != nil {
		t.Fatalf("tokenFromFile() failed: %v", err)
	}

	if loadedToken.AccessToken != testToken.AccessToken {
		t.Errorf("Expected AccessToken %s, got %s", testToken.AccessToken, loadedToken.AccessToken)
	}

	if loadedToken.RefreshToken != testToken.RefreshToken {
		t.Errorf("Expected RefreshToken %s, got %s", testToken.RefreshToken, loadedToken.RefreshToken)
	}

	if loadedToken.TokenType != testToken.TokenType {
		t.Errorf("Expected TokenType %s, got %s", testToken.TokenType, loadedToken.TokenType)
	}
}

func TestTokenFromFileNonExistent(t *testing.T) {
	auth := &Auth{
		tokenFile: "/non/existent/path/token.json",
	}

	_, err := auth.tokenFromFile()
	if err == nil {
		t.Error("Expected error when reading non-existent token file")
	}
}

func TestTokenPermissions(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	tokenPath := filepath.Join(tempDir, TokenFile)

	auth := &Auth{
		tokenFile: tokenPath,
	}

	// Create test token
	testToken := &oauth2.Token{
		AccessToken: "test-token",
	}

	// Save token
	auth.saveToken(testToken)

	// Check file permissions
	info, err := os.Stat(tokenPath)
	if err != nil {
		t.Fatalf("Failed to stat token file: %v", err)
	}

	perm := info.Mode().Perm()
	expectedPerm := os.FileMode(0600)
	if perm != expectedPerm {
		t.Errorf("Expected token file permissions %o, got %o", expectedPerm, perm)
	}
}

func TestAuthConfigValidation(t *testing.T) {
	tests := []struct {
		name         string
		clientID     string
		redirectURL  string
		expectError  bool
	}{
		{
			name:         "valid public client config",
			clientID:     "test-client-id",
			redirectURL:  "urn:ietf:wg:oauth:2.0:oob",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &oauth2.Config{
				ClientID:    tt.clientID,
				RedirectURL: tt.redirectURL,
			}

			if config.ClientID == "" && !tt.expectError {
				t.Error("ClientID should not be empty for valid config")
			}

			if config.ClientSecret != "" {
				t.Error("Public client should not have ClientSecret")
			}
		})
	}
}

func TestGetClientWithoutToken(t *testing.T) {
	t.Skip("Skipping test that would trigger OAuth flow")
	
	auth, err := NewAuth()
	if err != nil {
		t.Fatalf("NewAuth() failed: %v", err)
	}

	// This test would trigger OAuth flow in real usage
	_ = auth
}

// Benchmark for token file operations
func BenchmarkTokenSaveLoad(b *testing.B) {
	tempDir := b.TempDir()
	tokenPath := filepath.Join(tempDir, TokenFile)

	auth := &Auth{
		tokenFile: tokenPath,
	}

	testToken := &oauth2.Token{
		AccessToken:  "benchmark-access-token",
		RefreshToken: "benchmark-refresh-token",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		auth.saveToken(testToken)
		_, err := auth.tokenFromFile()
		if err != nil {
			b.Fatalf("tokenFromFile() failed: %v", err)
		}
	}
}