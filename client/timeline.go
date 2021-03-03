package client

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"hummingbard/gomatrix"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/unrolled/secure"
)

type CachedRoomEvents struct {
	Events []gomatrix.Event
	Time   time.Time
}

type PathItem struct {
	Item string `json:"item"`
	Path string `json:"path"`
}

type ChildRoom struct {
	RoomID   string       `json:"room_id"`
	Alias    string       `json:"alias"`
	Path     string       `json:"path"`
	Children []*ChildRoom `json:"children"`
}

type Room struct {
	Path           string        `json:"path"`
	PathItems      []PathItem    `json:"path_items"`
	Alias          string        `json:"alias"`
	ID             string        `json:"id"`
	EventID        string        `json:"event_id"`
	ThreadInRoomID string        `json:"thread_in_room_id"`
	Description    template.HTML `json:"description"`
	Owner          User          `json:"owner"`
	Admin          string        `json:"admin"`
	CreatedAt      string        `json:"created_at"`
	Name           string        `json:"name"`
	Topic          template.HTML `json:"topic"`
	Type           string        `json:"type"`
	Avatar         string        `json:"avatar"`
	Members        int           `json:"members"`
	State          interface{}   `json:"state"`
	CSS            template.CSS  `json:"css"`
	Header         string        `json:"header"`
	Children       []*ChildRoom  `json:"children"`
	Pages          []*ChildRoom  `json:"pages"`
	Archaic        bool          `json:"archaic"`
}

type TimelinePage struct {
	Room             Room        `json:"room"`
	PublicRooms      interface{} `json:"public_rooms"`
	RoomState        interface{} `json:"room_state"`
	Posts            interface{} `json:"posts"`
	InitialPosts     interface{} `json:"initial_posts"`
	PermalinkedPost  interface{} `json:"permalinked_post"`
	IsPermalink      bool        `json:"is_permalink"`
	IsReplyPermalink bool        `json:"is_reply_permalink"`
	RootEvent        string      `json:"root_event"`
	IsArticle        bool        `json:"is_article"`
	IsPage           bool        `json:"is_page"`
	BasePage         `json:"base_page"`
	IsMember         bool        `json:"is_member"`
	HomeServerURL    string      `json:"home_server_url"`
	IsUserProfile    bool        `json:"is_user_profile"`
	IsOwner          bool        `json:"is_owner"`
	IsAdmin          bool        `json:"is_admin"`
	LastEvent        interface{} `json:"last_event"`
	Depth            int         `json:"depth"`
	Sort             string      `json:"sort"`
}

func (c *Client) BuildSpaceChildren(roomID string, spaces *gomatrix.RespSpaces, page bool) []*ChildRoom {

	rooms := []*ChildRoom{}

	for _, child := range spaces.Events {

		if *child.StateKey == roomID {
			isPage, ok := child.Content["page"].(bool)
			if page && (!ok || !isPage) {
				continue
			}
			if !page && (ok || isPage) {
				continue
			}

			y := ChildRoom{
				RoomID: child.RoomID,
			}
			if alias, ok := child.Content["canonical_alias"].(string); ok {
				y.Path = c.GetLocalPartPath(alias)
			}
			if stripped, ok := child.Content["stripped"].(string); ok {
				y.Alias = stripped
			}
			rooms = append(rooms, &y)
		}
	}

	var findChildren func(id string) []*ChildRoom

	findChildren = func(id string) []*ChildRoom {
		r := []*ChildRoom{}

		for _, child := range spaces.Events {
			if *child.StateKey == id {
				isPage, ok := child.Content["page"].(bool)
				if page && (!ok || !isPage) {
					continue
				}
				if !page && (ok || isPage) {
					continue
				}

				y := ChildRoom{
					RoomID: child.RoomID,
				}
				if alias, ok := child.Content["canonical_alias"].(string); ok {
					y.Path = c.GetLocalPartPath(alias)
				}
				if stripped, ok := child.Content["stripped"].(string); ok {
					y.Alias = stripped
				}
				y.Children = findChildren(child.RoomID)

				r = append(r, &y)
			}
		}
		return r
	}

	for _, room := range rooms {
		room.Children = findChildren(room.RoomID)
	}

	return rooms
}

func (c *Client) RejectLargeRoom(path string) bool {
	_, ok := c.Cache.LargeRooms.Get(path)
	if ok {
		return true
	}
	return false
}

func (c *Client) Timeline(w http.ResponseWriter, r *http.Request) {
	var token, userid string

	us := LoggedInUser(r)

	token = c.DefaultUser.AccessToken
	userid = c.DefaultUser.UserID

	serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

	// we're not using user token for viewing timeline events until peeking is
	// implemented in Dendrite
	/*
		if us != nil && us.AccessToken != "" {
			token = us.AccessToken
			userid = us.UserID
		}
	*/

	/*
		if fed && us != nil {
			wk, err := WellKnown(c.URLScheme(use.ServerName))
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}
			serverName = c.URLScheme(wk.ServerName)
		}
	*/

	cli, err := gomatrix.NewClient(serverName, userid, token)
	if err != nil {
		log.Println(err)
		c.Error(w, r)
		return
	}

	pathItems, err := c.PathItems(r)
	if err != nil {
		c.NotFound(w, r)
		return
	}

	path := pathItems.Path

	//temporary check for large rooms on matrix.org, need a better way to do this..
	{
		if strings.Contains(path, `matrix.org`) {
			if j := c.RejectLargeRoom(path); j {
				c.RoomTooLarge(w, r)
				return
			}
		}
	}

	//is this a federated room path?
	fed, use := c.IsFederated(path)

	//let's see if the path has extra items, eg. post permalink
	//split by /
	pi := strings.Split(path, "/")

	//if the last item in path starts with '$' we know it's an event_id
	//permalink. we strip out the event_id and use the rest to check room path
	lastItem := pi[len(pi)-1]
	secondLastItem := ""
	if len(pi) > 1 {
		secondLastItem = pi[len(pi)-2]
	}

	//if more than single item in path, it means we have a permalinked post
	if len(pi) > 1 {
		if lastItem[0] == '$' {
			{
				path := pathItems.Path
				sp := strings.Split(path, "/")
				sp = sp[:len(sp)-1]
				if secondLastItem[0] == '$' {
					sp = sp[:len(sp)-2]
				}
				if len(sp) > 1 {

					path = strings.Join(sp, "_")
					room := fmt.Sprintf(`#%s:%s`, path, c.Config.Client.Domain)
					_, err = cli.ResolveAlias(room)
					if err != nil {
						log.Println(err)
						c.PermalinkTimeline(w, r, true)
						return
					}
				}
			}
			c.PermalinkTimeline(w, r, false)
			return
		}

	}

	// we want to see if nested rooms exist, or else send to permalink
	var room string

	profileRoom := false

	//build full room path
	//if path starts with @ then it's a user timeline, otherwise room
	children := []*ChildRoom{}
	pages := []*ChildRoom{}

	if path[0] == '@' && len(pi) == 1 {
		//this is for local user's timeline room
		room = fmt.Sprintf(`#%s:%s`, path, c.Config.Client.Domain)
		//this is for users on other homeservers
		if fed {
			room = fmt.Sprintf(`#%s:%s`, use.LocalPart, use.ServerName)
		}
		profileRoom = true
	} else {
		sp := strings.Split(path, "/")
		path = strings.Join(sp, "_")
		room = fmt.Sprintf(`#%s:%s`, path, c.Config.Client.Domain)
		fed, use := FederationRoom(path)
		if fed {
			room = fmt.Sprintf(`#%s:%s`, use.LocalPart, use.ServerName)
		}
	}

	//let's check whether this is a slugged blog-post item
	{
		path := pathItems.Path
		sp := strings.Split(path, "/")

		if len(sp) > 1 {

			path = strings.Join(sp, "_")
			room = fmt.Sprintf(`#%s:%s`, path, c.Config.Client.Domain)

			path := pathItems.Items
			fed, use := FederationRoom(path[0])
			if fed {
				path := pathItems.Items[1:]
				rem := strings.Join(path, "_")
				j := use.LocalPart + "_" + rem
				room = fmt.Sprintf(`#%s:%s`, j, use.ServerName)
			}

			_, err := cli.ResolveAlias(room)
			if err != nil {
				log.Println(err)
				c.PermalinkTimeline(w, r, true)
				return
			}
		}
	}

	//check if it's a user profile timeline
	var profile *gomatrix.RespProfile
	if fed {
		profile, err = cli.GetProfile(path)
		if err != nil {
			log.Println(err)
			c.NotFound(w, r)
			return
		}
	}

	//check if room exists
	var ra *gomatrix.RespAliasResolve

	ra, err = cli.ResolveAlias(room)
	if err != nil {
		log.Println(err)
		if profile != nil {
			c.UserNotRegistered(w, r, profile, path, use)
			return
		} else {
			c.NotFound(w, r)
			return
		}
	}
	log.Println("resp is ", ra.RoomID)

	roomID := string(ra.RoomID)

	messagesFrom := string(ra.RoomID)

	// Default account should join this room
	err = c.OperatorJoinRoom(roomID)
	if err != nil {
		log.Println(err)
	}

	//grab room state
	state, err := cli.RoomState(roomID)
	if err != nil {
		log.Println(err)
		c.NotFound(w, r)
		return
	}

	//get room messages
	events := []gomatrix.Event{}

	rc := c.RoomCreateEventFromState(state)

	cli.Prefix = "/_matrix/client/"

	opts := map[string]interface{}{
		"event_id":         rc,
		"room_id":          string(ra.RoomID),
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

	if relationships != nil && len(relationships.Events) > 0 {
		events = relationships.Events
	}

	posts := c.ProcessMessages(events, state, us)

	isPage := c.IsPage(state)

	//if this a page, we want to pull in state from the parent room so that
	//sidebar info is from parent
	if len(pi) > 1 && isPage {
		newRoom := ""
		pi = pi[:len(pi)-1]
		path = strings.Join(pi, "/")
		if path[0] == '@' {
			//this is for local user's timeline room
			newRoom = fmt.Sprintf(`#%s:%s`, path, c.Config.Client.Domain)
			//this is for users on other homeservers
			if fed {
				newRoom = fmt.Sprintf(`#%s:%s`, use.LocalPart, use.ServerName)
			}
			profileRoom = true
		} else {
			sp := strings.Split(path, "/")
			path = strings.Join(sp, "_")
			newRoom = fmt.Sprintf(`#%s:%s`, path, c.Config.Client.Domain)
			path := pathItems.Items
			fed, use := FederationRoom(path[0])
			if fed {
				path := pathItems.Items[:len(pathItems.Items)-1]
				if len(path) == 1 {
					newRoom = fmt.Sprintf(`#%s`, path[0])
				} else {
					path := pathItems.Items[1:]
					rem := strings.Join(path, "_")
					j := use.LocalPart + "_" + rem
					newRoom = fmt.Sprintf(`#%s:%s`, j, use.ServerName)
				}
			}
		}
		ra, err = cli.ResolveAlias(newRoom)
		if err != nil {
			log.Println(err)
			if profile != nil {
				c.UserNotRegistered(w, r, profile, path, use)
				return
			} else {
				c.NotFound(w, r)
				return
			}
		}

		roomID = string(ra.RoomID)
		state, err = cli.RoomState(roomID)
		if err != nil {
			log.Println(err)
			c.NotFound(w, r)
			return
		}
	}

	//let's pull in space info
	{

		cli.Prefix = "/_matrix/client/"
		spaces, err := cli.Spaces(roomID, &gomatrix.ReqSpaces{
			MaxRoomsPerSpace: 100,
			Limit:            666,
		})
		if err != nil {
			log.Println(err)
		}

		if spaces != nil {
			if len(spaces.Events) > 0 {
				children = c.BuildSpaceChildren(roomID, spaces, false)
				pages = c.BuildSpaceChildren(roomID, spaces, true)
			}
		}

		cli.Prefix = "/_matrix/client/r0"
	}

	//build meta-header nav items
	path = pathItems.Path
	pitems := []PathItem{}
	sp := strings.Split(path, "/")
	for i, x := range sp {
		pitems = append(pitems, PathItem{
			Item: x,
			Path: strings.Join(sp[:i+1], "/"),
		})
	}

	//page construct
	t := TimelinePage{
		Room: Room{
			Path:      path,
			PathItems: pitems,
			ID:        messagesFrom,
			Children:  children,
			Pages:     pages,
		},
		Posts:         posts,
		IsUserProfile: profileRoom,
	}

	if len(events) > 0 {
		t.LastEvent = events[len(events)-1].Timestamp
	}

	if c.Config.Mode == "development" {
		t.HomeServerURL = c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)
	} else {
		t.HomeServerURL = fmt.Sprintf(`https://%s`, c.Config.Matrix.Server)
	}

	// has logged in user joined the room?
	if us != nil && len(us.JoinedRooms) > 0 {
		for i, _ := range us.JoinedRooms {
			if us.JoinedRooms[i].RoomID == t.Room.ID {
				t.IsMember = true
			}
		}
	}

	//process room state and add to room construct
	srq := ProcessStateRequest{
		State: state,
		User:  us,
		Page:  &t,
	}
	c.ProcessState(&srq)

	//if this is a page we only want the first item
	if isPage {
		t.IsPage = true
		if len(posts) > 0 {
			t.Posts = t.Posts.([]gomatrix.Event)[:1]
		}
	}

	t.LoggedInUser = us

	nonce := secure.CSPNonce(r.Context())
	t.Nonce = nonce

	ip, err := json.Marshal(t.Posts)
	if err != nil {
		c.Error(w, r)
		return
	}

	t.InitialPosts = string(ip)

	if t.Room.Type == "gallery" {
		c.Templates.ExecuteTemplate(w, "gallery", t)
		return
	}

	c.Templates.ExecuteTemplate(w, "timeline", t)
}

func (c *Client) UserNotRegistered(w http.ResponseWriter, r *http.Request, p *gomatrix.RespProfile, id string, u *UserID) {
	user := LoggedInUser(r)

	type profile struct {
		DisplayName *string
		UserID      string
		AvatarURL   *string
		ServerName  string
	}

	type page struct {
		BasePage
		LoggedInUser *User
		User         profile
	}

	nonce := secure.CSPNonce(r.Context())

	t := &page{
		User: profile{
			UserID:      id,
			DisplayName: p.Displayname,
			ServerName:  u.ServerName,
		},
	}
	if p.AvatarURL != nil {
		a := c.BuildAvatar(*p.AvatarURL)
		t.User.AvatarURL = &a
	}

	t.Nonce = nonce
	t.LoggedInUser = user

	c.Templates.ExecuteTemplate(w, "user-not-registered", t)
}

func (c *Client) PermalinkTimeline(w http.ResponseWriter, r *http.Request, slugged bool) {
	var token, userid string

	us := LoggedInUser(r)

	token = c.DefaultUser.AccessToken
	userid = c.DefaultUser.UserID

	pathItems, err := c.PathItems(r)
	if err != nil {
		c.NotFound(w, r)
		return
	}

	path := pathItems.Path

	//let's see if the path has extra items, eg. post permalink
	//split by /
	pi := strings.Split(path, "/")

	//we know this is a reply permalink if the last two items are events
	//starting with '$'
	replyPermalink := len(pi) >= 3 && pi[len(pi)-2][0] == '$' && pi[len(pi)-1][0] == '$'

	if slugged {
		replyPermalink = len(pi) >= 3 && pi[len(pi)-2][0] != '$' && pi[len(pi)-1][0] == '$'
	}

	lastItem := pi[len(pi)-1]
	lastEvent := ""
	if replyPermalink {
		lastEvent = pi[len(pi)-1]
	}

	if !replyPermalink {
		pi = pi[:len(pi)-1]
	}

	path = strings.Join(pi, "/")

	if replyPermalink {
		lastItem = pi[len(pi)-2]
		p := pi[:len(pi)-2]
		path = strings.Join(p, "/")
	}

	fed, use := c.IsFederated(path)

	//if the last item in path starts with '$' we know it's an event_id
	//permalink. we strip out the event_id and use the rest to check room path

	//if slugged, let's get the event id from DB
	if slugged {
		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		eventID, err := c.GetSlugEventID(ctx, path, lastItem)
		if err != nil {
			log.Println(err)
			c.NotFound(w, r)
			return
		}

		if eventID != "" {
			lastItem = eventID
		}
	}

	serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

	cli, err := gomatrix.NewClient(serverName, userid, token)
	if err != nil {
		log.Println(err)
		c.Error(w, r)
		return
	}

	var room string

	profileRoom := false

	//build full room path
	//if path starts with @ then it's a user timeline, otherwise room
	if path[0] == '@' {
		room = fmt.Sprintf(`#%s:%s`, path, c.Config.Client.Domain)
		if fed {
			room = fmt.Sprintf(`#%s:%s`, use.LocalPart, use.ServerName)
		}
		profileRoom = true
	} else {
		sp := strings.Split(path, "/")
		pth := strings.Join(sp, "_")
		if replyPermalink {
			p := pi[:len(pi)-2]
			pth = strings.Join(p, "_")
		}
		room = fmt.Sprintf(`#%s:%s`, pth, c.Config.Client.Domain)
		fed, use := FederationRoom(pth)

		if fed {
			room = fmt.Sprintf(`#%s:%s`, use.LocalPart, use.ServerName)
		}
	}

	//check if room exists
	var ra *gomatrix.RespAliasResolve

	ra, err = cli.ResolveAlias(room)
	if err != nil {
		log.Println(err)
		c.NotFound(w, r)
		return
	}

	roomID := string(ra.RoomID)

	// Default account should join this room
	err = c.OperatorJoinRoom(roomID)
	if err != nil {
		log.Println(err)
	}

	post := &gomatrix.Event{}

	event, err := cli.RoomEvent(roomID, lastItem)
	if err != nil {
		log.Println(err)
		c.NotFound(w, r)
		return
	}

	state, err := cli.RoomState(roomID)
	if err != nil {
		log.Println(err)
		c.NotFound(w, r)
		return
	}

	events := []gomatrix.Event{*event}
	pr := c.ProcessMessages(events, state, us)
	//c.ProcessEvent(event, us)
	if len(pr) == 0 {
		c.NotFound(w, r)
		return
	} else {

		post = &pr[0]
	}

	relationships := []gomatrix.Event{}

	cli.Prefix = "/_matrix/client/"

	query := r.URL.Query()
	sort := query.Get("sort")

	opts := map[string]interface{}{
		"event_id":         lastItem,
		"room_id":          roomID,
		"depth_first":      false,
		"recent_first":     false,
		"include_parent":   false,
		"include_children": true,
		"direction":        "down",
		"limit":            500,
		"max_depth":        2,
		"max_breadth":      7,
	}

	if replyPermalink {
		opts["event_id"] = lastEvent
		opts["max_depth"] = 100
		opts["max_breadth"] = 100
	}

	rela, err := cli.GetRelationships(opts)
	if err != nil {
		log.Println(err)
	}

	relationships = rela.Events

	cli.Prefix = "/_matrix/client/r0"

	children := []*ChildRoom{}
	pages := []*ChildRoom{}
	{

		cli.Prefix = "/_matrix/client/"
		spaces, err := cli.Spaces(roomID, &gomatrix.ReqSpaces{
			MaxRoomsPerSpace: 100,
			Limit:            666,
		})
		if err != nil {
			log.Println(err)
		}

		if spaces != nil {
			if len(spaces.Events) > 0 {
				children = c.BuildSpaceChildren(roomID, spaces, false)
				pages = c.BuildSpaceChildren(roomID, spaces, true)
			}
		}

		cli.Prefix = "/_matrix/client/r0"
	}

	path = pathItems.Path

	pitems := []PathItem{}
	sp := strings.Split(path, "/")
	sp = sp[:len(sp)-1]
	for i, x := range sp {
		if strings.Contains(x, "$") {
			continue
		}
		pitems = append(pitems, PathItem{
			Item: x,
			Path: strings.Join(sp[:i+1], "/"),
		})
	}

	t := TimelinePage{
		Room: Room{
			Path:      path,
			PathItems: pitems,
			ID:        roomID,
			Children:  children,
			Pages:     pages,
		},
		IsUserProfile: profileRoom,
	}
	processed := c.ProcessMessages(relationships, state, us)

	t.RootEvent = path

	pid := post.ID
	if replyPermalink {
		t.IsReplyPermalink = true

		path := pathItems.Path
		x := strings.Split(path, "/")
		x = x[:len(x)-1]
		t.RootEvent = strings.Join(x, "/")

		ev, err := cli.RoomEvent(roomID, lastEvent)

		if err != nil {
			log.Println(err)
			c.NotFound(w, r)
			return
		}

		rel, ok := ev.Content["m.relationship"].(map[string]interface{})
		if ok {
			pid = rel["event_id"].(string)
		}
	}

	log.Println("pid is ", pid)
	t.Posts = c.SortReplies(processed, pid, sort)

	t.Room.ThreadInRoomID = string(ra.RoomID)
	t.IsPermalink = true
	t.Room.EventID = post.ID

	t.PermalinkedPost = &pr[0]

	//build article item, cache it
	if slugged {
		article := &pr[0]

		path := pathItems.Path

		content, ok := c.Cache.Articles.Get(path)
		if ok {
			article.Content["bodyHTML"] = content
		} else {

			if v, ok := article.Content["com.hummingbard.article"].(map[string]interface{}); ok {
				if cu, ok := v["content_uri"].(string); ok {
					cli.Prefix = "/_matrix/media/r0"

					dow, err := cli.Download(StripMXCPrefix(cu), "")
					if err != nil {
						log.Println(err)
						c.Error(w, r)
						return
					}

					article.Content["bodyHTML"] = template.HTML(dow)
					c.Cache.Articles.Set(path, template.HTML(dow), 1)
				}
			}
		}

	}

	if c.Config.Mode == "development" {
		t.HomeServerURL = c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)
	} else {
		t.HomeServerURL = fmt.Sprintf(`https://%s`, c.Config.Matrix.Server)
	}

	if us != nil && len(us.JoinedRooms) > 0 {
		for i, _ := range us.JoinedRooms {
			if us.JoinedRooms[i].RoomID == t.Room.ID || us.JoinedRooms[i].RoomID == event.RoomID {
				t.IsMember = true

			}
		}
	}

	srq := ProcessStateRequest{
		State: state,
		User:  us,
		Page:  &t,
	}

	c.ProcessState(&srq)

	t.LoggedInUser = us

	nonce := secure.CSPNonce(r.Context())
	t.Nonce = nonce

	ip, err := json.Marshal(processed)
	if err != nil {
		c.Error(w, r)
		return
	}

	t.Sort = sort

	t.InitialPosts = string(ip)

	if t.Room.Type == "gallery" {
		c.Templates.ExecuteTemplate(w, "gallery-item", t)
		return
	}

	if post.IsArticle {
		t.IsArticle = true
		c.Templates.ExecuteTemplate(w, "article", t)
		return
	}

	c.Templates.ExecuteTemplate(w, "timeline", t)
}

func (c *Client) FetchReplies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
		}

		type payload struct {
			RoomID       string `json:"room_id"`
			ThreadRoomID string `json:"thread_room_id"`
			EventID      string `json:"event_id"`
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

		type Response struct {
			Replies  interface{} `json:"replies"`
			Unsorted interface{} `json:"unsorted"`
		}

		matrix, err := c.TempMatrixClient(c.DefaultUser.UserID, c.DefaultUser.AccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		state, err := matrix.RoomState(pay.RoomID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		matrix.Prefix = "/_matrix/client/"

		query := r.URL.Query()
		sort := query.Get("sort")

		opts := map[string]interface{}{
			"event_id":         pay.EventID,
			"room_id":          pay.ThreadRoomID,
			"depth_first":      false,
			"recent_first":     false,
			"include_parent":   false,
			"include_children": true,
			"direction":        "down",
			"limit":            500,
			"max_depth":        1,
			"max_breadth":      53,
		}

		rel, err := matrix.GetRelationships(opts)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff := Response{}

		if rel != nil {
			processed := c.ProcessMessages(rel.Events, state, user)
			replies := c.SortReplies(processed, pay.EventID, sort)
			ff.Unsorted = processed
			ff.Replies = replies
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

func (c *Client) GetMoreMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
		}

		type payload struct {
			Id        string      `json:"id"`
			End       interface{} `json:"end"`
			Permalink bool        `json:"permalink"`
			Public    bool        `json:"public"`
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

		type Response struct {
			Posts     interface{} `json:"posts"`
			LastEvent interface{} `json:"last_event"`
		}

		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		cli, err := gomatrix.NewClient(serverName, c.DefaultUser.UserID, c.DefaultUser.AccessToken)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		state, err := cli.RoomState(pay.Id)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}
		rc := c.RoomCreateEventFromState(state)

		cli.Prefix = "/_matrix/client/"

		opts := map[string]interface{}{
			"event_id":         rc,
			"room_id":          string(pay.Id),
			"depth_first":      false,
			"recent_first":     true,
			"include_parent":   false,
			"include_children": true,
			"direction":        "down",
			"limit":            5,
			"max_depth":        0,
			"max_breadth":      0,
		}

		if ev, ok := pay.End.(float64); ok {
			le := strconv.FormatInt(int64(ev), 10)
			opts["last_event"] = le
		}

		events := []gomatrix.Event{}

		relationships, err := cli.GetRelationships(opts)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}
		cli.Prefix = "/_matrix/client/r0"

		var lastEvent interface{}

		lastEvent = relationships.Events[len(relationships.Events)-1].Timestamp

		events = relationships.Events

		if pay.Public {

			if ev, ok := pay.End.(string); ok {

				fil := `{types:["com.hummingbard.post"]}`
				msg, err := cli.Messages(string(pay.Id), ev, "", 'b', 13, fil)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), 400)
					return
				}
				events = msg.Chunk
				lastEvent = msg.End
			}
		}

		ff := Response{
			Posts:     c.ProcessMessages(events, state, user),
			LastEvent: lastEvent,
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

func (c *Client) OperatorJoinRoom(roomID string) error {
	//see whether default accout has already joined this room
	rooms, ok := c.Cache.Rooms.Get(c.DefaultUser.UserID)
	if ok {
		for i, _ := range rooms.([]string) {
			if roomID == rooms.([]string)[i] {
				return nil
			}
		}
	}

	_, err := c.Matrix.JoinRoom(roomID, "", nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
