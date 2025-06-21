package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	
	// OAuth2 client configuration (embedded in binary)
	// Replace with your actual OAuth2 credentials
	ClientID     = "your-client-id.apps.googleusercontent.com"
	ClientSecret = "your-client-secret"
	RedirectURL  = "urn:ietf:wg:oauth:2.0:oob"
)

type Auth struct {
	config    *oauth2.Config
	tokenFile string
}

func NewAuth() (*Auth, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("unable to create config directory: %v", err)
	}

	tokenPath := filepath.Join(configDir, TokenFile)

	config := &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		RedirectURL:  RedirectURL,
		Scopes:       []string{drive.DriveScope},
		Endpoint:     google.Endpoint,
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

func (a *Auth) getTokenFromWeb() (*oauth2.Token, error) {
	authURL := a.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	tok, err := a.config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return tok, nil
}

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

func (a *Auth) saveToken(token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", a.tokenFile)
	f, err := os.OpenFile(a.tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

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
