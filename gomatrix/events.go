package gomatrix

import (
	"html"
	"regexp"
	"time"
)

// Event represents a single Matrix event.
type Event struct {
	StateKey     *string                `json:"state_key,omitempty"` // The state key for the event. Only present on State Events.
	Sender       string                 `json:"sender"`              // The user ID of the sender of the event
	Type         string                 `json:"type"`                // The event type
	Timestamp    int64                  `json:"origin_server_ts"`    // The unix timestamp when this message was sent by the origin server
	Time         time.Time              `json:"time"`
	Date         string                 `json:"date"`
	When         string                 `json:"when"`
	ID           string                 `json:"event_id"`               // The unique ID of this event
	RoomID       string                 `json:"room_id"`                // The room the event was sent to. May be nil (e.g. for presence)
	Redacts      string                 `json:"redacts,omitempty"`      // The event ID that was redacted if a m.room.redaction event
	Unsigned     map[string]interface{} `json:"unsigned"`               // The unsigned portions of the event, such as age and prev_content
	Content      map[string]interface{} `json:"content"`                // The JSON content of the event.
	PrevContent  map[string]interface{} `json:"prev_content,omitempty"` // The JSON prev_content of the event.
	Author       Author                 `json:"author,omitempty"`
	SharedPost   *Event                 `json:"shared_post,omitempty"`
	Redacted     bool                   `json:"redacted"`
	Owner        bool                   `json:"owner"`
	ShortID      string                 `json:"short_id"`
	IsArticle    bool                   `json:"is_article"`
	Replies      []*Event               `json:"replies"`
	TotalReplies int                    `json:"total_replies"`
	Reactions    *[]*Reaction           `json:"reactions"`
}

type Reaction struct {
	Key    string `json:"key"`
	Author Author `json:"author"`
}

type Author struct {
	DisplayName string `json:"display_name,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	FormattedID string `json:"formatted_id,omitempty"`
}

// Body returns the value of the "body" key in the event content if it is
// present and is a string.
func (event *Event) Body() (body string, ok bool) {
	value, exists := event.Content["body"]
	if !exists {
		return
	}
	body, ok = value.(string)
	return
}

// MessageType returns the value of the "msgtype" key in the event content if
// it is present and is a string.
func (event *Event) MessageType() (msgtype string, ok bool) {
	value, exists := event.Content["msgtype"]
	if !exists {
		return
	}
	msgtype, ok = value.(string)
	return
}

// This is the standard Hummngbard post event.
type Post struct {
	MsgType         string            `json:"msgtype"`
	Body            string            `json:"body"`
	FormattedBody   string            `json:"formatted_body"`
	Format          string            `json:"format"`
	Reply           bool              `json:"reply"`
	Links           []Link            `json:"links,omitempty"`
	Attachments     []Attachment      `json:"attachments,omitempty"`
	Images          []Image           `json:"images,omitempty"`
	ThreadRoomID    *string           `json:"thread_room_id,omitempty"`
	ThreadRoomAlias *string           `json:"thread_room_alias"`
	ThreadInRoomID  *string           `json:"thread_in_room_id,omitempty"`
	EventID         *string           `json:"event_id,omitempty"`
	Root            *bool             `json:"root,omitempty"`
	RoomAlias       *string           `json:"room_alias,omitempty"`
	RoomPath        string            `json:"room_path,omitempty"`
	MRelationship   map[string]string `json:"m.relationship,omitempty"`
	MRelatesTo      map[string]string `json:"m_relates_to"`
	MNewContent     map[string]string `json:"m_new_content"`
	NSFW            bool              `json:"nsfw,omitempty"`
	Anonymous       bool              `json:"anonymous,omitempty"`
	SharedPost      interface{}       `json:"shared_post,omitempty"`
	Article         *Article          `json:"com.hummingbard.article,omitempty"`
	ShareReply      *bool             `json:"share_reply"`
	ReplyPermalink  *string           `json:"reply_permalink"`
}

type Article struct {
	Title         string `json:"title,omitempty"`
	Slug          string `json:"slug,omitempty"`
	Subtitle      string `json:"subtitle,omitempty"`
	Description   string `json:"description,omitempty"`
	CanonicalLink string `json:"canonical_link,omitempty"`
	FeaturedImage *Image `json:"featured_image,omitempty"`
	ContentURI    string `json:"content_uri"`
}

type ReqReaction struct {
	MRelatesTo    map[string]string `json:"m_relates_to"`
	MRelationship map[string]string `json:"m.relationship,omitempty"`
}

type Link struct {
	Href        string `json:"href"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author,omitempty"`
	Image       string `json:"image,omitempty"`
	IsYoutube   bool   `json:"is_youtube"`
	YoutubeID   string `json:"youtube_id"`
}

type Attachment struct {
	Filename string `json:"filename"`
	Mimetype string `json:"mimetype"`
	Size     uint   `json:"size"`
	MXC      string `json:"mxc"`
}

type Image struct {
	Caption     string `json:"caption,omitempty"`
	Description string `json:"description,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Height      uint   `json:"height,omitempty"`
	Width       uint   `json:"width,omitempty"`
	AspectRatio uint   `json:"aspect_ratio,omitempty"`
	Mimetype    string `json:"mimetype,omitempty"`
	Size        uint   `json:"size,omitempty"`
	MXC         string `json:"mxc,omitempty"`
}

// TextMessage is the contents of a Matrix formated message event.
type TextMessage struct {
	MsgType       string `json:"msgtype"`
	Body          string `json:"body"`
	FormattedBody string `json:"formatted_body"`
	Format        string `json:"format"`
}

// ThumbnailInfo contains info about an thumbnail image - http://matrix.org/docs/spec/client_server/r0.2.0.html#m-image
type ThumbnailInfo struct {
	Height   uint   `json:"h,omitempty"`
	Width    uint   `json:"w,omitempty"`
	Mimetype string `json:"mimetype,omitempty"`
	Size     uint   `json:"size,omitempty"`
}

// ImageInfo contains info about an image - http://matrix.org/docs/spec/client_server/r0.2.0.html#m-image
type ImageInfo struct {
	Height        uint          `json:"h,omitempty"`
	Width         uint          `json:"w,omitempty"`
	Mimetype      string        `json:"mimetype,omitempty"`
	Size          uint          `json:"size,omitempty"`
	ThumbnailInfo ThumbnailInfo `json:"thumbnail_info,omitempty"`
	ThumbnailURL  string        `json:"thumbnail_url,omitempty"`
}

// VideoInfo contains info about a video - http://matrix.org/docs/spec/client_server/r0.2.0.html#m-video
type VideoInfo struct {
	Mimetype      string        `json:"mimetype,omitempty"`
	ThumbnailInfo ThumbnailInfo `json:"thumbnail_info"`
	ThumbnailURL  string        `json:"thumbnail_url,omitempty"`
	Height        uint          `json:"h,omitempty"`
	Width         uint          `json:"w,omitempty"`
	Duration      uint          `json:"duration,omitempty"`
	Size          uint          `json:"size,omitempty"`
}

// VideoMessage is an m.video  - http://matrix.org/docs/spec/client_server/r0.2.0.html#m-video
type VideoMessage struct {
	MsgType string    `json:"msgtype"`
	Body    string    `json:"body"`
	URL     string    `json:"url"`
	Info    VideoInfo `json:"info"`
}

// ImageMessage is an m.image event
type ImageMessage struct {
	MsgType string    `json:"msgtype"`
	Body    string    `json:"body"`
	URL     string    `json:"url"`
	Info    ImageInfo `json:"info"`
}

// An HTMLMessage is the contents of a Matrix HTML formated message event.
type HTMLMessage struct {
	Body          string `json:"body"`
	MsgType       string `json:"msgtype"`
	Format        string `json:"format"`
	FormattedBody string `json:"formatted_body"`
}

// FileInfo contains info about an file - http://matrix.org/docs/spec/client_server/r0.2.0.html#m-file
type FileInfo struct {
	Mimetype string `json:"mimetype,omitempty"`
	Size     uint   `json:"size,omitempty"` //filesize in bytes
}

// FileMessage is an m.file event - http://matrix.org/docs/spec/client_server/r0.2.0.html#m-file
type FileMessage struct {
	MsgType       string    `json:"msgtype"`
	Body          string    `json:"body"`
	URL           string    `json:"url"`
	Filename      string    `json:"filename"`
	Info          FileInfo  `json:"info,omitempty"`
	ThumbnailURL  string    `json:"thumbnail_url,omitempty"`
	ThumbnailInfo ImageInfo `json:"thumbnail_info,omitempty"`
}

// LocationMessage is an m.location event - http://matrix.org/docs/spec/client_server/r0.2.0.html#m-location
type LocationMessage struct {
	MsgType       string    `json:"msgtype"`
	Body          string    `json:"body"`
	GeoURI        string    `json:"geo_uri"`
	ThumbnailURL  string    `json:"thumbnail_url,omitempty"`
	ThumbnailInfo ImageInfo `json:"thumbnail_info,omitempty"`
}

// AudioInfo contains info about an file - http://matrix.org/docs/spec/client_server/r0.2.0.html#m-audio
type AudioInfo struct {
	Mimetype string `json:"mimetype,omitempty"`
	Size     uint   `json:"size,omitempty"`     //filesize in bytes
	Duration uint   `json:"duration,omitempty"` //audio duration in ms
}

// AudioMessage is an m.audio event - http://matrix.org/docs/spec/client_server/r0.2.0.html#m-audio
type AudioMessage struct {
	MsgType string    `json:"msgtype"`
	Body    string    `json:"body"`
	URL     string    `json:"url"`
	Info    AudioInfo `json:"info,omitempty"`
}

var htmlRegex = regexp.MustCompile("<[^<]+?>")

// GetHTMLMessage returns an HTMLMessage with the body set to a stripped version of the provided HTML, in addition
// to the provided HTML.
func GetHTMLMessage(msgtype, htmlText string) HTMLMessage {
	return HTMLMessage{
		Body:          html.UnescapeString(htmlRegex.ReplaceAllLiteralString(htmlText, "")),
		MsgType:       msgtype,
		Format:        "org.matrix.custom.html",
		FormattedBody: htmlText,
	}
}
