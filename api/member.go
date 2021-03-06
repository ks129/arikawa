package api

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/utils/httputil"
	"github.com/diamondburned/arikawa/utils/json/option"
)

// Member returns a guild member object for the specified user..
func (c *Client) Member(guildID, userID discord.Snowflake) (*discord.Member, error) {
	var m *discord.Member
	return m, c.RequestJSON(&m, "GET", EndpointGuilds+guildID.String()+"/members/"+userID.String())
}

// Members returns a list of members of the guild with the passed id. This
// method automatically paginates until it reaches the passed limit, or, if the
// limit is set to 0, has fetched all members within the passed range.
//
// As the underlying endpoint has a maximum of 1000 members per request, at
// maximum a total of limit/1000 rounded up requests will be made, although
// they may be less, if no more members are available.
//
// When fetching the members, those with the smallest ID will be fetched first.
func (c *Client) Members(guildID discord.Snowflake, limit uint) ([]discord.Member, error) {
	return c.MembersAfter(guildID, 0, limit)
}

// MembersAfter returns a list of members of the guild with the passed id. This
// method automatically paginates until it reaches the passed limit, or, if the
// limit is set to 0, has fetched all members within the passed range.
//
// As the underlying endpoint has a maximum of 1000 members per request, at
// maximum a total of limit/1000 rounded up requests will be made, although
// they may be less, if no more members are available.
func (c *Client) MembersAfter(
	guildID, after discord.Snowflake, limit uint) ([]discord.Member, error) {

	var mems []discord.Member

	const hardLimit int = 1000

	unlimited := limit == 0

	for fetch := uint(hardLimit); limit > 0 || unlimited; fetch = uint(hardLimit) {
		if limit > 0 {
			if fetch > limit {
				fetch = limit
			}
			limit -= fetch
		}

		m, err := c.membersAfter(guildID, after, fetch)
		if err != nil {
			return mems, err
		}
		mems = append(mems, m...)

		// There aren't any to fetch, even if this is less than limit.
		if len(mems) < hardLimit {
			break
		}

		after = mems[hardLimit-1].User.ID
	}

	return mems, nil
}

func (c *Client) membersAfter(
	guildID, after discord.Snowflake, limit uint) ([]discord.Member, error) {

	switch {
	case limit == 0:
		limit = 0
	case limit > 1000:
		limit = 1000
	}

	var param struct {
		After discord.Snowflake `schema:"after,omitempty"`
		Limit uint              `schema:"limit"`
	}

	param.Limit = limit
	param.After = after

	var mems []discord.Member
	return mems, c.RequestJSON(
		&mems, "GET",
		EndpointGuilds+guildID.String()+"/members",
		httputil.WithSchema(c, param),
	)
}

// https://discord.com/developers/docs/resources/guild#add-guild-member-json-params
type AddMemberData struct {
	// Token is an oauth2 access token granted with the guilds.join to the
	// bot's application for the user you want to add to the guild.
	Token string `json:"access_token"`
	// Nick is the value to set users nickname to.
	//
	// Requires MANAGE_NICKNAMES.
	Nick option.String `json:"nick,omitempty"`
	// Roles is an array of role ids the member is assigned.
	//
	// Requires MANAGE_ROLES.
	Roles *[]discord.Snowflake `json:"roles,omitempty"`
	// Mute specifies whether the user is muted in voice channels.
	//
	// Requires MUTE_MEMBERS.
	Mute option.Bool `json:"mute,omitempty"`
	// Deaf specifies whether the user is deafened in voice channels.
	//
	// Requires DEAFEN_MEMBERS.
	Deaf option.Bool `json:"deaf,omitempty"`
}

// AddMember adds a user to the guild, provided you have a valid oauth2 access
// token for the user with the guilds.join scope. Returns a 201 Created with
// the guild member as the body, or 204 No Content if the user is already a
// member of the guild.
//
// Fires a Guild Member Add Gateway event.
//
// The Authorization header must be a Bot token (belonging to the same
// application used for authorization), and the bot must be a member of the
// guild with CREATE_INSTANT_INVITE permission.
func (c *Client) AddMember(
	guildID, userID discord.Snowflake, data AddMemberData) (*discord.Member, error) {
	var mem *discord.Member
	return mem, c.RequestJSON(
		&mem, "PUT",
		EndpointGuilds+guildID.String()+"/members/"+userID.String(),
		httputil.WithJSONBody(data),
	)
}

// https://discord.com/developers/docs/resources/guild#add-guild-member-json-params
type ModifyMemberData struct {
	// Nick is the value to set users nickname to.
	//
	// Requires MANAGE_NICKNAMES.
	Nick option.String `json:"nick,omitempty"`
	// Roles is an array of role ids the member is assigned.
	//
	// Requires MANAGE_ROLES.
	Roles *[]discord.Snowflake `json:"roles,omitempty"`
	// Mute specifies whether the user is muted in voice channels.
	//
	// Requires MUTE_MEMBERS.
	Mute option.Bool `json:"mute,omitempty"`
	// Deaf specifies whether the user is deafened in voice channels.
	//
	// Requires DEAFEN_MEMBERS.
	Deaf option.Bool `json:"deaf,omitempty"`

	// Voice channel is the id of channel to move user to (if they are
	// connected to voice).
	//
	// Requires MOVE_MEMBER
	VoiceChannel discord.Snowflake `json:"channel_id,omitempty"`
}

// ModifyMember modifies attributes of a guild member. If the channel_id is set
// to null, this will force the target user to be disconnected from voice.
//
// Fires a Guild Member Update Gateway event.
func (c *Client) ModifyMember(guildID, userID discord.Snowflake, data ModifyMemberData) error {

	return c.FastRequest(
		"PATCH",
		EndpointGuilds+guildID.String()+"/members/"+userID.String(),
		httputil.WithJSONBody(data),
	)
}

// PruneCount returns the number of members that would be removed in a prune
// operation. Days must be 1 or more, default 7.
//
// Requires KICK_MEMBERS.
func (c *Client) PruneCount(guildID discord.Snowflake, days uint) (uint, error) {
	if days == 0 {
		days = 7
	}

	var param struct {
		Days uint `schema:"days"`
	}

	param.Days = days

	var resp struct {
		Pruned uint `json:"pruned"`
	}

	return resp.Pruned, c.RequestJSON(
		&resp, "GET",
		EndpointGuilds+guildID.String()+"/prune",
		httputil.WithSchema(c, param),
	)
}

// Prune begins a prune. Days must be 1 or more, default 7.
//
// Requires KICK_MEMBERS.
func (c *Client) Prune(guildID discord.Snowflake, days uint) error {
	if days == 0 {
		days = 7
	}

	var param struct {
		Days     uint `schema:"days"`
		RetCount bool `schema:"compute_prune_count"`
	}

	param.Days = days
	param.RetCount = false

	return c.FastRequest(
		"POST",
		EndpointGuilds+guildID.String()+"/prune",
		httputil.WithSchema(c, param),
	)
}

// PruneWithCounts returns the number of members that is removed. Days must be 1 or more,
// default 7.
//
// Requires KICK_MEMBERS.
func (c *Client) PruneWithCount(guildID discord.Snowflake, days uint) (uint, error) {
	if days == 0 {
		days = 7
	}

	var param struct {
		Days     uint `schema:"days"`
		RetCount bool `schema:"compute_prune_count"`
	}

	param.Days = days
	param.RetCount = true

	var resp struct {
		Pruned uint `json:"pruned"`
	}

	return resp.Pruned, c.RequestJSON(
		&resp, "POST",
		EndpointGuilds+guildID.String()+"/prune",
		httputil.WithSchema(c, param),
	)
}

// Kick removes a member from a guild.
//
// Requires KICK_MEMBERS permission.
// Fires a Guild Member Remove Gateway event.
func (c *Client) Kick(guildID, userID discord.Snowflake) error {
	return c.FastRequest(
		"DELETE",
		EndpointGuilds+guildID.String()+"/members/"+userID.String(),
	)
}

// Bans returns a list of ban objects for the users banned from this guild.
//
// Requires the BAN_MEMBERS permission.
func (c *Client) Bans(guildID discord.Snowflake) ([]discord.Ban, error) {
	var bans []discord.Ban
	return bans, c.RequestJSON(
		&bans, "GET",
		EndpointGuilds+guildID.String()+"/bans",
	)
}

// GetBan returns a ban object for the given user.
//
// Requires the BAN_MEMBERS permission.
func (c *Client) GetBan(guildID, userID discord.Snowflake) (*discord.Ban, error) {
	var ban *discord.Ban
	return ban, c.RequestJSON(
		&ban, "GET",
		EndpointGuilds+guildID.String()+"/bans/"+userID.String(),
	)
}

// https://discord.com/developers/docs/resources/guild#create-guild-ban-query-string-params
type BanData struct {
	// DeleteDays is the number of days to delete messages for (0-7).
	DeleteDays option.Uint `schema:"delete_message_days,omitempty"`
	// Reason is the reason for the ban.
	Reason option.String `schema:"reason,omitempty"`
}

// Ban creates a guild ban, and optionally delete previous messages sent by the
// banned user.
//
// Requires the BAN_MEMBERS permission.
func (c *Client) Ban(guildID, userID discord.Snowflake, data BanData) error {
	if *data.DeleteDays > 7 {
		*data.DeleteDays = 7
	}

	return c.FastRequest(
		"PUT",
		EndpointGuilds+guildID.String()+"/bans/"+userID.String(),
		httputil.WithSchema(c, data),
	)
}

// Unban removes the ban for a user.
//
// Requires the BAN_MEMBERS permissions.
// Fires a Guild Ban Remove Gateway event.
func (c *Client) Unban(guildID, userID discord.Snowflake) error {
	return c.FastRequest("DELETE", EndpointGuilds+guildID.String()+"/bans/"+userID.String())
}
