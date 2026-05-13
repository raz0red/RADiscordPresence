package presence

// StatusSnapshot captures the current state of the Worker, reported on every
// poll cycle to optional consumers such as the web UI.
type StatusSnapshot struct {
	IsPlaying        bool
	GameTitle        string
	ConsoleName      string
	CoverArtURL      string
	RichPresenceMsg  string
	Achievements     string // e.g. "14/47 (Hardcore)" or ""
	GameStartTime    int64
	DiscordConnected bool
	RAConnected      bool
	Username         string
	GameID           int
	// User profile info
	AvatarURL           string
	TotalPoints         int
	TotalSoftcorePoints int
	TotalTruePoints     int
	Rank                int
	TotalRanked         int
	// LastError holds the most recent poll error message, or "" on success.
	LastError string
}

// StatusUpdater receives status snapshots from the Worker on each poll cycle.
// web.Hub implements this interface.
type StatusUpdater interface {
	UpdateStatus(StatusSnapshot)
}
