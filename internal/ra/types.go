package ra

import "time"

// PresenceTimeout is how stale a RichPresenceMsgDate can be before
// we consider the user no longer actively playing.
const PresenceTimeout = 130 * time.Second

// UserSummary is the relevant subset of the API_GetUserSummary response.
type UserSummary struct {
	User                string `json:"User"`
	LastGameID          int    `json:"LastGameID"`
	RichPresenceMsg     string `json:"RichPresenceMsg"`
	RichPresenceMsgDate string `json:"RichPresenceMsgDate"`
	Status              string `json:"Status"`
	UserPic             string `json:"UserPic"`
	TotalPoints         int    `json:"TotalPoints"`
	TotalSoftcorePoints int    `json:"TotalSoftcorePoints"`
	TotalTruePoints     int    `json:"TotalTruePoints"`
	Rank                int    `json:"Rank"`
	TotalRanked         int    `json:"TotalRanked"`
}

// AvatarURL returns the full URL for the user's profile picture.
func (u UserSummary) AvatarURL() string {
	if u.UserPic == "" {
		return ""
	}
	return "https://media.retroachievements.org" + u.UserPic
}

// IsActive returns true if the user has fresh rich presence data.
func (u UserSummary) IsActive() bool {
	if u.LastGameID == 0 || u.RichPresenceMsg == "" {
		return false
	}
	t, err := time.Parse("2006-01-02 15:04:05", u.RichPresenceMsgDate)
	if err != nil {
		return false
	}
	return time.Since(t.UTC()) < PresenceTimeout
}

// Game is the relevant subset of the API_GetGame response.
type Game struct {
	Title       string `json:"Title"`
	ConsoleName string `json:"ConsoleName"`
	ImageIcon   string `json:"ImageIcon"`
}

// ArtURL returns the full URL for the game's cover art image.
func (g Game) ArtURL() string {
	if g.ImageIcon == "" {
		return ""
	}
	return "https://media.retroachievements.org" + g.ImageIcon
}

// UserProgress is the relevant subset of the API_GetUserProgress response.
type UserProgress struct {
	NumPossibleAchievements int `json:"NumPossibleAchievements"`
	NumAchieved             int `json:"NumAchieved"`
	NumAchievedHardcore     int `json:"NumAchievedHardcore"`
}
