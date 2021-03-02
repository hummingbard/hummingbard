package client

import (
	"encoding/json"
	"fmt"
	"html/template"
	"hummingbard/gomatrix"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/unrolled/secure"
)

type PublicEvents struct {
	Events    []gomatrix.Event
	LastEvent string
}

func (c *Client) PublicFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		us := LoggedInUser(r)

		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		cli, err := gomatrix.NewClient(serverName, c.DefaultUser.UserID, c.DefaultUser.AccessToken)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		room := fmt.Sprintf(`#public:%s`, c.Config.Client.Domain)

		ra, err := cli.ResolveAlias(room)
		if err != nil {
			log.Println(err)
			c.NotFound(w, r)
			return
		}

		roomID := string(ra.RoomID)

		state, err := cli.RoomState(roomID)
		if err != nil {
			log.Println(err)
		}

		t := TimelinePage{
			Room: Room{
				Path: "public",
				ID:   string(ra.RoomID),
			},
			RoomState:     state,
			IsUserProfile: false,
		}

		var pev PublicEvents

		cachedEvents, ok := c.Cache.Events.Get(roomID)
		if ok {
			if x, ok := cachedEvents.(PublicEvents); ok {
				pev = x
				log.Println(len(pev.Events))
			}
		} else {
			msg, err := cli.Messages(roomID, "", "", 'b', 23, "")
			if err != nil {
				log.Println(err)
			}
			pev = PublicEvents{
				Events:    msg.Chunk,
				LastEvent: msg.End,
			}
			c.Cache.Events.Set(roomID, pev, 1)
		}

		t.Posts = c.ProcessMessages(pev.Events, state, us)
		t.LastEvent = pev.LastEvent

		nonce := secure.CSPNonce(r.Context())

		t.Nonce = nonce
		t.LoggedInUser = us

		ip, err := json.Marshal(t.Posts)
		if err != nil {
			c.Error(w, r)
			return
		}

		t.InitialPosts = string(ip)

		t.Room.Type = "feed"
		t.Room.Alias = "public"

		if c.Config.Mode == "development" {
			t.HomeServerURL = c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)
		} else {
			t.HomeServerURL = fmt.Sprintf(`https://%s`, c.Config.Matrix.Server)
		}

		c.Templates.ExecuteTemplate(w, "public", t)
	}
}

func (c *Client) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		us := LoggedInUser(r)

		if us != nil {
			c.IndexUser(w, r)
			return
		}

		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		cli, err := gomatrix.NewClient(serverName, c.DefaultUser.UserID, c.DefaultUser.AccessToken)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		room := fmt.Sprintf(`#public:%s`, c.Config.Client.Domain)

		ra, err := cli.ResolveAlias(room)
		if err != nil {
			log.Println(err)
			c.NotFound(w, r)
			return
		}

		roomID := string(ra.RoomID)

		state, err := cli.RoomState(roomID)
		if err != nil {
			log.Println(err)
		}

		t := TimelinePage{
			Room: Room{
				Path: "public",
				ID:   string(ra.RoomID),
			},
			RoomState:     state,
			IsUserProfile: false,
		}

		var pev PublicEvents

		cachedEvents, ok := c.Cache.Events.Get(roomID)
		if ok {
			if x, ok := cachedEvents.(PublicEvents); ok {
				pev = x
			}
		} else {
			msg, err := cli.Messages(roomID, "", "", 'b', 23, "")
			if err != nil {
				log.Println(err)
			}
			pev = PublicEvents{
				Events:    msg.Chunk,
				LastEvent: msg.End,
			}
			c.Cache.Events.Set(roomID, pev, 1)
		}

		t.Posts = c.ProcessMessages(pev.Events, state, us)
		t.LastEvent = pev.LastEvent

		rooms, err := c.GetPublicRoomsFromCache()
		if err != nil {
			log.Println(err)
		}
		if rooms != nil && len(rooms) > 0 {
			r := []*PublicRoom{}
			for _, x := range rooms {
				/*
					if i == 53 {
						break
					}
				*/
				if !strings.Contains(x.RoomPath, "/") {
					r = append(r, x)
				}
			}
			t.PublicRooms = r
		}

		//let's add them manually for now

		/*
			rooms := []string{
				"art",
				"design",
				"environment",
				"eargasm:matrix.org",
				"hummingbard",
				"hummingbard/bugs",
				"hummingbard/feature-requests",
				"music",
				"music/classical",
				"music/jazz",
				"music/jazz/fusion",
				"music/metal",
				"programming",
				"rhythm:matrix.org",
				"news",
				"technology",
				"videos",
			}
		*/

		nonce := secure.CSPNonce(r.Context())

		t.Nonce = nonce

		ip, err := json.Marshal(t.Posts)
		if err != nil {
			c.Error(w, r)
			return
		}

		t.InitialPosts = string(ip)

		t.Room.Type = "feed"
		t.Room.Alias = "public"

		if c.Config.Mode == "development" {
			t.HomeServerURL = c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)
		} else {
			t.HomeServerURL = fmt.Sprintf(`https://%s`, c.Config.Matrix.Server)
		}

		c.Templates.ExecuteTemplate(w, "index", t)
	}
}

func (c *Client) IndexUser(w http.ResponseWriter, r *http.Request) {
	user := LoggedInUser(r)

	s, err := GetSession(r, c)
	if err != nil {
		log.Println(err)
	}
	if s != nil {
		x := s.Flashes("signed-up")
		if len(x) > 0 {
			s.Save(r, w)
			c.WelcomePage(w, r)
			return
		}
	}

	serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

	matrix, err := gomatrix.NewClient(serverName, c.DefaultUser.UserID, c.DefaultUser.AccessToken)
	if err != nil {
		log.Println(err)
		c.Error(w, r)
		return
	}

	posts := []gomatrix.Event{}
	feedItems := []*IndexFeed{}

	limit := 14

	l := len(user.JoinedRooms)

	if len(user.JoinedRooms) > 0 {

		switch {
		case l > 14:
			limit = 2
		case l > 9 && l < 14:
			limit = 3
		case l > 5 && l <= 9:
			limit = 3
		case l > 2 && l <= 4:
			limit = 6
		case l == 2:
			limit = 8
		}

		for _, room := range user.JoinedRooms {
			matrix.Prefix = "/_matrix/client/r0"

			state, err := matrix.RoomState(room.RoomID)
			if err != nil {
				log.Println(err)
				continue
			}
			matrix.Prefix = "/_matrix/client/"

			rc := c.RoomCreateEventFromState(state)

			opts := map[string]interface{}{
				"event_id":         rc,
				"room_id":          room.RoomID,
				"depth_first":      false,
				"recent_first":     false,
				"include_parent":   false,
				"include_children": true,
				"direction":        "down",
				"limit":            limit,
				"max_depth":        0,
				"max_breadth":      0,
				"last_event":       "0",
			}

			relationships, err := matrix.GetRelationships(opts)
			if err != nil {
				log.Println(err)
				continue
			}
			if relationships != nil && len(relationships.Events) > 0 {

				processed := c.ProcessMessages(relationships.Events, state, user)
				posts = append(posts, processed...)
				lastEvent := relationships.Events[len(relationships.Events)-1].Timestamp
				feedItems = append(feedItems, &IndexFeed{
					RoomID:    room.RoomID,
					LastEvent: lastEvent,
				})
			}
		}
		sort.Slice(posts, func(i, j int) bool { return posts[j].Time.Before(posts[i].Time) })
	}

	type page struct {
		Room Room
		BasePage
		LoggedInUser  *User
		State         interface{}
		ProfileLink   template.URL
		Posts         interface{}
		HomeServerURL string
		FeedItems     interface{}
		InitialPosts  interface{}
		Depth         int
	}

	nonce := secure.CSPNonce(r.Context())

	t := &page{
		Room: Room{
			Path: "user-index",
		},
		Posts:     posts,
		FeedItems: feedItems,
	}

	if user.Federated {
		t.ProfileLink = template.URL(user.UserID)
	} else {
		sp := strings.Split(user.UserID, ":")
		t.ProfileLink = template.URL(sp[0])
	}

	if c.Config.Mode == "development" {
		t.HomeServerURL = c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)
	} else {
		t.HomeServerURL = fmt.Sprintf(`https://%s`, c.Config.Matrix.Server)
	}

	t.Nonce = nonce
	ip, err := json.Marshal(t.Posts)
	if err != nil {
		c.Error(w, r)
		return
	}

	t.InitialPosts = string(ip)
	t.LoggedInUser = user

	c.Templates.ExecuteTemplate(w, "index-user", t)
}

func (c *Client) WelcomePage(w http.ResponseWriter, r *http.Request) {
	user := LoggedInUser(r)

	type page struct {
		BasePage
		LoggedInUser  *User
		HomeServerURL string
		Rooms         interface{}
	}

	nonce := secure.CSPNonce(r.Context())

	t := &page{}
	rooms, err := c.GetPublicRoomsFromCache()
	if err != nil {
		log.Println(err)
	}
	if rooms != nil && len(rooms) > 0 {
		r := []*PublicRoom{}
		for i, x := range rooms {
			if i == 53 {
				break
			}
			if !(strings.Contains(x.RoomPath, "/") &&
				strings.Contains(x.RoomPath, ":")) {
				r = append(r, x)
			}
		}
		t.Rooms = r
	}

	t.Nonce = nonce
	t.LoggedInUser = user

	if c.Config.Mode == "development" {
		t.HomeServerURL = c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)
	} else {
		t.HomeServerURL = fmt.Sprintf(`https://%s`, c.Config.Matrix.Server)
	}

	c.Templates.ExecuteTemplate(w, "welcome", t)
}

func (c *Client) GetFeedEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
		}

		type payload struct {
			Feed []struct {
				RoomID    string `json:"room_id"`
				LastEvent int64  `json:"last_event"`
			} `json:"feed"`
		}

		var pay payload
		if r.Body == nil {
			log.Println(err)
			http.Error(w, "Please send a request body", 400)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		matrix, err := gomatrix.NewClient(serverName, c.DefaultUser.UserID, c.DefaultUser.AccessToken)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		posts := []gomatrix.Event{}
		feedItems := []*IndexFeed{}

		limit := 14

		l := len(pay.Feed)

		if len(pay.Feed) > 0 {

			switch {
			case l > 14:
				limit = 2
			case l > 9 && l < 14:
				limit = 3
			case l > 5 && l <= 9:
				limit = 3
			case l > 2 && l <= 4:
				limit = 6
			case l == 2:
				limit = 8
			}

			for _, room := range pay.Feed {
				matrix.Prefix = "/_matrix/client/r0"

				state, err := matrix.RoomState(room.RoomID)
				if err != nil {
					log.Println(err)
					continue
				}
				matrix.Prefix = "/_matrix/client/"

				rc := c.RoomCreateEventFromState(state)

				lastEvent := room.LastEvent
				le := strconv.FormatInt(lastEvent, 10)

				opts := map[string]interface{}{
					"event_id":         rc,
					"room_id":          room.RoomID,
					"depth_first":      false,
					"recent_first":     false,
					"include_parent":   false,
					"include_children": true,
					"direction":        "down",
					"limit":            limit,
					"max_depth":        0,
					"max_breadth":      0,
					"last_event":       le,
				}

				relationships, err := matrix.GetRelationships(opts)
				if err != nil {
					log.Println(err)
					c.Error(w, r)
					return
				}
				if relationships != nil && len(relationships.Events) > 0 {

					processed := c.ProcessMessages(relationships.Events, state, user)
					posts = append(posts, processed...)
					lastEvent := relationships.Events[len(relationships.Events)-1].Timestamp
					feedItems = append(feedItems, &IndexFeed{
						RoomID:    room.RoomID,
						LastEvent: lastEvent,
					})
				}
			}
		}

		sort.Slice(posts, func(i, j int) bool { return posts[j].Time.Before(posts[i].Time) })

		type Response struct {
			Posts     interface{} `json:"posts"`
			FeedItems interface{} `json:"feed_items"`
		}

		ff := Response{
			Posts:     posts,
			FeedItems: feedItems,
		}

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

type IndexFeed struct {
	RoomID    string `json:"room_id"`
	LastEvent int64  `json:"last_event"`
}

func (c *Client) ProcessUserFeed(join *gomatrix.RespSync, user *User) ([]gomatrix.Event, []*IndexFeed, error) {

	publicRoom := ""
	for _, p := range user.JoinedRooms {
		if strings.Contains(p.RoomAlias, "#public:") {
			publicRoom = p.RoomID
		}
	}

	events := []gomatrix.Event{}
	feedItems := []*IndexFeed{}

	for roomID, room := range join.Rooms.Join {

		for i, _ := range user.JoinedRooms {
			if user.JoinedRooms[i].RoomID == roomID && user.JoinedRooms[i].RoomAlias != "" &&
				user.JoinedRooms[i].RoomID != publicRoom {
				for _, y := range room.Timeline.Events {
					c.ProcessEvent(&y, user)
					events = append(events, y)
				}
				/*
					feedItems = append(feedItems, &IndexFeed{
						RoomID:    roomID,
						LastEvent: room.Timeline.PrevBatch,
					})
				*/
			}
		}
	}
	sort.Slice(events, func(i, j int) bool { return events[i].Time.After(events[j].Time) })

	return events, feedItems, nil
}
