# mattermost-plugin-disable-dm

A Mattermost plugin to disable direct messages and group chats with team-based whitelist support.

## Features

- **Block Direct Messages**: Prevent users from sending direct messages
- **Block Group Chats**: Prevent users from sending group chat messages
- **Team-Based Whitelist**: Allow specific teams to bypass restrictions (perfect for staff/admin teams)
- **Customizable Messages**: Configure the message shown to blocked users

## Use Cases

This plugin is ideal for organizations that need to control direct messaging, such as:
- **Correctional Facilities**: Block inmate-to-inmate DMs while allowing staff communication
- **Educational Institutions**: Restrict student DMs while permitting teacher communication
- **Enterprises**: Enforce communication policies by team membership

## Requirements

- Mattermost Server v9.0.0 or later
- Go v1.21 or later
- Make

## Installation and Setup

### Building the Plugin

1. Clone the repository:
```bash
git clone https://github.com/Brightscout/mattermost-plugin-disable-dm.git
cd mattermost-plugin-disable-dm
```

2. Build the plugin:
```bash
make dist
```

This will produce a `mattermost-plugin-disable-dm-v2.0.0.tar.gz` file in the `/dist` directory containing binaries for all platforms (Linux, macOS, Windows).

### Installing in Mattermost

1. Go to **System Console > Plugins > Plugin Management**
2. Upload the `.tar.gz` file from the `dist` directory
3. Enable the plugin

## Configuration

Navigate to **System Console > Plugins > Disable DM** to configure:

### Basic Settings

- **Reject DMs**: Enable to block direct messages (default: true)
- **Reject Group Chats**: Enable to block group chat messages (default: true)
- **Rejection Message**: Message shown to users when blocked (customizable)

### Team Whitelist Settings

- **Enable Team Whitelist**: When enabled, only users in whitelisted teams can send DMs/group messages
- **Whitelisted Team Names**: Comma-separated list of team names (not display names)

#### Example Configuration

For a prison system where staff should communicate freely:

1. Set **Reject DMs**: `true`
2. Set **Reject Group Chats**: `true`
3. Set **Enable Team Whitelist**: `true`
4. Set **Whitelisted Team Names**: `staff-team,admin-team`
5. Set **Rejection Message**: `Direct messaging between inmates is not permitted. Please use designated channels.`

**Result**: Only members of `staff-team` or `admin-team` can send DMs and group messages. All other users will be blocked.

#### Finding Team Names

Team names are the URL-friendly identifiers, not the display names:
- Navigate to a team in Mattermost
- Check the URL: `https://your-server.com/team-name/channels/...`
- Use `team-name` (from the URL) in the configuration, not the display name

## How It Works

1. The plugin intercepts all messages before they're posted using the `MessageWillBePosted` hook
2. **For Direct Messages (when blocking is enabled):**
   - If team whitelist is disabled: All DMs are blocked
   - If team whitelist is enabled: Only users in whitelisted teams can send DMs
3. **For Group Chats (when blocking is enabled):**
   - If team whitelist is disabled: All group messages are blocked
   - If team whitelist is enabled: Messages are only allowed if ALL members of the group are in whitelisted teams
   - This prevents mixed staff-inmate group chats
4. Blocked users see an ephemeral message (visible only to them) with the rejection reason
5. The message is not saved or visible to other users

## Important Things to Know

### What Gets Blocked
- **Messages**: All DM and group chat messages are blocked server-side for non-whitelisted users
- **Security**: The blocking happens at the message level, preventing communication

### What Doesn't Get Blocked
- **UI Channel Opening**: Users can still open DM channels in the UI (but cannot send messages)
- **Typing Indicators**: Users may see "User is typing..." indicators in DM channels
- **Channel URLs**: Direct URLs to DM channels remain accessible (but messages still blocked)

### Why This Is Acceptable
This is the standard behavior for message blocking plugins and is **secure** because:
- The actual **message sending is blocked server-side** - this is the security boundary
- Typing indicators without message delivery do not enable communication
- Users quickly learn DMs are disabled when they cannot send messages
- UI-level blocking would require a webapp plugin component and only provides cosmetic improvements
- The approach matches how enterprise platforms (Slack, Teams) implement similar restrictions

## Development

### Project Structure

```
mattermost-plugin-disable-dm/
├── server/
│   ├── plugin.go           # Main plugin logic and hooks
│   └── config/
│       ├── main.go          # Configuration management
│       └── manifest.go      # Plugin constants
├── plugin.json              # Plugin manifest and settings
├── Makefile                 # Build automation
└── README.md
```

### Building for Development

```bash
# Build plugin
make

# Run tests (if available)
make test

# Clean build artifacts
make clean
```

## Troubleshooting

### Plugin Won't Activate

- Check Mattermost server logs for errors
- Verify minimum server version (v9.0.0+)
- Ensure plugin uploads are enabled in System Console

### Whitelist Not Working

- Verify you're using team **names** (from URLs), not display names
- Check that users are actually members of the specified teams
- Review Mattermost server logs for permission errors

### Messages Still Going Through

- Confirm the plugin is enabled and activated
- Check that the relevant settings (Reject DMs / Reject Group Chats) are enabled
- Verify channel type (public/private channels are not affected, only DMs and group chats)

## License

See [LICENSE](LICENSE) file for details.

---

Made with &#9829; by [Brightscout](http://www.brightscout.com)
