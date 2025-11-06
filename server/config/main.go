package config

import (
	"errors"
	"strings"

	"github.com/mattermost/mattermost/server/public/plugin"
	"go.uber.org/atomic"
)

const (
	HeaderMattermostUserID = "Mattermost-User-Id"
)

var (
	config     atomic.Value
	Mattermost plugin.API
)

type Configuration struct {
	RejectDMs           bool   `json:"RejectDMs"`
	RejectGroupChats    bool   `json:"RejectGroupChats"`
	RejectionMessage    string `json:"RejectionMessage"`
	EnableTeamWhitelist bool   `json:"EnableTeamWhitelist"`
	WhitelistedTeams    string `json:"WhitelistedTeams"`

	// Processed field - not from JSON
	whitelistedTeamSet map[string]bool
}

func GetConfig() *Configuration {
	return config.Load().(*Configuration)
}

func SetConfig(c *Configuration) {
	config.Store(c)
}

func (c *Configuration) ProcessConfiguration() error {
	// any post-processing on configurations goes here

	c.RejectionMessage = strings.TrimSpace(c.RejectionMessage)
	c.WhitelistedTeams = strings.TrimSpace(c.WhitelistedTeams)

	// Parse the comma-separated team list into a set for fast lookup
	c.whitelistedTeamSet = make(map[string]bool)
	if c.WhitelistedTeams != "" {
		teams := strings.Split(c.WhitelistedTeams, ",")
		for _, team := range teams {
			teamName := strings.TrimSpace(team)
			if teamName != "" {
				c.whitelistedTeamSet[teamName] = true
			}
		}
	}

	return nil
}

func (c *Configuration) IsValid() error {
	// Add config validations here.
	// Check for required fields, formats, etc.

	if (c.RejectDMs || c.RejectGroupChats) && c.RejectionMessage == "" {
		return errors.New("rejection message cannot be empty")
	}

	if c.EnableTeamWhitelist && c.WhitelistedTeams == "" {
		return errors.New("team whitelist is enabled but no teams are specified")
	}

	return nil
}

// IsUserInWhitelistedTeam checks if a user is a member of any whitelisted team
// Returns true if the user should be ALLOWED to send DMs/group messages
func (c *Configuration) IsUserInWhitelistedTeam(userID string) bool {
	// If whitelist is not enabled, block everyone
	if !c.EnableTeamWhitelist {
		return false
	}

	// If no teams are whitelisted, no one can send DMs
	if len(c.whitelistedTeamSet) == 0 {
		return false
	}

	// Get all teams the user is a member of
	teams, appErr := Mattermost.GetTeamsForUser(userID)
	if appErr != nil {
		Mattermost.LogError("Failed to get teams for user", "user_id", userID, "error", appErr.Error())
		return false
	}

	// Check if user is in any of the whitelisted teams
	for _, team := range teams {
		if c.whitelistedTeamSet[team.Name] {
			return true
		}
	}

	return false
}
