package client

import (
	"context"
	"encoding/json"
	"fmt"
	"hummingbard/gomatrix"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/unrolled/secure"
)

func (c *Client) CreateRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		us := LoggedInUser(r)

		type page struct {
			BasePage
			UserExists bool
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce
		t.LoggedInUser = us

		c.Templates.ExecuteTemplate(w, "create", t)
	}
}

func (c *Client) SpaceUsernameAvailable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Username     string `json:"username"`
			SubSpace     bool   `json:"sub_space"`
			ParentRoomID string `json:"parent_room_id"`
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

		username := strings.ToLower(pay.Username)

		if pay.SubSpace {
			state, err := matrix.RoomState(pay.ParentRoomID)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}

			alias := c.CanonicalAliasFromState(state)

			alias = alias[1:]
			s := strings.Split(alias, ":")
			localPart := s[0]
			username = fmt.Sprintf(`%s_%s`, localPart, strings.ToLower(pay.Username))
		}

		canon := fmt.Sprintf(`#%s:%s`, username, c.Config.Client.Domain)

		if user.Federated {
			canon = fmt.Sprintf(`#%s:%s`, username, user.HomeServer)
		}

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

func (c *Client) ValidateRoomCreation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Username     string `json:"username"`
			Title        string `json:"title"`
			About        string `json:"about"`
			Type         string `json:"type"`
			Private      bool   `json:"private"`
			NSFW         bool   `json:"nsfw"`
			SubSpace     bool   `json:"sub_space"`
			ParentRoomID string `json:"parent_room_id"`
			Page         bool   `json:"page"`
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

		type NewRoom struct {
			RoomID         string `json:"room_id"`
			Alias          string `json:"alias"`
			CanonicalAlias string `json:"canonical_alias"`
			Path           string `json:"path"`
		}

		type Response struct {
			Created bool    `json:"created"`
			Room    NewRoom `json:"room,omitempty"`
		}

		ff := Response{
			Created: false,
		}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		canon := fmt.Sprintf(`#%s:%s`, strings.ToLower(pay.Username), c.Config.Client.Domain)

		var state []*gomatrix.Event

		alias := ""

		if pay.SubSpace {
			state, err = matrix.RoomState(pay.ParentRoomID)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}

			alias = c.CanonicalAliasFromState(state)
			s := strings.Split(alias, ":")
			localPart := s[0]
			username := fmt.Sprintf(`%s_%s`, localPart, strings.ToLower(pay.Username))
			canon = fmt.Sprintf(`#%s:%s`, username, c.Config.Client.Domain)
		}

		_, err = c.Matrix.ResolveAlias(canon)
		if err != nil {
			//create the room

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
						"room_type": pay.Type,
					},
				}, gomatrix.Event{
					Type: "m.room.type",
					Content: map[string]interface{}{
						"type": "m.space",
					},
				},
				pl,
			}

			username := strings.ToLower(pay.Username)
			path := ""

			if pay.SubSpace {

				alias = alias[1:]
				s := strings.Split(alias, ":")
				localPart := s[0]
				username = fmt.Sprintf(`%s_%s`, localPart, strings.ToLower(pay.Username))

				x := strings.Split(localPart, "_")
				j := strings.Join(x, "/")

				path = fmt.Sprintf(`%s/%s`, j, strings.ToLower(pay.Username))

				content := map[string]interface{}{
					"canonical_alias": fmt.Sprintf(`#%s:%s`, username, user.HomeServer),
					"local_part":      username,
					"stripped":        strings.ToLower(pay.Username),
					"room_path":       path,
					"via":             []string{c.Config.Client.Domain},
				}

				if pay.Page {
					content["page"] = true
				}

				initState = append(initState, gomatrix.Event{
					Type:     fmt.Sprintf(`%s.parent`, c.Config.Spaces.Prefix),
					StateKey: &pay.ParentRoomID,
					Content:  content,
				})
			}

			if pay.NSFW {
				initState = append(initState, gomatrix.Event{
					Type: "com.hummingbard.room.nsfw",
					Content: map[string]interface{}{
						"nsfw": true,
					},
				})
			}

			req := &gomatrix.ReqCreateRoom{
				Preset:        "public_chat",
				Visibility:    "public",
				RoomAliasName: strings.ToLower(pay.Username),
				Name:          pay.Title,
				Topic:         pay.About,
				CreationContent: map[string]interface{}{
					"m.federate": true,
				},
				InitialState: initState,
			}

			if pay.SubSpace {
				req.RoomAliasName = username
			}

			if pay.Private {
				req.Preset = "private_chat"
				req.Visibility = "private"
			}

			crr, err := matrix.CreateRoom(req)

			if err != nil || crr == nil {
				log.Println(err)
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}
			ff.Created = true

			canon := fmt.Sprintf(`#%s:%s`, pay.Username, user.HomeServer)

			ff.Room = NewRoom{
				RoomID:         crr.RoomID,
				Alias:          strings.ToLower(pay.Username),
				CanonicalAlias: canon,
				Path:           path,
			}

			err = c.UpdateJoinedRooms(matrix, r)
			if err != nil {
				log.Println(err)
			}

			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, 3*time.Second)

			if !pay.Page {
				alias := fmt.Sprintf(`#%s:%s`, strings.ToLower(username), user.HomeServer)
				path := username
				if strings.Contains(path, "_") {
					s := strings.Split(path, "_")
					path = strings.Join(s, "/")
				}

				c.Cache.Rooms.Set(crr.RoomID, alias, 1)

			}

			if crr != nil {
				c.OperatorJoinRoom(crr.RoomID)

				if pay.Type != "gallery" && pay.Type == "page" {
					text, html := InitialMessage()
					npe := gomatrix.CreatePostEvent{
						RoomID:        crr.RoomID,
						Text:          text,
						FormattedText: html,
					}

					_, err := matrix.CreatePost(&npe)
					if err != nil {
						log.Println(err)
						log.Println(err)
					}
				}
				err = c.RefreshRoomsCache()
				if err != nil {
					log.Println(err)
				}

			}
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
