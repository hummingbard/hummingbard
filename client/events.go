package client

import (
	"fmt"
	"html/template"
	"hummingbard/gomatrix"
	"log"
	"sort"
	"strings"
	"time"
)

func (c *Client) SortRelationships(events []gomatrix.Event, rootID string) []*gomatrix.Event {
	items := []*gomatrix.Event{}

	for i, _ := range events {
		log.Println("type is", events[i].Type)
		if events[i].ID == rootID {
			items = append(items, &events[i])
		}
		if mr, ok := events[i].Content["m.relationship"].(map[string]interface{}); ok {
			if mr["event_id"] == rootID {

				items = append(items, &events[i])
			}
		}
	}
	var findReplies func(id string) []*gomatrix.Event

	findReplies = func(id string) []*gomatrix.Event {
		replies := []*gomatrix.Event{}

		for i, _ := range events {
			if mr, ok := events[i].Content["m.relationship"].(map[string]interface{}); ok {
				if mr["event_id"] == id {

					x := events[i]
					x.Replies = findReplies(events[i].ID)
					replies = append(replies, &x)
				}
			}
		}
		return replies
	}

	for _, event := range items {
		if event.ID == rootID {
			continue
		}
		event.Replies = findReplies(event.ID)
	}

	return items
}

func (c *Client) SortReplies(events []gomatrix.Event, rootID string, s string) []*gomatrix.Event {

	items := []*gomatrix.Event{}

	for i, _ := range events {
		if mr, ok := events[i].Content["m.relationship"].(map[string]interface{}); ok {
			if mr["event_id"] == rootID {

				items = append(items, &events[i])
			}
		}
	}
	var findReplies func(id string) ([]*gomatrix.Event, int)

	findReplies = func(id string) ([]*gomatrix.Event, int) {
		replies := []*gomatrix.Event{}

		for i, _ := range events {
			if mr, ok := events[i].Content["m.relationship"].(map[string]interface{}); ok {
				if mr["event_id"] == id {

					x := events[i]
					reps, count := findReplies(events[i].ID)
					x.Replies = reps
					x.TotalReplies = len(reps) + count
					replies = append(replies, &x)
				}
			}
		}
		c := 0
		for _, x := range replies {
			c += x.TotalReplies
		}
		sort.Slice(replies, func(i, j int) bool { return replies[j].Time.After(replies[i].Time) })
		return replies, c
	}

	for _, event := range items {
		if event.ID == rootID {
			continue
		}
		reps, count := findReplies(event.ID)
		event.Replies = reps
		event.TotalReplies = len(reps) + count
	}

	switch s {
	case "replies":
		sort.Slice(items, func(i, j int) bool { return items[i].TotalReplies > items[j].TotalReplies })
	case "recent":
		sort.Slice(items, func(i, j int) bool { return items[i].Time.After(items[j].Time) })
	default:
		sort.Slice(items, func(i, j int) bool { return items[j].Time.After(items[i].Time) })
	}

	return items
}

func (c *Client) ProcessMessages(resp []gomatrix.Event, state []*gomatrix.Event, user *User, showAll bool) []gomatrix.Event {

	members := []*gomatrix.Event{}

	for _, m := range state {
		if m.Type == "m.room.member" {
			members = append(members, m)
		}
	}

	events := []gomatrix.Event{}

	for i, _ := range resp {

		if !showAll && resp[i].Content["m.relationship"] == nil {
			continue
		}
		if showAll {
			reply, ok := resp[i].Content["reply"].(bool)
			if ok && reply {
				continue
			}
		}

		c.ProcessEvent(&resp[i], user)

		_, sok := resp[i].Content["shared_post"].(map[string]interface{})

		for m, _ := range members {
			if resp[i].Sender == members[m].Sender {
				dn, ok := members[m].Content["displayname"].(string)
				if ok {
					resp[i].Author.DisplayName = dn
				}
				au, ok := members[m].Content["avatar_url"].(string)
				if ok {
					resp[i].Author.AvatarURL = c.BuildAvatar(au)
				}

			}
			if sok && resp[i].SharedPost.Sender == members[m].Sender {

				dn, ok := members[m].Content["displayname"].(string)
				if ok {
					resp[i].SharedPost.Author.DisplayName = dn
				}
				au, ok := members[m].Content["avatar_url"].(string)
				if ok {

					resp[i].SharedPost.Author.AvatarURL = c.BuildAvatar(au)
				}
			}
		}

		if !showAll {
			if resp[i].Type == "com.hummingbard.post" {
				events = append(events, resp[i])
			}
		} else {
			if resp[i].Type == "com.hummingbard.post" || resp[i].Type == "m.room.message" {
				events = append(events, resp[i])
			}
		}

	}

	sort.Slice(events, func(i, j int) bool { return events[i].Time.After(events[j].Time) })

	resp = events
	return resp
}

func (c *Client) ProcessEvent(x *gomatrix.Event, user *User) {
	red, ok := x.Unsigned["redacted_by"].(string)
	if ok && len(red) > 0 {
		x.Redacted = true
		return
	}

	y := time.Unix(x.Timestamp/1000, x.Timestamp/10000*1000)
	x.Time = y
	ft := y.Format("Mon Jan 2 15:04:05 MST 2006")
	x.Date = ft
	x.When = FormatTime(y)

	//convert content body markdown to HTML
	body, ok := x.Content["formatted_body"].(string)
	if ok {
		x.Content["bodyHTML"] = template.HTML(body)
	} else {
		body, ok := x.Content["body"].(string)
		if ok {
			x.Content["bodyHTML"] = template.HTML(body)
		}
	}

	//format room_path
	roomPath, ok := x.Content["room_path"].(string)
	if ok {
		if strings.Contains(roomPath, "_") &&
			strings.Contains(roomPath, ":") {

			sp := strings.Split(roomPath, ":")

			xx := strings.Split(sp[0], "_")

			rem := xx[1:]
			rest := strings.Join(rem, "/")

			path := fmt.Sprintf(`%s:%s/%s`, xx[0], sp[1], rest)
			x.Content["room_path"] = path
		}
		x.Content["bodyHTML"] = template.HTML(body)
	}

	if strings.Contains(x.Sender, c.Config.Client.Domain) {
		s := strings.Split(x.Sender, ":")
		x.Author.FormattedID = s[0]
	} else {
		x.Author.FormattedID = x.Sender
	}

	if ch, ok := x.Unsigned["children"].(map[string]interface{}); ok {
		x.TotalReplies = int(ch["m.reference"].(float64)) - 1
	}

	shp, ok := x.Content["shared_post"].(map[string]interface{})

	if ok {
		sp := &gomatrix.Event{
			Sender:  shp["sender"].(string),
			ID:      shp["event_id"].(string),
			RoomID:  shp["room_id"].(string),
			Type:    shp["type"].(string),
			Content: shp["content"].(map[string]interface{}),
		}

		content, ok := shp["content"].(map[string]interface{})
		if ok {

			body, ok := content["formatted_body"]
			if ok {
				sp.Content["bodyHTML"] = template.HTML(body.(string))
			} else {
				body, ok := content["body"].(string)
				if ok {
					sp.Content["bodyHTML"] = template.HTML(body)
				}
			}

			article, ok := content["com.hummingbard.article"]
			if ok && article != nil {
				sp.IsArticle = true
			}
		}

		ts, ok := shp["origin_server_ts"].(float64)
		if ok {
			tst := int64(ts)
			sp.Timestamp = tst
			y := time.Unix(tst/1000, tst/10000*1000)
			ft := y.Format("Mon Jan 2 15:04:05 MST 2006")
			sp.Date = ft
			sp.When = FormatTime(y)
		}

		x.SharedPost = sp

	}

	if user != nil && len(user.UserID) > 0 {
		if user.UserID == x.Sender {
			x.Owner = true
		}
	}

	body, ok = x.Content["body"].(string)
	if ok {

		x.Content["body_length"] = len(body)
	}

	x.ShortID = fmt.Sprintf(`ev%s`, RandomString(9))

	t := false
	s := false

	article, ok := x.Content["com.hummingbard.article"].(map[string]interface{})
	if ok {
		if v, ok := article["title"].(string); ok && len(v) > 0 {
			t = true
		}
		if v, ok := article["slug"].(string); ok && len(v) > 0 {
			s = true
		}
		if v, ok := article["featured_image"].(map[string]interface{}); ok {
			if r, ok := v["mxc"].(string); ok && len(r) > 0 {
				v["mxc"] = c.BuildImage(r)
			}
		}
		if t && s {
			x.IsArticle = true
		}
	}
}
