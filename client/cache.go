package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"hummingbard/gomatrix"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type JoinedPublicRoom struct {
	RoomID string      `json:"room_id"`
	Events interface{} `json:"events"`
}

func (c *Client) GetAllPublicRooms() ([]*JoinedPublicRoom, error) {

	fil, err := c.Matrix.CreateFilter([]byte(`
{
	"room": {
		"timeline": {
			"limit": 0,
			"types": ["com.hummingbard.post"]
		}
	}
}
	`))

	sre, err := c.Matrix.SyncRequest(0, "", fil.FilterID, true, "offline")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rms := []*JoinedPublicRoom{}
	for roomID, room := range sre.Rooms.Join {

		st, _ := json.Marshal(room.State.Events)
		roomType := gjson.Parse(string(st)).Get(`#(type="com.hummingbard.room")`).Get("content.room_type")

		if roomType.String() == "page" || roomType.String() == "post" {
			continue
		}
		alias := gjson.Parse(string(st)).Get(`#(type="m.room.canonical_alias")`).Get("content.alias")
		if len(alias.String()) > 0 &&
			!strings.Contains(alias.String(), "@") &&
			!strings.Contains(alias.String(), "#thread") &&
			!strings.Contains(alias.String(), "#public") {
			rms = append(rms, &JoinedPublicRoom{RoomID: roomID, Events: room.State.Events})
		}
	}
	return rms, nil
}

type PublicRoom struct {
	RoomID    string        `json:"room_id"`
	RoomPath  string        `json:"room_path"`
	RoomAlias string        `json:"room_alias"`
	Avatar    string        `json:"avatar"`
	Name      string        `json:"name"`
	Topic     template.HTML `json:"topic"`
}

func (c *Client) GetPublicRooms() ([]*PublicRoom, error) {

	fil, err := c.Matrix.CreateFilter([]byte(`
{
	"room": {
		"timeline": {
			"limit": 0,
			"types": ["com.hummingbard.post"]
		}
	}
}
	`))

	sre, err := c.Matrix.SyncRequest(0, "", fil.FilterID, true, "offline")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rms := []*PublicRoom{}
	for roomID, room := range sre.Rooms.Join {

		st, _ := json.Marshal(room.State.Events)
		roomType := gjson.Parse(string(st)).Get(`#(type="com.hummingbard.room")`).Get("content.room_type")

		if roomType.String() == "page" || roomType.String() == "post" {
			continue
		}

		alias := gjson.Parse(string(st)).Get(`#(type="m.room.canonical_alias")`).Get("content.alias")
		if len(alias.String()) > 0 &&
			!strings.Contains(alias.String(), "@") &&
			!strings.Contains(alias.String(), "#thread") &&
			!strings.Contains(alias.String(), "#public") {

			room := &PublicRoom{
				RoomID:    roomID,
				RoomAlias: alias.String(),
			}

			sp := strings.Split(alias.String(), ":")
			rp := sp[0][1:]
			p := strings.Split(rp, "_")
			al := strings.Join(p, "/")

			if strings.Contains(alias.String(), fmt.Sprintf(`:%s`, c.Config.Client.Domain)) {
				room.RoomPath = al
			} else {
				room.RoomPath = fmt.Sprintf(`%s:%s`, al, sp[1])
			}

			avatar := gjson.Parse(string(st)).Get(`#(type="m.room.avatar")`).Get("content.url")
			if avatar.String() != "" {
				room.Avatar = c.BuildAvatar(avatar.String())
			}

			name := gjson.Parse(string(st)).Get(`#(type="m.room.name")`).Get("content.name")
			if name.String() != "" {
				room.Name = name.String()
			}

			topic := gjson.Parse(string(st)).Get(`#(type="m.room.topic")`).Get("content.topic")
			if topic.String() != "" {
				html, err := ToStrictHTML(topic.String())
				if err != nil {
					log.Println(err)
					room.Topic = template.HTML(topic.String())
				} else {
					room.Topic = html
				}
			}

			rms = append(rms, room)
		}
	}
	return rms, nil
}

func (c *Client) GetPublicRoomsFromCache() ([]*PublicRoom, error) {
	rooms, ok := c.Cache.PublicRooms.Get("public")
	if !ok {
		return nil, errors.New("couldn't find public room in cache")
	}

	if v, ok := rooms.([]*PublicRoom); ok {
		if len(v) > 0 {
			return v, nil
		}
	}

	return nil, nil
}

func (c *Client) RefreshCache() {
	err := c.RefreshRoomsCache()
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) RefreshRoomsCache() error {

	items, err := c.GetPublicRooms()
	if err != nil {
		log.Println(err)
		return err
	}
	sort.Slice(items, func(i, j int) bool { return items[i].RoomAlias < items[j].RoomAlias })
	c.Cache.PublicRooms.Set("public", items, 1)

	return nil
}

func (c *Client) GetRoomAliasFromCache(roomID string) (string, error) {

	room, ok := c.Cache.Rooms.Get(roomID)
	if !ok {
		return "", errors.New("couldn't find room in cache")
	}

	return room.(gomatrix.PublicRoom).CanonicalAlias, nil
}

func (c *Client) ProcessJoinedRooms(rooms []string) ([]JoinedRoom, error) {

	joinedRooms := []JoinedRoom{}
	for _, x := range rooms {
		r := JoinedRoom{
			RoomID: x,
		}
		room, ok := c.Cache.Rooms.Get(x)
		if ok {
			r.RoomAlias = room.(gomatrix.PublicRoom).CanonicalAlias
		}
		joinedRooms = append(joinedRooms, r)
	}

	return joinedRooms, nil
}

func (c *Client) RefreshRoomEvents(roomID string) (bool, error) {
	token := c.DefaultUser.AccessToken
	userid := c.DefaultUser.UserID

	serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

	cli, err := gomatrix.NewClient(serverName, userid, token)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	state, err := cli.RoomState(roomID)
	if err != nil {
		log.Println(err)
		return false, err
	}
	cli.Prefix = "/_matrix/client/"
	rc := c.RoomCreateEventFromState(state)

	opts := map[string]interface{}{
		"event_id":         rc,
		"room_id":          roomID,
		"depth_first":      false,
		"recent_first":     true,
		"include_parent":   false,
		"include_children": true,
		"direction":        "down",
		"limit":            14,
		"max_depth":        0,
		"max_breadth":      0,
		"last_event":       "0",
	}

	relationships, err := cli.GetRelationships(opts)
	if err != nil {
		log.Println(err)
	}

	cli.Prefix = "/_matrix/client/r0"

	cachedEvents := CachedRoomEvents{
		Events: relationships.Events,
		Time:   time.Now(),
	}

	c.Cache.Events.Set(roomID, cachedEvents, 1)

	return true, nil
}

func (c *Client) RefreshPublicEvents(roomID string) (bool, error) {
	serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

	cli, err := gomatrix.NewClient(serverName, c.DefaultUser.UserID, c.DefaultUser.AccessToken)
	if err != nil {
		log.Println(err)
		return false, err
	}

	var pev PublicEvents

	msg, err := cli.Messages(roomID, "", "", 'b', 23, "")
	if err != nil {
		log.Println(err)
	}
	pev = PublicEvents{
		Events:    msg.Chunk,
		LastEvent: msg.End,
	}
	c.Cache.Events.Set(roomID, pev, 1)

	return true, nil
}
