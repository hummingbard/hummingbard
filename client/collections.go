package client

import (
	"encoding/json"
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

	//check if @user exists
	{
		var room string
		//build full room path
		//if path starts with @ then it's a user timeline, otherwise room
		p := pathItems.Items
		room = fmt.Sprintf(`#%s:%s`, p[0], c.Config.Client.Domain)
		if fed {
			room = fmt.Sprintf(`#%s:%s`, use.LocalPart, use.ServerName)
		}

		//check if room exists
		_, err = cli.ResolveAlias(room)
		if err != nil {
			log.Println(err)
			c.NotFound(w, r)
			return
		}
	}

	var room string

	room = fmt.Sprintf(`#%s:%s`, path, c.Config.Client.Domain)
	if fed {
		room = fmt.Sprintf(`#%s:%s`, use.LocalPart, use.ServerName)
	}

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

		if spaces != nil && len(spaces.Events) > 0 {
			for _, space := range spaces.Events {
				log.Println(space)
			}
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

func (c *Client) CreateCollection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Name        string `json:"name"`
			Path        string `json:"path"`
			Description string `json:"description"`
		}

		var pay payload
		if r.Body == nil {
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
			Created    bool        `json:"created"`
			Collection interface{} `json:"collection"`
		}

		ff := Response{Created: false}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		//check if #@user_collections parent space exists, if not create it
		localPart := GetLocalPart(user.UserID)
		canon := fmt.Sprintf(`#@%s_%s:%s`, localPart, "collections", c.Config.Client.Domain)
		if user.Federated {
			canon = fmt.Sprintf(`#@%s_%s:%s`, localPart, "collections", user.HomeServer)
		}

		log.Println("canon is ", canon)
		log.Println("canon is ", canon)
		log.Println("canon is ", canon)

		av, err := matrix.ResolveAlias(canon)

		log.Println("does it exist", av)
		log.Println("does it exist", av)
		log.Println("does it exist", av)
		log.Println("does it exist", av)

		roomID := ""

		if av != nil {
			roomID = string(av.RoomID)
		}

		if av == nil {

			pl := gomatrix.Event{
				Type: "m.room.power_levels",
				Content: map[string]interface{}{
					"ban": 50,
					"events": map[string]interface{}{
						"m.room.name":         100,
						"m.room.power_levels": 100,
					},
					"events_default": 0,
					"invite":         50,
					"kick":           50,
					"notifications": map[string]interface{}{
						"room": 20,
					},
					"redact":        50,
					"state_default": 50,
					"users": map[string]interface{}{
						user.UserID:          100,
						c.DefaultUser.UserID: 100,
					},
					"users_default": 0,
				},
			}

			initState := []gomatrix.Event{
				gomatrix.Event{
					Type: "m.room.history_visibility",
					Content: map[string]interface{}{
						"history_visibility": "world_readable",
					},
				}, gomatrix.Event{
					Type: "m.room.guest_access",
					Content: map[string]interface{}{
						"guest_access": "can_join",
					},
				}, gomatrix.Event{
					Type: "com.hummingbard.room",
					Content: map[string]interface{}{
						"room_type": "collection",
					},
				}, gomatrix.Event{
					Type: "m.room.type",
					Content: map[string]interface{}{
						"type": "m.space",
					},
				}, gomatrix.Event{
					Type:     fmt.Sprintf(`%s.parent`, c.Config.Spaces.Prefix),
					StateKey: &user.RoomID,
					Content: map[string]interface{}{
						"via": []string{c.Config.Client.Domain},
					},
				},
				pl,
			}

			localPart := GetLocalPart(user.UserID)
			alias := fmt.Sprintf(`@%s_%s`, localPart, "collections")

			req := &gomatrix.ReqCreateRoom{
				Preset:        "public_chat",
				Visibility:    "public",
				RoomAliasName: alias,
				Name:          fmt.Sprintf(`%s's Collections`, localPart),
				Topic:         fmt.Sprintf(`%s's Collections`, localPart),
				CreationContent: map[string]interface{}{
					"m.federate": true,
				},
				InitialState: initState,
			}

			crr, err := matrix.CreateRoom(req)

			if err != nil || crr == nil {
				log.Println(err)
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}
			if crr != nil {
				roomID = crr.RoomID
				c.OperatorJoinRoom(crr.RoomID)
			}
		}
		pl := gomatrix.Event{
			Type: "m.room.power_levels",
			Content: map[string]interface{}{
				"ban": 50,
				"events": map[string]interface{}{
					"m.room.name":         100,
					"m.room.power_levels": 100,
				},
				"events_default": 0,
				"invite":         50,
				"kick":           50,
				"notifications": map[string]interface{}{
					"room": 20,
				},
				"redact":        50,
				"state_default": 50,
				"users": map[string]interface{}{
					user.UserID:          100,
					c.DefaultUser.UserID: 100,
				},
				"users_default": 0,
			},
		}

		alias := fmt.Sprintf(`@%s_%s_%s`, localPart, "collections", strings.ToLower(pay.Path))

		canon = fmt.Sprintf(`#%s:%s`, alias, user.HomeServer)

		initState := []gomatrix.Event{
			gomatrix.Event{
				Type: "m.room.history_visibility",
				Content: map[string]interface{}{
					"history_visibility": "world_readable",
				},
			}, gomatrix.Event{
				Type: "m.room.guest_access",
				Content: map[string]interface{}{
					"guest_access": "can_join",
				},
			}, gomatrix.Event{
				Type: "com.hummingbard.room",
				Content: map[string]interface{}{
					"room_type": "collection",
				},
			}, gomatrix.Event{
				Type: "m.room.type",
				Content: map[string]interface{}{
					"type": "m.space",
				},
			}, gomatrix.Event{
				Type:     fmt.Sprintf(`%s.parent`, c.Config.Spaces.Prefix),
				StateKey: &roomID,
				Content: map[string]interface{}{
					"via":             []string{c.Config.Client.Domain},
					"name":            pay.Name,
					"description":     pay.Description,
					"canonical_alias": canon,
					"stripped":        strings.ToLower(pay.Path),
				},
			},
			pl,
		}

		req := &gomatrix.ReqCreateRoom{
			Preset:        "public_chat",
			Visibility:    "public",
			RoomAliasName: alias,
			Name:          pay.Name,
			Topic:         pay.Description,
			CreationContent: map[string]interface{}{
				"m.federate": true,
			},
			InitialState: initState,
		}

		crr, err := matrix.CreateRoom(req)

		if err != nil || crr == nil {
			log.Println(err)
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}
		if crr != nil {
			roomID = crr.RoomID
			c.OperatorJoinRoom(crr.RoomID)

			ff.Created = true
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

func (c *Client) CollectionAvailable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Name string `json:"name"`
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
			Available bool `json:"available"`
		}
		ff := Response{Available: false}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		username := strings.ToLower(pay.Name)
		username = strings.ReplaceAll(username, " ", "")

		localPart := GetLocalPart(user.UserID)

		canon := fmt.Sprintf(`#@%s_%s_%s:%s`, localPart, "collections", username, c.Config.Client.Domain)

		if user.Federated {
			canon = fmt.Sprintf(`#@%s_%s_%s:%s`, localPart, "collections", username, user.HomeServer)
		}

		log.Println("canon is ", canon)
		log.Println("canon is ", canon)
		log.Println("canon is ", canon)

		av, err := matrix.ResolveAlias(canon)

		if av == nil {
			ff.Available = true
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
