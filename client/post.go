package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hummingbard/gomatrix"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func (c *Client) CreateNewPost() http.HandlerFunc {
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
			Post      struct {
				Content struct {
					Text string `json:"text"`
					HTML string `json:"html"`
				} `json:"content"`
				Links []struct {
					Href     string `json:"href"`
					Metadata struct {
						Title       string `json:"title"`
						Author      string `json:"author"`
						Description string `json:"description"`
						Image       string `json:"image"`
						IsYoutube   bool   `json:"is_youtube"`
						YoutubeID   string `json:"youtube_id"`
					} `json:"metadata"`
				} `json:"links"`
				Images []struct {
					Caption     string `json:"caption"`
					Description string `json:"description"`
					Filename    string `json:"filename"`
					Size        uint   `json:"size"`
					Mimetype    string `json:"mimetype"`
					MXC         string `json:"mxc"`
					Width       uint   `json:"width"`
					Height      uint   `json:"height"`
				} `json:"images"`
				Attachments []struct {
					Filename string `json:"filename"`
					Size     uint   `json:"size"`
					Mimetype string `json:"mimetype"`
					MXC      string `json:"mxc"`
				} `json:"attachments"`
				Article struct {
					Enabled       bool   `json:"enabled"`
					Title         string `json:"title"`
					Subtitle      string `json:"subtitle"`
					Description   string `json:"description"`
					CanonicalLink string `json:"canonical_link"`
					FeaturedImage *struct {
						Caption string `json:"caption"`
						MXC     string `json:"mxc"`
						Width   uint   `json:"width"`
						Height  uint   `json:"height"`
					} `json:"featured_image"`
				} `json:"article"`
			} `json:"post"`
			Page             bool   `json:"page"`
			Reply            bool   `json:"reply"`
			EventID          string `json:"event_id"`
			NSFW             bool   `json:"nsfw"`
			Anonymous        bool   `json:"anonymous"`
			SharedPostID     string `json:"shared_post_id"`
			Share            bool   `json:"share"`
			SharedPostRoomID string `json:"shared_post_room_id"`
			ShareReply       bool   `json:"share_reply"`
			ReplyPermalink   string `json:"reply_permalink"`
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
			Post interface{} `json:"post"`
		}

		fu, us := c.FederationUser(user.UserID)
		//port is only for my dev environment, this needs to go
		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		//if federation user, we query homeserver at the /well-known endpoint
		//for full server path
		if fu {
			wk, err := WellKnown(c.URLScheme(us.ServerName))
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}
			serverName = c.URLScheme(wk.ServerName)
		}
		if pay.Anonymous {
			serverName = c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		}

		userid := user.UserID
		accessToken := user.MatrixAccessToken

		matrix, err := gomatrix.NewClient(serverName, userid, accessToken)

		roomCreateEventID := ""

		if !pay.Reply {

			state, err := matrix.RoomState(pay.RoomID)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}

			rc := c.RoomCreateEventFromState(state)
			if len(rc) > 0 {
				roomCreateEventID = rc
			}
		}

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		post := &gomatrix.Event{}
		if pay.Share && len(pay.SharedPostID) > 0 {
			ev, err := c.Matrix.RoomEvent(pay.SharedPostRoomID, pay.SharedPostID)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}
			post = ev
		}

		_, err = matrix.JoinRoom(pay.RoomAlias, "", nil)
		if err != nil {
			log.Println(err)
		}

		threadRoomAlias := fmt.Sprintf(`thread%s`, RandomNumber(32))
		reAlias := fmt.Sprintf(`#%s:%s`, threadRoomAlias, user.HomeServer)

		roomID := ""

		text, err := SanitizeHTML(pay.Post.Content.Text)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		html, err := ToStrictHTML(pay.Post.Content.Text)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}
		if pay.Post.Article.Enabled {
			html, err = ToHTML(pay.Post.Content.Text)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}
		}

		//create post in new thread rroom
		npe := gomatrix.CreatePostEvent{
			RoomID:         roomID,
			RoomAlias:      pay.RoomAlias,
			Unsanitized:    pay.Post.Content.Text,
			ThreadInRoomID: pay.RoomID,
			NSFW:           pay.NSFW,
			Anonymous:      pay.Anonymous,
		}

		if !pay.Post.Article.Enabled {
			npe.Text = text
			npe.FormattedText = string(html)
		}

		if pay.Post.Article.Enabled {
			b := []byte(html)
			l := int64(len(b))
			r := bytes.NewReader(b)
			upl, err := matrix.UploadToContentRepo(r, "text/html", l)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			if upl != nil {
				npe.ArticleContent = upl.ContentURI
			}

		}

		slug := ""

		if pay.Post.Article.Enabled {
			npe.IsArticle = true
			npe.Title = pay.Post.Article.Title
			npe.Subtitle = pay.Post.Article.Subtitle
			npe.CanonicalLink = pay.Post.Article.CanonicalLink
			npe.Description = pay.Post.Article.Description
			if pay.Post.Article.FeaturedImage != nil &&
				len(pay.Post.Article.FeaturedImage.MXC) > 0 {

				f := (float32(pay.Post.Article.FeaturedImage.Height) / float32(pay.Post.Article.FeaturedImage.Width)) * 100
				asp := uint(f)
				y := gomatrix.Image{
					Caption:     pay.Post.Article.FeaturedImage.Caption,
					Width:       pay.Post.Article.FeaturedImage.Width,
					Height:      pay.Post.Article.FeaturedImage.Height,
					AspectRatio: asp,
					MXC:         StripMXCPrefix(pay.Post.Article.FeaturedImage.MXC),
				}
				npe.FeaturedImage = &y
			}

			slug = Slugify(pay.Post.Article.Title)

			roomPath := c.RoomPathFromAlias(pay.RoomAlias)

			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, 7*time.Second)
			exists, err := c.DoesSlugExist(ctx, roomPath, slug)
			if err != nil {
				log.Println(err)
			}

			if exists {
				slug = fmt.Sprintf(`%s-%s`, slug, RandomNumber(8))
			}

			npe.Slug = slug

			go func() {
				roomPath := c.RoomPathFromAlias(pay.RoomAlias)
				fullPath := filepath.Join(roomPath, slug)
				c.Cache.Articles.Set(fullPath, html, 1)
			}()
		}

		if post != nil && pay.Share && len(pay.SharedPostID) > 0 {
			c.ProcessEvent(post, user)
			npe.Share = true
			npe.SharedPost = post

			if pay.ShareReply && len(pay.ReplyPermalink) > 0 {
				npe.ShareReply = true
				npe.ReplyPermalink = pay.ReplyPermalink
				post.Content["share_reply"] = true
				post.Content["reply_permalink"] = pay.ReplyPermalink
			}
		}

		{

			rp := pay.RoomAlias
			if pay.RoomAlias[0] == '#' {
				rp = pay.RoomAlias[1:]
			}
			sp := strings.Split(rp, ":")
			if sp[1] != c.Config.Client.Domain {
				npe.Federated = true
			}
		}

		if len(pay.Post.Links) > 0 {
			links := []gomatrix.Link{}
			for _, x := range pay.Post.Links {
				y := gomatrix.Link{
					Href:        x.Href,
					Title:       x.Metadata.Title,
					Description: x.Metadata.Description,
					IsYoutube:   x.Metadata.IsYoutube,
					YoutubeID:   x.Metadata.YoutubeID,
				}

				links = append(links, y)
			}
			npe.Links = links
		}

		if len(pay.Post.Images) > 0 {
			images := []gomatrix.Image{}
			for _, x := range pay.Post.Images {
				f := (float32(x.Height) / float32(x.Width)) * 100
				asp := uint(f)
				y := gomatrix.Image{
					Caption:     x.Caption,
					Description: x.Description,
					Filename:    x.Filename,
					Size:        x.Size,
					Mimetype:    x.Mimetype,
					Width:       x.Width,
					Height:      x.Height,
					AspectRatio: asp,
					MXC:         StripMXCPrefix(x.MXC),
				}

				images = append(images, y)
			}
			npe.Images = images
		}

		if len(pay.Post.Attachments) > 0 {
			attachments := []gomatrix.Attachment{}
			for _, x := range pay.Post.Attachments {
				y := gomatrix.Attachment{
					Filename: x.Filename,
					Size:     x.Size,
					Mimetype: x.Mimetype,
					MXC:      StripMXCPrefix(x.MXC),
				}

				attachments = append(attachments, y)
			}
			npe.Attachments = attachments
		}

		npe.RoomID = pay.RoomID

		npe.Reply = pay.Reply
		if pay.Reply && pay.EventID != "" {
			npe.EventID = pay.EventID
		}

		if !pay.Reply {
			npe.EventID = roomCreateEventID
			npe.ThreadRoomID = roomID
			npe.ThreadRoomAlias = reAlias
			npe.ThreadInRoomID = ""
		}

		cre, err := matrix.CreatePost(&npe)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		if slug != "" {
			go func() {
				roomPath := c.RoomPathFromAlias(pay.RoomAlias)
				ctx := context.Background()
				ctx, _ = context.WithTimeout(ctx, 7*time.Second)
				_, err = c.UpdateEventSlug(ctx, roomPath, slug, cre.EventID)
				if err != nil {
					log.Println(err)
				}
			}()
		}

		//duplicate it in #public room if it isn't a thread reply
		if !pay.Reply && !pay.Page && !c.Config.Privacy.DisablePublic {

			go func() {
				npe.EventID = cre.EventID
				un := fmt.Sprintf(`#public:%s`, c.Config.Client.Domain)
				res, err := matrix.ResolveAlias(un)
				if err != nil || res == nil {
					log.Println(err)
				}

				if res != nil {
					npe.ThreadRoomID = ""
					npe.ThreadInRoomID = ""
					npe.Reply = false
					npe.RoomID = string(res.RoomID)
					if !strings.Contains(pay.RoomAlias, `#test`) {

						_, err = matrix.CreatePost(&npe)
						if err != nil {
							log.Println(err)
							http.Error(w, err.Error(), 400)
							return
						}
					}
					go func() {
						_, err = c.RefreshPublicEvents(string(res.RoomID))
						if err != nil {
							log.Println(err)
						}
					}()

				}

			}()
		}

		event, err := matrix.RoomEvent(pay.RoomID, cre.EventID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		c.ProcessEvent(event, user)

		ff := Response{
			Post: event,
		}

		if fu {

			go func() {
				time.Sleep(7 * time.Second)
				_, err = c.RefreshRoomEvents(pay.RoomID)
				if err != nil {
					log.Println(err)
				}
			}()
		} else {

			_, err = c.RefreshRoomEvents(pay.RoomID)
			if err != nil {
				log.Println(err)
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

func (c *Client) EditPost() http.HandlerFunc {
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
			EventID string `json:"event_id"`
			Content string `json:"content"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		fu, us := c.FederationUser(user.UserID)
		//port is only for my dev environment, this needs to go
		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		//if federation user, we query homeserver at the /well-known endpoint
		//for full server path
		if fu {
			wk, err := WellKnown(c.URLScheme(us.ServerName))
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}
			serverName = c.URLScheme(wk.ServerName)
		}

		matrix, err := gomatrix.NewClient(serverName, user.UserID, user.MatrixAccessToken)

		_, err = matrix.RoomEvent(pay.RoomID, pay.EventID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		text, err := SanitizeHTML(pay.Content)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		html, err := ToStrictHTML(pay.Content)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		npe := gomatrix.CreatePostEvent{
			RoomID:        pay.RoomID,
			Unsanitized:   pay.Content,
			Text:          text,
			FormattedText: string(html),
			Edit:          true,
			EditEventID:   pay.EventID,
		}

		_, err = matrix.CreatePost(&npe)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		type Response struct {
			Updated bool `json:"updated"`
		}

		ff := Response{
			Updated: false,
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

func (c *Client) ReactToPost() http.HandlerFunc {
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
			EventID string `json:"event_id"`
			Key     string `json:"key"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Reacted bool `json:"reacted"`
		}

		fu, us := c.FederationUser(user.UserID)
		//port is only for my dev environment, this needs to go
		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		//if federation user, we query homeserver at the /well-known endpoint
		//for full server path
		if fu {
			wk, err := WellKnown(c.URLScheme(us.ServerName))
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}
			serverName = c.URLScheme(wk.ServerName)
		}

		matrix, err := gomatrix.NewClient(serverName, user.UserID, user.MatrixAccessToken)

		//create post in new thread rroom
		npe := gomatrix.ReactToPostEvent{
			RoomID:  pay.RoomID,
			EventID: pay.EventID,
			Key:     pay.Key,
		}

		_, err = matrix.ReactToPost(&npe)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff := Response{
			Reacted: true,
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

func (c *Client) RedactPost() http.HandlerFunc {
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
			EventID string `json:"event_id"`
			Reason  string `json:"reason"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Redacted bool `json:"redacted"`
		}

		fu, us := c.FederationUser(user.UserID)
		//port is only for my dev environment, this needs to go
		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		//if federation user, we query homeserver at the /well-known endpoint
		//for full server path
		if fu {
			wk, err := WellKnown(c.URLScheme(us.ServerName))
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}
			serverName = c.URLScheme(wk.ServerName)
		}

		matrix, err := gomatrix.NewClient(serverName, user.UserID, user.MatrixAccessToken)

		_, err = matrix.RedactEvent(pay.RoomID, pay.EventID, &gomatrix.ReqRedact{
			Reason: pay.Reason,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff := Response{
			Redacted: true,
		}

		go func() {
			_, err = c.RefreshRoomEvents(pay.RoomID)
			if err != nil {
				log.Println(err)
			}
		}()

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

func (c *Client) ReportPost() http.HandlerFunc {
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
			EventID string `json:"event_id"`
			Reason  string `json:"reason"`
			Score   int    `json:"score"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Reported bool `json:"reported"`
		}

		fu, us := c.FederationUser(user.UserID)
		//port is only for my dev environment, this needs to go
		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		//if federation user, we query homeserver at the /well-known endpoint
		//for full server path
		if fu {
			wk, err := WellKnown(c.URLScheme(us.ServerName))
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}
			serverName = c.URLScheme(wk.ServerName)
		}

		matrix, err := gomatrix.NewClient(serverName, user.UserID, user.MatrixAccessToken)

		_, err = matrix.ReportEvent(pay.RoomID, pay.EventID, &gomatrix.ReqReport{
			Reason: pay.Reason,
			Score:  pay.Score,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff := Response{
			Reported: true,
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
