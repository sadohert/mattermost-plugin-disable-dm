package main

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/Brightscout/mattermost-plugin-disable-dm/server/config"
)

type Plugin struct {
	plugin.MattermostPlugin
}

func (p *Plugin) OnActivate() error {
	config.Mattermost = p.API
	config.Mattermost.LogInfo("Plugin activating")

	if err := p.OnConfigurationChange(); err != nil {
		config.Mattermost.LogError("Failed to activate plugin: " + err.Error())
		return err
	}

	config.Mattermost.LogInfo("Plugin activated successfully")
	return nil
}

func (p *Plugin) OnConfigurationChange() error {
	if config.Mattermost != nil {
		config.Mattermost.LogInfo("Configuration change detected")

		var configuration config.Configuration

		if err := config.Mattermost.LoadPluginConfiguration(&configuration); err != nil {
			config.Mattermost.LogError("Error in LoadPluginConfiguration: " + err.Error())
			return err
		}

		config.Mattermost.LogInfo("Configuration loaded",
			"reject_dms", configuration.RejectDMs,
			"reject_group_chats", configuration.RejectGroupChats,
			"enable_whitelist", configuration.EnableTeamWhitelist,
			"whitelisted_teams", configuration.WhitelistedTeams,
			"rejection_message", configuration.RejectionMessage)

		if err := configuration.ProcessConfiguration(); err != nil {
			config.Mattermost.LogError("Error in ProcessConfiguration: " + err.Error())
			return err
		}

		config.Mattermost.LogInfo("Configuration processed",
			"whitelist_team_count", len(configuration.WhitelistedTeams))

		if err := configuration.IsValid(); err != nil {
			config.Mattermost.LogError("Error in Validating Configuration: " + err.Error())
			return err
		}

		config.SetConfig(&configuration)
		config.Mattermost.LogInfo("Configuration applied successfully")
	}
	return nil
}

func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	conf := config.GetConfig()

	channel, appError := config.Mattermost.GetChannel(post.ChannelId)
	if appError != nil {
		config.Mattermost.LogError("Failed to get channel", "channel_id", post.ChannelId, "error", appError.Error())
		return post, ""
	}

	// Check if this is a DM that should be blocked
	if channel.Type == model.ChannelTypeDirect && conf.RejectDMs {
		// For DMs, check if sender is whitelisted
		isWhitelisted := conf.IsUserInWhitelistedTeam(post.UserId)
		if !isWhitelisted {
			config.Mattermost.LogInfo("Blocked DM",
				"user_id", post.UserId,
				"channel_type", string(channel.Type))

			config.Mattermost.SendEphemeralPost(post.UserId, &model.Post{
				Message:   conf.RejectionMessage,
				ChannelId: post.ChannelId,
			})
			return nil, conf.RejectionMessage
		}
	}

	// Check if this is a group chat that should be blocked
	if channel.Type == model.ChannelTypeGroup && conf.RejectGroupChats {
		// For group chats, check if ALL members are whitelisted
		members, appErr := config.Mattermost.GetChannelMembers(post.ChannelId, 0, 100)
		if appErr != nil {
			config.Mattermost.LogError("Failed to get channel members", "channel_id", post.ChannelId, "error", appErr.Error())
			return post, ""
		}

		// Check if all members (including sender) are whitelisted
		for _, member := range members {
			if !conf.IsUserInWhitelistedTeam(member.UserId) {
				config.Mattermost.LogInfo("Blocked group message - not all members whitelisted",
					"sender_id", post.UserId,
					"non_whitelisted_user", member.UserId,
					"channel_type", string(channel.Type))

				config.Mattermost.SendEphemeralPost(post.UserId, &model.Post{
					Message:   conf.RejectionMessage,
					ChannelId: post.ChannelId,
				})
				return nil, conf.RejectionMessage
			}
		}
	}

	// Allow the message to be posted
	return post, ""
}

func main() {
	plugin.ClientMain(&Plugin{})
}
