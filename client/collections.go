package client

import (
	"fmt"
	"hummingbard/gomatrix"
	"log"
	"net/http"
	"strings"

	"github.com/unrolled/secure"
)

func (c *Client) UserCollections(w http.ResponseWriter, r *http.Request) {

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

	if path[0] == '@' && (len(pi) > 2) {
		c.UserCollectionsItem(w, r)
		return
	}

	path = strings.Join(pi, "_")

	fed, use := c.IsFederated(path)

	//if the last item in path starts with '$' we know it's an event_id
	//permalink. we strip out the event_id and use the rest to check room path

	serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

	cli, err := gomatrix.NewClient(serverName, userid, token)
	if err != nil {
		log.Println(err)
		c.Error(w, r)
		return
	}

	var room string

	//build full room path
	//if path starts with @ then it's a user timeline, otherwise room
	room = fmt.Sprintf(`#%s:%s`, path, c.Config.Client.Domain)
	if fed {
		room = fmt.Sprintf(`#%s:%s`, use.LocalPart, use.ServerName)
	}

	//check if room exists
	var ra *gomatrix.RespAliasResolve

	ra, err = cli.ResolveAlias(room)
	if err != nil {
		log.Println(err)
	}

	var roomID string

	if ra != nil && len(ra.RoomID) > 0 {
		roomID = string(ra.RoomID)

	}

	// Default account should join this room
	err = c.OperatorJoinRoom(roomID)
	if err != nil {
		log.Println(err)
	}

	state, err := cli.RoomState(roomID)
	if err != nil {
		log.Println(err)
	}

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
				children = c.BuildSpaceChildren(roomID, spaces, false, false)
				pages = c.BuildSpaceChildren(roomID, spaces, true, false)
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
			PathItems: pitems,
			ID:        roomID,
			Children:  children,
			Pages:     pages,
		},
		RoomState: state,
	}

	{
		p := pathItems.Items
		t.Room.Path = p[0]
	}

	//build article item, cache it

	if c.Config.Mode == "development" {
		t.HomeServerURL = c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)
	} else {
		t.HomeServerURL = fmt.Sprintf(`https://%s`, c.Config.Matrix.Server)
	}

	srq := ProcessStateRequest{
		State: state,
		User:  us,
		Page:  &t,
	}

	c.ProcessState(&srq)

	t.LoggedInUser = us

	if us != nil {
		path := pathItems.Items
		if us.UserID == path[0] {
			t.IsOwner = true
		}
		if !strings.Contains(path[0], ":") {
			c := fmt.Sprintf(`%s:%s`, path[0], c.Config.Client.Domain)
			if us.UserID == c {
				t.IsOwner = true
			}
		}
	}

	nonce := secure.CSPNonce(r.Context())
	t.Nonce = nonce

	c.Templates.ExecuteTemplate(w, "collections", t)
}

func (c *Client) UserCollectionsItem(w http.ResponseWriter, r *http.Request) {

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

	lastItem := pi[len(pi)-1]

	path = strings.Join(pi, "_")

	fed, use := c.IsFederated(path)

	//if the last item in path starts with '$' we know it's an event_id
	//permalink. we strip out the event_id and use the rest to check room path

	serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

	cli, err := gomatrix.NewClient(serverName, userid, token)
	if err != nil {
		log.Println(err)
		c.Error(w, r)
		return
	}

	var room string

	//build full room path
	//if path starts with @ then it's a user timeline, otherwise room
	room = fmt.Sprintf(`#%s:%s`, path, c.Config.Client.Domain)
	if fed {
		room = fmt.Sprintf(`#%s:%s`, use.LocalPart, use.ServerName)
	}

	//check if room exists
	var ra *gomatrix.RespAliasResolve

	log.Println("what is room?", room)
	log.Println("what is room?", room)
	log.Println("what is room?", room)
	log.Println("what is room?", room)

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
				children = c.BuildSpaceChildren(roomID, spaces, false, false)
				pages = c.BuildSpaceChildren(roomID, spaces, true, false)
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
		RoomState: state,
	}
	processed := c.ProcessMessages(relationships, state, us)

	t.RootEvent = path

	pid := post.ID

	t.Posts = c.SortReplies(processed, pid, sort)

	t.Room.ThreadInRoomID = string(ra.RoomID)
	t.IsPermalink = true
	t.Room.EventID = post.ID

	t.PermalinkedPost = &pr[0]

	//build article item, cache it

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

	c.Templates.ExecuteTemplate(w, "collections", t)
}
