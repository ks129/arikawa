package api

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/utils/httputil"
	"github.com/diamondburned/arikawa/utils/json/option"
)

var EndpointChannels = Endpoint + "channels/"

// Channels returns a list of guild channel objects.
func (c *Client) Channels(guildID discord.Snowflake) ([]discord.Channel, error) {
	var chs []discord.Channel
	return chs, c.RequestJSON(&chs, "GET", EndpointGuilds+guildID.String()+"/channels")
}

// https://discord.com/developers/docs/resources/guild#create-guild-channel-json-params
type CreateChannelData struct {
	// Name is the channel name (2-100 characters).
	//
	// Channel Type: All
	Name string `json:"name"`
	// Type is the type of channel.
	//
	// Channel Type: All
	Type discord.ChannelType `json:"type,omitempty"`
	// Topic is the channel topic (0-1024 characters).
	//
	// Channel Types: Text, News
	Topic string `json:"topic,omitempty"`
	// VoiceBitrate is the bitrate (in bits) of the voice channel.
	// 8000 to 96000 (128000 for VIP servers)
	//
	// Channel Types: Voice
	VoiceBitrate uint `json:"bitrate,omitempty"`
	// VoiceUserLimit is the user limit of the voice channel.
	// 0 refers to no limit, 1 to 99 refers to a user limit.
	//
	// Channel Types: Voice
	VoiceUserLimit uint `json:"user_limit,omitempty"`
	// UserRateLimit is the amount of seconds a user has to wait before sending
	// another message (0-21600).
	// Bots, as well as users with the permission manage_messages or
	// manage_channel, are unaffected.
	//
	// Channel Types: Text
	UserRateLimit discord.Seconds `json:"rate_limit_per_user,omitempty"`
	// Position is the sorting position of the channel.
	//
	// Channel Types: All
	Position option.Int `json:"position,omitempty"`
	// Permissions are the channel's permission overwrites.
	//
	// Channel Types: All
	Permissions []discord.Overwrite `json:"permission_overwrites,omitempty"`
	// CategoryID is the 	id of the parent category for a channel.
	//
	// Channel Types: Text, News, Store, Voice
	CategoryID discord.Snowflake `json:"parent_id,string,omitempty"`
	// NSFW specifies whether the channel is nsfw.
	//
	// Channel Types: Text, News, Store.
	NSFW bool `json:"nsfw,omitempty"`
}

// CreateChannel creates a new channel object for the guild.
//
// Requires the MANAGE_CHANNELS permission.
// Fires a Channel Create Gateway event.
func (c *Client) CreateChannel(
	guildID discord.Snowflake, data CreateChannelData) (*discord.Channel, error) {
	var ch *discord.Channel
	return ch, c.RequestJSON(
		&ch, "POST",
		EndpointGuilds+guildID.String()+"/channels",
		httputil.WithJSONBody(data),
	)
}

type MoveChannelData struct {
	// ID is the channel id.
	ID discord.Snowflake `json:"id"`
	// Position is the sorting position of the channel
	Position option.Int `json:"position"`
}

// MoveChannel modifies the position of channels in the guild.
//
// Requires MANAGE_CHANNELS.
func (c *Client) MoveChannel(guildID discord.Snowflake, datum []MoveChannelData) error {
	return c.FastRequest(
		"PATCH",
		EndpointGuilds+guildID.String()+"/channels", httputil.WithJSONBody(datum),
	)
}

// Channel gets a channel by ID. Returns a channel object.
func (c *Client) Channel(channelID discord.Snowflake) (*discord.Channel, error) {
	var channel *discord.Channel
	return channel, c.RequestJSON(&channel, "GET", EndpointChannels+channelID.String())
}

// https://discord.com/developers/docs/resources/channel#modify-channel-json-params
type ModifyChannelData struct {
	// Name is the 2-100 character channel name.
	//
	// Channel Types: All
	Name string `json:"name,omitempty"`
	// Type is the type of the channel.
	// Only conversion between text and news is supported and only in guilds
	// with the "NEWS" feature
	//
	// Channel Types: Text, News
	Type *discord.ChannelType `json:"type,omitempty"`
	// Postion is the position of the channel in the left-hand listing
	//
	// Channel Types: All
	Position option.NullableInt `json:"position,omitempty"`
	// Topic is the 0-1024 character channel topic.
	//
	// Channel Types: Text, News
	Topic option.NullableString `json:"topic,omitempty"`
	// NSFW specifies whether the channel is nsfw.
	//
	// Channel Types: Text, News, Store.
	NSFW option.NullableBool `json:"nsfw,omitempty"`
	// UserRateLimit is the amount of seconds a user has to wait before sending
	// another message (0-21600).
	// Bots, as well as users with the permission manage_messages or
	// manage_channel, are unaffected.
	//
	// Channel Types: Text
	UserRateLimit option.NullableUint `json:"rate_limit_per_user,omitempty"`
	// VoiceBitrate is the bitrate (in bits) of the voice channel.
	// 8000 to 96000 (128000 for VIP servers)
	//
	// Channel Types: Voice
	VoiceBitrate option.NullableUint `json:"bitrate,omitempty"`
	// VoiceUserLimit is the user limit of the voice channel.
	// 0 refers to no limit, 1 to 99 refers to a user limit.
	//
	// Channel Types: Voice
	VoiceUserLimit option.NullableUint `json:"user_limit,omitempty"`
	// Permissions are the channel or category-specific permissions.
	//
	// Channel Types: All
	Permissions *[]discord.Overwrite `json:"permission_overwrites,omitempty"`
	// CategoryID is the id of the new parent category for a channel.
	// Channel Types: Text, News, Store, Voice
	CategoryID discord.Snowflake `json:"parent_id,string,omitempty"`
}

// ModifyChannel updates a channel's settings.
//
// Requires the MANAGE_CHANNELS permission for the guild.
func (c *Client) ModifyChannel(channelID discord.Snowflake, data ModifyChannelData) error {
	return c.FastRequest("PATCH", EndpointChannels+channelID.String(), httputil.WithJSONBody(data))
}

// DeleteChannel deletes a channel, or closes a private message. Requires the
// MANAGE_CHANNELS permission for the guild. Deleting a category does not
// delete its child channels: they will have their parent_id removed and a
// Channel Update Gateway event will fire for each of them.
//
// Fires a Channel Delete Gateway event.
func (c *Client) DeleteChannel(channelID discord.Snowflake) error {
	return c.FastRequest("DELETE", EndpointChannels+channelID.String())
}

// EditChannelPermission edits the channel's permission overwrites for a user
// or role in a channel. Only usable for guild channels.
//
// Requires the MANAGE_ROLES permission.
func (c *Client) EditChannelPermission(
	channelID discord.Snowflake, overwrite discord.Overwrite) error {

	url := EndpointChannels + channelID.String() + "/permissions/" + overwrite.ID.String()
	overwrite.ID = 0

	return c.FastRequest("PUT", url, httputil.WithJSONBody(overwrite))
}

// DeleteChannelPermission deletes a channel permission overwrite for a user or
// role in a channel. Only usable for guild channels.
//
// Requires the MANAGE_ROLES permission.
func (c *Client) DeleteChannelPermission(channelID, overwriteID discord.Snowflake) error {
	return c.FastRequest(
		"DELETE",
		EndpointChannels+channelID.String()+"/permissions/"+overwriteID.String(),
	)
}

// Typing posts a typing indicator to the channel. Undocumented, but the client
// usually clears the typing indicator after 8-10 seconds (or after a message).
func (c *Client) Typing(channelID discord.Snowflake) error {
	return c.FastRequest("POST", EndpointChannels+channelID.String()+"/typing")
}

// PinnedMessages returns all pinned messages in the channel as an array of
// message objects.
func (c *Client) PinnedMessages(channelID discord.Snowflake) ([]discord.Message, error) {
	var pinned []discord.Message
	return pinned, c.RequestJSON(&pinned, "GET", EndpointChannels+channelID.String()+"/pins")
}

// PinMessage pins a message in a channel.
//
// Requires the MANAGE_MESSAGES permission.
func (c *Client) PinMessage(channelID, messageID discord.Snowflake) error {
	return c.FastRequest("PUT", EndpointChannels+channelID.String()+"/pins/"+messageID.String())
}

// UnpinMessage deletes a pinned message in a channel.
//
// Requires the MANAGE_MESSAGES permission.
func (c *Client) UnpinMessage(channelID, messageID discord.Snowflake) error {
	return c.FastRequest("DELETE", EndpointChannels+channelID.String()+"/pins/"+messageID.String())
}

// AddRecipient adds a user to a group direct message. As accessToken is needed,
// clearly this endpoint should only be used for OAuth. AccessToken can be
// obtained with the "gdm.join" scope.
func (c *Client) AddRecipient(
	channelID, userID discord.Snowflake, accessToken, nickname string) error {

	var params struct {
		AccessToken string `json:"access_token"`
		Nickname    string `json:"nickname"`
	}

	params.AccessToken = accessToken
	params.Nickname = nickname

	return c.FastRequest(
		"PUT",
		EndpointChannels+channelID.String()+"/recipients/"+userID.String(),
		httputil.WithJSONBody(params),
	)
}

// RemoveRecipient removes a user from a group direct message.
func (c *Client) RemoveRecipient(channelID, userID discord.Snowflake) error {
	return c.FastRequest(
		"DELETE",
		EndpointChannels+channelID.String()+"/recipients/"+userID.String(),
	)
}

// Ack is the read state of a channel. This is undocumented.
type Ack struct {
	Token string `json:"token"`
}

// Ack marks the read state of a channel. This is undocumented. The method will
// write to the ack variable passed in. If this method is called asynchronously,
// then ack should be mutex guarded.
func (c *Client) Ack(channelID, messageID discord.Snowflake, ack *Ack) error {
	return c.RequestJSON(
		ack, "POST",
		EndpointChannels+channelID.String()+"/messages/"+messageID.String()+"/ack",
		httputil.WithJSONBody(ack),
	)
}
