// Package ra implements the RetroAchievements Web API client.
package ra

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/raz0red/radpresence/internal/buildinfo"
)

const baseURL = "https://retroachievements.org/API"

func userAgent() string {
	return fmt.Sprintf("RADPresence/%s (https://github.com/raz0red/radpresence)", buildinfo.Version)
}

// Client calls the RetroAchievements Web API.
type Client struct {
	username string
	apiKey   string
	http     *http.Client
}

// New returns a Client authenticated with the given credentials.
func New(username, apiKey string) *Client {
	return &Client{
		username: username,
		apiKey:   apiKey,
		http:     &http.Client{Timeout: 15 * time.Second},
	}
}

// UpdateCredentials replaces the username and API key used for subsequent requests.
func (c *Client) UpdateCredentials(username, apiKey string) {
	c.username = username
	c.apiKey = apiKey
}

func (c *Client) get(endpoint string, params url.Values) ([]byte, error) {
	params.Set("y", c.apiKey)
	u := fmt.Sprintf("%s/%s.php?%s", baseURL, endpoint, params.Encode())
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", endpoint, err)
	}
	req.Header.Set("User-Agent", userAgent())
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", endpoint, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: HTTP %d", endpoint, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: reading body: %w", endpoint, err)
	}
	return body, nil
}

// GetUserSummary fetches the user's current session summary.
func (c *Client) GetUserSummary() (UserSummary, error) {
	data, err := c.get("API_GetUserSummary", url.Values{
		"u": {c.username}, "g": {"0"}, "a": {"0"},
	})
	if err != nil {
		return UserSummary{}, err
	}
	var s UserSummary
	return s, json.Unmarshal(data, &s)
}

// GetGame fetches metadata for the given game ID.
func (c *Client) GetGame(gameID int) (Game, error) {
	data, err := c.get("API_GetGame", url.Values{
		"z": {c.username}, "i": {strconv.Itoa(gameID)},
	})
	if err != nil {
		return Game{}, err
	}
	var g Game
	return g, json.Unmarshal(data, &g)
}

// GetUserProgress fetches achievement progress for the given game.
func (c *Client) GetUserProgress(gameID int) (UserProgress, error) {
	data, err := c.get("API_GetUserProgress", url.Values{
		"u": {c.username}, "i": {strconv.Itoa(gameID)},
	})
	if err != nil {
		return UserProgress{}, err
	}
	// Response is a map keyed by game ID string.
	var raw map[string]UserProgress
	if err := json.Unmarshal(data, &raw); err != nil {
		return UserProgress{}, fmt.Errorf("parsing UserProgress: %w", err)
	}
	p, ok := raw[strconv.Itoa(gameID)]
	if !ok {
		return UserProgress{}, fmt.Errorf("game %d missing from progress response", gameID)
	}
	return p, nil
}
