# Mattermost Disable DM Plugin - Validation Test Plan

**Plugin Version:** 2.0.0
**Date:** 2025-11-06
**Tester:** _____________

## Prerequisites

- [ ] Mattermost Server v9.0.0 or later
- [ ] Plugin uploaded and activated
- [ ] Debug logging enabled in Mattermost config
- [ ] Access to `mattermost.log` file

## Test Environment Setup

### Enable Debug Logging

1. Navigate to: **System Console > Environment > Logging**
2. Set **File Log Level** to `DEBUG`
3. Save changes

### Plugin Configuration Access

Navigate to: **System Console > Plugins > Disable DM**

---

## Test Scenario 1: Basic DM Blocking (No Whitelist)

### Configuration
```
Reject DMs: ✅ true
Reject Group Chats: ✅ true
Enable Team Whitelist: ❌ false
Whitelisted Team Names: (empty)
Rejection Message: "Direct messages have been disabled by an administrator."
```

### Test Steps

1. **Save Configuration**
   - [ ] Configuration saved successfully
   - [ ] Check logs for: `"Configuration change detected"`
   - [ ] Check logs for: `"Configuration applied successfully"`
   - [ ] Verify log shows: `"enable_whitelist": false`

2. **Test DM Blocking**
   - [ ] User A sends DM to User B
   - [ ] Expected: Message is blocked
   - [ ] Expected: User A sees ephemeral message with rejection text
   - [ ] Expected: User B does NOT receive message
   - [ ] Expected: User B does NOT receive desktop notification
   - [ ] Check logs for: `"Blocking message"` with `"channel_type": "D"`

3. **Test Group Chat Blocking**
   - [ ] User A sends message in group chat (3+ users)
   - [ ] Expected: Message is blocked
   - [ ] Expected: User A sees ephemeral message
   - [ ] Expected: Other users do NOT receive message
   - [ ] Check logs for: `"Blocking message"` with `"channel_type": "G"`

4. **Test Public Channel (Should NOT Block)**
   - [ ] User A sends message in public channel
   - [ ] Expected: Message goes through normally
   - [ ] Expected: No blocking occurs
   - [ ] Check logs: `"should_reject": false`

### Expected Log Output
```
[INFO] Configuration change detected
[INFO] Configuration loaded reject_dms=true reject_group_chats=true enable_whitelist=false
[INFO] Configuration applied successfully
[DEBUG] MessageWillBePosted called post_id=... channel_type=D
[DEBUG] Rejection check should_reject=true is_direct=true
[DEBUG] Whitelist disabled - blocking user
[INFO] Blocking message user_id=... channel_type=D reason="User not in whitelisted team"
```

**Result:** ⬜ PASS / ⬜ FAIL
**Notes:**

---

## Test Scenario 2: Team Whitelist - Staff Only

### Configuration
```
Reject DMs: ✅ true
Reject Group Chats: ✅ true
Enable Team Whitelist: ✅ true
Whitelisted Team Names: "staff-team"
Rejection Message: "Direct messaging between inmates is not permitted."
```

### Test Environment Setup

Create test teams:
- [ ] Team "staff-team" exists
- [ ] Team "inmates-team" exists
- [ ] User "Staff1" is member of "staff-team" only
- [ ] User "Inmate1" is member of "inmates-team" only
- [ ] User "Inmate2" is member of "inmates-team" only

### Test Steps

1. **Save Configuration**
   - [ ] Configuration saved successfully
   - [ ] Check logs for: `"enable_whitelist": true`
   - [ ] Check logs for: `"whitelisted_teams": "staff-team"`
   - [ ] Check logs for: `"whitelist_team_count": 1`

2. **Test Staff → Staff DM (Should ALLOW)**
   - [ ] Staff1 sends DM to another staff member
   - [ ] Expected: Message goes through
   - [ ] Expected: No rejection message
   - [ ] Check logs for: `"User allowed - member of whitelisted team"` with `"team_name": "staff-team"`

3. **Test Inmate → Inmate DM (Should BLOCK)**
   - [ ] Inmate1 sends DM to Inmate2
   - [ ] Expected: Message is blocked
   - [ ] Expected: Inmate1 sees rejection message
   - [ ] Expected: Inmate2 does NOT receive message or notification
   - [ ] Check logs for: `"User not in any whitelisted team - blocking"`

4. **Test Inmate → Staff DM (Should BLOCK sender)**
   - [ ] Inmate1 sends DM to Staff1
   - [ ] Expected: Message is blocked
   - [ ] Expected: Inmate1 sees rejection message
   - [ ] Check logs for: `"Blocking message"` for Inmate1's user_id

5. **Test Group Chat - Mixed Users (Block non-staff)**
   - [ ] Create group chat with Staff1, Inmate1, Inmate2
   - [ ] Inmate1 tries to send message
   - [ ] Expected: Blocked
   - [ ] Staff1 tries to send message
   - [ ] Expected: Allowed

### Expected Log Output
```
[INFO] Configuration change detected
[INFO] Configuration loaded enable_whitelist=true whitelisted_teams="staff-team"
[DEBUG] IsUserInWhitelistedTeam called enable_whitelist=true whitelist_count=1
[DEBUG] User team membership team_count=1
[DEBUG] Checking team team_name="staff-team" is_whitelisted=true
[INFO] User allowed - member of whitelisted team team_name="staff-team"
```

**Result:** ⬜ PASS / ⬜ FAIL
**Notes:**

---

## Test Scenario 3: Multiple Whitelisted Teams

### Configuration
```
Reject DMs: ✅ true
Reject Group Chats: ✅ true
Enable Team Whitelist: ✅ true
Whitelisted Team Names: "staff-team,admin-team,security-team"
```

### Test Environment Setup

- [ ] User "Admin1" is member of "admin-team"
- [ ] User "Security1" is member of "security-team"
- [ ] User "Inmate1" is member of "inmates-team"

### Test Steps

1. **Save Configuration**
   - [ ] Check logs show all 3 teams parsed
   - [ ] Verify: `"whitelist_team_count": 3` (Note: this counts string length, not teams)

2. **Test Each Whitelisted Team**
   - [ ] Staff1 (staff-team) → Staff2: Allowed
   - [ ] Admin1 (admin-team) → Admin2: Allowed
   - [ ] Security1 (security-team) → Security2: Allowed
   - [ ] Inmate1 → Inmate2: Blocked

3. **Test Cross-Team (Whitelisted to Whitelisted)**
   - [ ] Staff1 → Admin1: Allowed
   - [ ] Admin1 → Security1: Allowed

**Result:** ⬜ PASS / ⬜ FAIL
**Notes:**

---

## Test Scenario 4: Disable DMs Only (Allow Group Chats)

### Configuration
```
Reject DMs: ✅ true
Reject Group Chats: ❌ false
Enable Team Whitelist: ❌ false
```

### Test Steps

1. **Test DM**
   - [ ] User A → User B DM: Blocked
   - [ ] Check logs: `"is_direct": true, "should_reject": true`

2. **Test Group Chat**
   - [ ] User A sends in group chat: Allowed
   - [ ] Check logs: `"is_group": true, "should_reject": false`

**Result:** ⬜ PASS / ⬜ FAIL
**Notes:**

---

## Test Scenario 5: Configuration Edge Cases

### Test 5a: Empty Whitelist Team String
```
Enable Team Whitelist: ✅ true
Whitelisted Team Names: ""
```

- [ ] Configuration should fail validation
- [ ] Check logs for: `"team whitelist is enabled but no teams are specified"`

### Test 5b: Whitespace in Team Names
```
Whitelisted Team Names: " staff-team , admin-team , security-team "
```

- [ ] Configuration should parse correctly (whitespace trimmed)
- [ ] Teams should be identified properly

### Test 5c: Invalid Team Name
```
Whitelisted Team Names: "nonexistent-team"
```

- [ ] User tries to send DM
- [ ] Expected: Blocked (no users in this team)
- [ ] Check logs: `"User not in any whitelisted team"`

**Result:** ⬜ PASS / ⬜ FAIL
**Notes:**

---

## Test Scenario 6: Desktop Notifications

### Critical Test
When a message is blocked, notifications should NOT be sent.

1. **User A (blocked) sends DM to User B**
2. **Verification Points:**
   - [ ] User B does NOT see desktop notification
   - [ ] User B does NOT see unread badge
   - [ ] User B does NOT see message in channel
   - [ ] Only User A sees ephemeral rejection message

**Result:** ⬜ PASS / ⬜ FAIL
**Notes:**

---

## Troubleshooting Guide

### Issue: Messages Not Being Blocked

Check:
1. Plugin is activated: System Console > Plugins > Plugin Management
2. Configuration saved properly
3. Log shows: `"MessageWillBePosted called"` when sending message
4. Channel type detection: Check if `"channel_type"` is `"D"` or `"G"`
5. Whitelist logic: Check `"is_whitelisted"` value in logs

### Issue: Notifications Still Sent for Blocked Messages

This indicates the hook is NOT preventing the message. Check:
1. Hook is returning `(nil, "rejection message")` not `(post, "")`
2. Log shows: `"Blocking message"` entry
3. Verify second return value is the rejection message string

### Issue: Configuration Not Loading

Check:
1. Log shows: `"Configuration change detected"`
2. Log shows: `"Configuration loaded"` with correct values
3. No errors in: `LoadPluginConfiguration` or `ProcessConfiguration`

---

## Log Analysis Commands

### View Plugin Logs
```bash
# Real-time monitoring
tail -f /path/to/mattermost.log | grep -i "disable-dm\|MessageWillBePosted\|Blocking message"

# Search for specific user's actions
grep "user_id=<USER_ID>" mattermost.log | grep -i "disable-dm"

# Check configuration loading
grep "Configuration loaded" mattermost.log | tail -5

# Find blocking events
grep "Blocking message" mattermost.log
```

---

## Test Summary

| Scenario | Result | Issues Found |
|----------|--------|--------------|
| 1. Basic DM Blocking | ⬜ PASS / ⬜ FAIL | |
| 2. Team Whitelist | ⬜ PASS / ⬜ FAIL | |
| 3. Multiple Teams | ⬜ PASS / ⬜ FAIL | |
| 4. DMs Only | ⬜ PASS / ⬜ FAIL | |
| 5. Edge Cases | ⬜ PASS / ⬜ FAIL | |
| 6. Notifications | ⬜ PASS / ⬜ FAIL | |

**Overall Result:** ⬜ PASS / ⬜ FAIL

**Critical Issues Found:**

**Recommendations:**
