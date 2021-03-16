package client

import (
	"encoding/json"
	"fmt"
	"hummingbard/gomatrix"
	"log"
	"net/http"
	"strings"
)

func (c *Client) JoinRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Id    string `json:"id"`
			Alias string `json:"alias"`
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
			Joined bool `json:"joined"`
		}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff := Response{
			Joined: false,
		}

		jr, err := matrix.JoinRoom(pay.Alias, "", nil)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		} else if jr != nil {
			ff.Joined = true
		}

		rms, err := c.GetUserJoinedRooms(matrix)
		if err != nil {
			c.Error(w, r)
			return
		}

		err = c.RefreshJoinedRooms(r, rms)
		if err != nil {
			log.Println(err)
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

func (c *Client) LeaveRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Id    string `json:"id"`
			Alias string `json:"alias"`
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
			Left bool `json:"left"`
		}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff := Response{
			Left: false,
		}

		jr, err := matrix.LeaveRoom(pay.Id)
		if err != nil {
			log.Println(err)
		} else if jr != nil {
			ff.Left = true
		}

		rms, err := c.GetUserJoinedRooms(matrix)
		if err != nil {
			c.Error(w, r)
			return
		}

		err = c.RefreshJoinedRooms(r, rms)
		if err != nil {
			log.Println(err)
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

func (c *Client) GetRoomState() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			RoomID    string `json:"room_id"`
			RoomAlias string `json:"room_alias"`
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

		type Response struct {
			State interface{} `json:"state"`
		}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff := Response{}

		state, err := matrix.RoomState(pay.RoomID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff.State = state

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

func (c *Client) GetRoomMembers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			RoomID string `json:"room_id"`
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

		type member struct {
			DisplayName *string `json:"display_name"`
			AvatarURL   *string `json:"avatar_url"`
		}

		type Response struct {
			Members map[string]member `json:"members"`
		}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff := Response{Members: map[string]member{}}

		members, err := matrix.JoinedMembers(pay.RoomID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		for r, x := range members.Joined {
			if strings.Contains(r, `@anonymous`) ||
				strings.Contains(r, fmt.Sprintf(`@%s`, c.Config.Client.Domain)) {
				continue
			}
			ff.Members[r] = x
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

func (c *Client) UpdateRoomInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			RoomID  string `json:"room_id"`
			Profile bool   `json:"profile"`
			Info    struct {
				Title  string `json:"title"`
				About  string `json:"about"`
				Avatar string `json:"avatar"`
			} `json:"info"`
			Appearance struct {
				Header string `json:"header"`
				CSS    string `json:"css"`
			} `json:"appearance"`
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
			Updated bool `json:"updated"`
		}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff := Response{}

		owner := false
		if pay.Profile && user.RoomID == pay.RoomID {
			owner = true
		}

		if len(pay.Info.Title) > 0 {

			_, err = matrix.SendStateEvent(pay.RoomID, "m.room.name", "", map[string]interface{}{
				"name": pay.Info.Title,
			})
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}

			if owner {

				err = matrix.SetDisplayName(pay.Info.Title)
				if err != nil {
					log.Println(err)
					log.Println(err)
					log.Println(err)
					http.Error(w, err.Error(), 400)
					return
				}
				c.UpdateDisplayName(r, pay.Info.Title)
			}
		}

		if len(pay.Info.About) > 0 {
			_, err = matrix.SendStateEvent(pay.RoomID, "m.room.topic", "", map[string]interface{}{
				"topic": pay.Info.About,
			})
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}
		}

		_, err = matrix.SendStateEvent(pay.RoomID, "com.hummingbard.room.style", "", map[string]interface{}{
			"css": pay.Appearance.CSS,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		if len(pay.Appearance.Header) > 0 {
			_, err = matrix.SendStateEvent(pay.RoomID, "com.hummingbard.room.header", "", map[string]interface{}{
				"url": pay.Appearance.Header,
			})
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}
		}

		_, err = matrix.SendStateEvent(pay.RoomID, "m.room.avatar", "", map[string]interface{}{
			"url": pay.Info.Avatar,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}
		if owner {

			av := pay.Info.Avatar
			if len(pay.Info.Avatar) == 0 {
				av = "mxc://"
			}

			err = matrix.SetAvatarURL(av)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}
			c.UpdateAvatar(r, pay.Info.Avatar)
		}

		ff.Updated = true

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

func (c *Client) GetRoomInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			RoomID string `json:"room_id"`
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
			State interface{} `json:"state"`
		}

		userID := c.DefaultUser.UserID
		accessToken := c.DefaultUser.AccessToken

		if user != nil && len(user.UserID) > 0 {
			userID = user.UserID
			token = user.MatrixAccessToken
		}

		matrix, err := c.TempMatrixClient(userID, accessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff := Response{}

		state, err := matrix.RoomState(pay.RoomID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff.State = state

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

func (c *Client) FetchPublicSpaces() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type payload struct {
			Offset int `json:"offset"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Spaces interface{} `json:"spaces"`
		}

		rooms, err := c.GetPublicRoomsFromCache()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff := Response{}

		l := 43

		if pay.Offset != 0 {
			if len(rooms) > pay.Offset+l {
				ff.Spaces = rooms[pay.Offset : pay.Offset+l]
			} else {
				ff.Spaces = rooms[pay.Offset:]
			}
		} else {
			if len(rooms) >= l {
				ff.Spaces = rooms[:l]
			} else {
				ff.Spaces = rooms
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

func (c *Client) Queryspace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Query string `json:"query"`
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
			Spaces interface{} `json:"spaces"`
		}

		rooms, err := c.GetPublicRoomsFromCache()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		fed, use := FederationRoom(pay.Query)

		if fed {
			wk, err := WellKnown(c.URLScheme(use.ServerName))
			if err == nil {
				serverName := c.URLScheme(wk.ServerName)

				cli, err := gomatrix.NewClient(serverName, user.UserID, user.MatrixAccessToken)
				if err != nil {
					log.Println(err)
				}
				canon := pay.Query
				if canon[0] != '#' {
					canon = fmt.Sprintf(`#%s`, canon)
				}
				ra, err := cli.ResolveAlias(canon)
				if err != nil {
					log.Println(err)
				}
				if ra != nil {
					log.Println(ra.RoomID)
				}
			}
		}

		ff := Response{}

		rm := []*PublicRoom{}
		for _, room := range rooms {
			if strings.Contains(room.RoomAlias, pay.Query) ||
				strings.Contains(room.Name, pay.Query) ||
				strings.Contains(string(room.Topic), pay.Query) {
				rm = append(rm, room)
			}
		}
		ff.Spaces = rm

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}
