package client

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

func WellKnown(s string) (*WellKnownServer, error) {
	resp, err := http.Get(fmt.Sprintf(`%s/.well-known/matrix/server`, s))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var res WellKnownServer

	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

type UserID struct {
	LocalPart  string
	ServerName string
}

//Determind whether username belongs to this homesever or an external one
//external usernames should be in full @username:homeserver.org form
//local usernames should be @username only for better UI
func (c *Client) FederationUser(username string) (bool, *UserID) {
	//validUsernameRegex := regexp.MustCompile(`^@[0-9a-z_\-=./]+:[0-9a-z_\-=./]+\.[a-z]{2,}$`)
	validUsernameRegex := regexp.MustCompile(`^@.+?:.+$`)
	if validUsernameRegex.MatchString(username) {
		username = username[1:]
		parts := strings.Split(username, ":")

		if parts[1] == c.Config.Matrix.Server {
			return false, nil
		}
		return true, &UserID{
			LocalPart:  parts[0],
			ServerName: parts[1],
		}
	}

	return false, nil
}

func FederationRoom(username string) (bool, *UserID) {
	//validUsernameRegex := regexp.MustCompile(`^@[0-9a-z_\-=./]+:[0-9a-z_\-=./]+\.[a-z]{2,}$`)
	validUsernameRegex := regexp.MustCompile(`^.+?:.+$`)
	if validUsernameRegex.MatchString(username) {
		parts := strings.Split(username, ":")
		log.Println(parts)
		return true, &UserID{
			LocalPart:  parts[0],
			ServerName: parts[1],
		}
	}

	return false, nil
}

func GetHomeServerPart(s string) string {
	if strings.Contains(s, ":") {
		sp := strings.Split(s, ":")
		return sp[len(sp)-1]
	}
	return s
}

func (c *Client) IsFederated(username string) (bool, *UserID) {
	//federated user paths should have the same format as email, like so: username@homeserver.com
	//obviously a very loose regexp
	//validUsernameRegex := regexp.MustCompile(`^.+?@.+$`)
	validUsernameRegex := regexp.MustCompile(`^@.+?:.+$`)
	if validUsernameRegex.MatchString(username) {

		//lets's split the localpart and server_name
		parts := strings.Split(username, ":")

		//if severname is the same as out homeserver, return
		if parts[1] == c.Config.Client.Domain ||
			parts[1] == c.Config.Matrix.Server {
			return false, nil
		}

		//return federated path
		return true, &UserID{
			LocalPart:  parts[0],
			ServerName: parts[1],
		}
	}

	return false, nil
}

func FileID(fileID string) string {
	fi := strings.Replace(fileID, "mxc://", "", 1)
	sp := strings.Split(fi, "/")
	return sp[1]
}

func (c *Client) URLScheme(url string) string {
	if c.Config.Matrix.Server != url {
		return fmt.Sprintf(`https://%s`, url)
	}
	return fmt.Sprintf(`http://%s`, url)
}

func UnsafeHTML(x string) (template.HTML, error) {
	unsafe := blackfriday.Run([]byte(x))
	return template.HTML(unsafe), nil
}

func ToHTML(x string) (template.HTML, error) {
	unsafe := blackfriday.Run([]byte(x))
	p := bluemonday.UGCPolicy()
	safe := p.Sanitize(string(unsafe))
	return template.HTML(safe), nil
}

func ToStrictHTML(x string) (template.HTML, error) {
	unsafe := blackfriday.Run([]byte(x))

	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	p.RequireParseableURLs(true)
	p.AllowRelativeURLs(true)

	p.AllowStandardAttributes()

	p.AllowImages()

	p.AllowURLSchemes("mailto", "https")

	p.AllowAttrs("href").OnElements("a")

	p.AllowElements("blockquote")

	p.AllowElements("p")
	p.AllowElements("b", "strong")
	p.AllowElements("i", "em")
	p.AllowAttrs("class").OnElements("span")

	p.AllowElements("br")

	p.AllowElements("hr")
	p.AllowElements("ul")
	p.AllowElements("ol")
	p.AllowElements("li")
	p.AllowElements("br")

	p.AllowAttrs("id").OnElements("li")
	p.AllowAttrs("class").OnElements("li")

	p.AllowElements("sub")
	p.AllowElements("sup")

	p.AllowElements("s")
	p.AllowElements("del")

	p.AllowElements("pre")
	p.AllowElements("code")

	safe := p.Sanitize(string(unsafe))
	return template.HTML(safe), nil
}

func SanitizeHTML(x string) (string, error) {
	p := bluemonday.StrictPolicy()
	p.AllowElements("br")
	safe := p.Sanitize(x)
	return safe, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const NumberBytes = "0123456789"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func RandomString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func RandomNumber(n int) string {
	b := make([]byte, n)

	src := rand.NewSource(time.Now().UnixNano())
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(NumberBytes) {
			b[i] = NumberBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func StripMXCPrefix(s string) string {
	s = strings.Replace(s, "mxc://", "", -1)
	return s
}

func (c *Client) RoomPathFromAlias(alias string) string {
	federated := false
	sp := strings.Split(alias, ":")
	if sp[1] != c.Config.Client.Domain {
		federated = true
	}

	path := ""
	rp := alias[1:]

	if !federated {
		sp := strings.Split(rp, ":")
		s := strings.Split(sp[0], "_")
		if len(sp) > 1 {
			p := strings.Join(s, "/")
			path = p
		} else {
			path = s[0]
		}
	} else {
		path = rp
	}
	return path
}

func FormatTime(t time.Time) string {

	thisYear := "Jan _2"
	pastYears := "Jan _2, 2006"
	//max := 24 * time.Hour

	now := time.Now()

	difference := now.Sub(t)

	// If it's within last 12 hours

	if difference < time.Minute {
		return "Just Now"
	}

	if difference < time.Hour {
		difference = difference.Round(time.Minute)
		x := math.Trunc(difference.Minutes())
		return fmt.Sprintf(`%.fm`, x)
	}

	if difference <= time.Hour*23 {
		difference = difference.Round(time.Hour)
		x := math.Trunc(difference.Hours())
		return fmt.Sprintf(`%.fh`, x)
	}

	if t.Year() == now.Year() {
		x := t.Format(thisYear)
		return x
	}

	return t.Format(pastYears)
}

func InitialMessage() (string, string) {
	plain_text := `If I could write the beauty of your eyes,
			And in fresh numbers number all your graces,
			The age to come would say ‘This poet lies;
			Such heavenly touches ne’er touch’d earthly faces.’`

	html := `If I could write the beauty of your eyes,<br>
			And in fresh numbers number all your graces,<br>
			The age to come would say ‘<em>This poet lies</em>;<br>
			Such heavenly touches ne’er touch’d earthly faces.’<br>`
	return plain_text, html
}

func (c *Client) BuildDownloadLink(mxc string) string {

	avurl := StripMXCPrefix(mxc)

	if len(avurl) == 0 {
		return ""
	}

	serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

	avurl = fmt.Sprintf(`%s/_matrix/media/r0/download/%s`, serverName, avurl)

	if c.Config.Mode == "production" {
		serverName = c.Config.Matrix.Server
		avurl = fmt.Sprintf(`https://%s/_matrix/media/r0/download/%s`, serverName, StripMXCPrefix(mxc))
	}

	return avurl
}

func (c *Client) BuildAvatar(mxc string) string {

	avurl := StripMXCPrefix(mxc)

	if len(avurl) == 0 {
		return ""
	}

	serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

	avurl = fmt.Sprintf(`%s/_matrix/media/r0/thumbnail/%s?width=32&height=32&method=crop`, serverName, avurl)

	if c.Config.Mode == "production" {
		serverName = c.Config.Matrix.Server
		avurl = fmt.Sprintf(`https://%s/_matrix/media/r0/thumbnail/%s?width=32&height=32&method=crop`, serverName, StripMXCPrefix(mxc))
	}

	return avurl
}

func (c *Client) BuildProfileAvatar(mxc string) string {

	avurl := StripMXCPrefix(mxc)

	if len(avurl) == 0 {
		return ""
	}

	serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

	avurl = fmt.Sprintf(`%s/_matrix/media/r0/thumbnail/%s?width=100&height=100&method=crop`, serverName, avurl)

	if c.Config.Mode == "production" {
		serverName = c.Config.Matrix.Server
		avurl = fmt.Sprintf(`https://%s/_matrix/media/r0/thumbnail/%s?width=32&height=32&method=crop`, serverName, StripMXCPrefix(mxc))
	}

	return avurl
}

func (c *Client) BuildImage(mxc string) string {

	avurl := StripMXCPrefix(mxc)

	if len(avurl) == 0 {
		return ""
	}

	serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

	avurl = fmt.Sprintf(`%s/_matrix/media/r0/thumbnail/%s?width=800&height=600&method=scale`, serverName, avurl)

	if c.Config.Mode == "production" {
		serverName = c.Config.Matrix.Server
		avurl = fmt.Sprintf(`https://%s/_matrix/media/r0/thumbnail/%s?width=800&height=600&method=scale`, serverName, StripMXCPrefix(mxc))
	}

	return avurl
}

func RejectUsername(username string) bool {
	usernames := []string{
		"admin",
		"matrix",
	}

	exists := false

	for _, x := range usernames {
		if x == username {
			return true
		}
	}

	return exists
}

func GetLocalPart(s string) string {
	s = s[1:]
	x := strings.Split(s, ":")
	return x[0]
}
func (c *Client) GetLocalPartPath(s string, profile bool) string {
	s = s[1:]
	x := strings.Split(s, ":")

	g := strings.Split(x[0], "_")

	if profile {
		g = g[1:]
	}

	if !strings.Contains(x[1], c.Config.Client.Domain) {
		g[0] = fmt.Sprintf(`%s:%s`, g[0], x[1])
	}

	return strings.Join(g, "/")
}
