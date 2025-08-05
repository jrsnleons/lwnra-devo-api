package facebook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// TokenManager handles Facebook token refresh
type TokenManager struct {
	appID        string
	appSecret    string
	currentToken string
	pageID       string
}

// TokenResponse represents Facebook token API response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}

// NewTokenManager creates a new token manager
func NewTokenManager(appID, appSecret, initialToken string) *TokenManager {
	return &TokenManager{
		appID:        appID,
		appSecret:    appSecret,
		currentToken: initialToken,
		pageID:       "164421594332429", // Living Word NRA page ID
	}
}

// RefreshToken exchanges current token for a long-lived token
func (tm *TokenManager) RefreshToken() (string, error) {
	// Exchange current token for long-lived token
	apiURL := fmt.Sprintf(
		"https://graph.facebook.com/oauth/access_token?grant_type=fb_exchange_token&client_id=%s&client_secret=%s&fb_exchange_token=%s",
		tm.appID,
		tm.appSecret,
		url.QueryEscape(tm.currentToken),
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Facebook API error (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	tm.currentToken = tokenResp.AccessToken
	return tokenResp.AccessToken, nil
}

// GetPageToken gets a long-lived page access token
func (tm *TokenManager) GetPageToken(pageID string) (string, error) {
	apiURL := fmt.Sprintf(
		"https://graph.facebook.com/%s?fields=access_token&access_token=%s",
		pageID,
		url.QueryEscape(tm.currentToken),
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to get page token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Facebook API error (status %d): %s", resp.StatusCode, string(body))
	}

	var pageResp struct {
		AccessToken string `json:"access_token"`
		ID          string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&pageResp); err != nil {
		return "", fmt.Errorf("failed to decode page response: %w", err)
	}

	return pageResp.AccessToken, nil
}

// ValidateToken checks if current token is still valid
func (tm *TokenManager) ValidateToken() (bool, error) {
	apiURL := fmt.Sprintf(
		"https://graph.facebook.com/me?access_token=%s",
		url.QueryEscape(tm.currentToken),
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// GetTokenInfo returns information about the current token
func (tm *TokenManager) GetTokenInfo() (map[string]interface{}, error) {
	// Use debug_token endpoint for more detailed info
	apiURL := fmt.Sprintf(
		"https://graph.facebook.com/debug_token?input_token=%s&access_token=%s",
		url.QueryEscape(tm.currentToken),
		url.QueryEscape(tm.currentToken),
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get token info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Facebook API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode token info: %w", err)
	}

	return response.Data, nil
}

// IsTokenExpiringSoon checks if token expires within 7 days
func (tm *TokenManager) IsTokenExpiringSoon() (bool, error) {
	tokenInfo, err := tm.GetTokenInfo()
	if err != nil {
		return false, err
	}

	expiresAtFloat, exists := tokenInfo["expires_at"].(float64)
	if !exists {
		// Token might be permanent or very long-lived
		return false, nil
	}

	expiresAt := int64(expiresAtFloat)
	if expiresAt == 0 {
		// Token never expires
		return false, nil
	}

	// Check if expires within 7 days (604800 seconds)
	currentTime := time.Now().Unix()
	return (expiresAt - currentTime) < 604800, nil
}

// GetCurrentToken returns the current token
func (tm *TokenManager) GetCurrentToken() string {
	return tm.currentToken
}

// SetCurrentToken updates the current token
func (tm *TokenManager) SetCurrentToken(token string) {
	tm.currentToken = token
}
