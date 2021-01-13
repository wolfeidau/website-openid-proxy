package stores

import "time"

// SessionRecord is used to store information about the currently authenticated user session
type SessionRecord struct {
	CreatedAt    time.Time
	AccessToken  string
	IDToken      string
	RefreshToken string
}
