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

type WellKnownServer struct {
	ServerName string `json:"m.server"`
}

func (c *Client) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type page struct {
			BasePage
			UserExists                bool
			ServerDown                bool
			LoginError                bool
			LoginUsername             string
			LoginFederated            bool
			FederationDisabled        bool
			FederationDisabledWarning bool
			PasswordResetSuccess      bool
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
		}
		if s != nil {
			x := s.Flashes("login-error")
			if len(x) > 0 {
				t.LoginError = true
			}
			u := s.Flashes("login-username")
			if len(u) > 0 {
				t.LoginUsername = u[0].(string)
			}
			f := s.Flashes("login-federated")
			if len(f) > 0 {
				t.LoginFederated = f[0].(bool)
			}
			p := s.Flashes("password-reset-success")
			if len(p) > 0 {
				t.PasswordResetSuccess = true
			}
			s.Save(r, w)
		}

		query := r.URL.Query()
		account := query.Get("account")

		if account == "matrix" {
			t.LoginFederated = true
		}

		t.Nonce = nonce

		c.Templates.ExecuteTemplate(w, "login", t)
	}
}

//Log user in
func (c *Client) ValidateLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		username := r.FormValue("username")
		password := r.FormValue("password")
		federated := r.FormValue("federated") == "on"

		if username == "" || password == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		fu, us := c.IsFederated(username)
		//port is only for my dev environment, this needs to go, or i'm just
		//confused
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
			log.Println(wk)
			serverName = c.URLScheme(wk.ServerName)
			username = fmt.Sprintf(`%s:%s`, us.LocalPart, us.ServerName)
		}

		matrix, err := gomatrix.NewClient(serverName, "", "")
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		rl := &gomatrix.ReqLogin{
			Type:     "m.login.password",
			User:     username,
			Password: password,
		}

		resp, err := matrix.Login(rl)
		if err != nil {
			log.Println(err)

			s.AddFlash("Username or Password Wrong", "login-error")
			s.AddFlash(username, "login-username")
			s.AddFlash(federated, "login-federated")
			s.Save(r, w)

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		matrix.SetCredentials(resp.UserID, resp.AccessToken)

		//check if a room exists for this username with canonical room alis in
		//the format #@username:server.org
		un := fmt.Sprintf(`#@%s:%s`, username, c.Config.Client.Domain)
		if fu {
			un = fmt.Sprintf(`#%s:%s`, us.LocalPart, us.ServerName)
		}

		res, err := matrix.ResolveAlias(un)
		if err != nil {
			log.Println(err)
		}

		if res != nil {
			go func() {
				c.OperatorJoinRoom(string(res.RoomID))
			}()
		}

		// If user's room doesn't exist, we create it
		newUser := false
		if res == nil && fu {
			newUser = true

			go func() {

				u := username

				if fu {
					u = us.LocalPart
				}

				crr, err := matrix.CreateRoom(&gomatrix.ReqCreateRoom{
					Visibility:    "public",
					Preset:        "public_chat",
					RoomAliasName: fmt.Sprintf(`%s`, u),
					Name:          fmt.Sprintf(`%s's Timeline`, u),
					CreationContent: map[string]interface{}{
						"m.federate": true,
					},
					InitialState: []gomatrix.Event{gomatrix.Event{
						Type: "m.room.history_visibility",
						Content: map[string]interface{}{
							"history_visibility": "world_readable",
						},
					}, gomatrix.Event{
						Type: "m.room.guest_access",
						Content: map[string]interface{}{
							"guest_access": "can_join",
						},
					}},
				})
				if err != nil || crr == nil {
					log.Println(err)
				}

				if crr != nil {

					go func() {
						c.OperatorJoinRoom(string(crr.RoomID))

						_, err := matrix.SendStateEvent(crr.RoomID, "m.room.power_levels", "", map[string]interface{}{
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
								resp.UserID:          100,
								c.DefaultUser.UserID: 100,
							},
							"users_default": 0,
						})
						if err != nil {
							log.Println(err)
						}

					}()

					c.UpdateUserRoomID(r, crr.RoomID)
				}

			}()
		}

		pub := fmt.Sprintf(`#public:%s`, c.Config.Client.Domain)
		_, err = matrix.JoinRoom(pub, "", nil)
		if err != nil {
			log.Println(err)
		}

		rms, err := c.GetUserJoinedRooms(matrix)
		if err != nil {
			c.Error(w, r)
			return
		}

		profile, err := matrix.GetProfile(resp.UserID)
		if err != nil {
			log.Println(err)
		}

		prefs, err := matrix.GetAccountPreferences(resp.UserID)
		if err != nil {
			log.Println(err)
		}

		// store user session
		token := RandomString(64)

		u := User{
			AccessToken:       token,
			MatrixAccessToken: resp.AccessToken,
			DeviceID:          resp.DeviceID,
			HomeServer:        resp.HomeServer,
			UserID:            resp.UserID,
			JoinedRooms:       rms,
			WellKnown:         serverName,
			Federated:         fu,
		}

		if prefs != nil {
			u.Preferences = *prefs
		}

		if profile != nil {
			if profile.Displayname != nil && len(*profile.Displayname) > 0 {
				u.DisplayName = *profile.Displayname
			}
			if profile.AvatarURL != nil && len(*profile.AvatarURL) > 0 {
				u.AvatarURL = StripMXCPrefix(*profile.AvatarURL)
			}
		}

		if res != nil && res.RoomID != "" {
			u.RoomID = string(res.RoomID)
		}

		serialized, err := json.Marshal(u)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		err = c.Store.Set(token, resp.UserID, 0).Err()
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		err = c.Store.Set(resp.UserID, serialized, 0).Err()
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		s.Values["access_token"] = token
		s.Values["device_id"] = resp.DeviceID

		s.AddFlash("User logged in", "login-success")
		if newUser {
			s.AddFlash("Signed Up", "signed-up")
		}
		s.Save(r, w)

		//redirect to index
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//signup page
func (c *Client) Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//us := LoggedInUser(r)

		type page struct {
			BasePage
			UserExists           bool
			ServerDown           bool
			SignupError          bool
			Interactive          bool
			HomeServer           string
			RegistrationDisabled bool
			Email                string
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
		}
		if s != nil {
			x := s.Flashes("user-exists")
			if len(x) > 0 {
				t.UserExists = true
				s.Save(r, w)
			}
			y := s.Flashes("server-down")
			if len(y) > 0 {
				t.ServerDown = true
				s.Save(r, w)
			}
			i := s.Flashes("interactive")
			if len(i) > 0 {
				t.Interactive = true
				t.HomeServer = i[0].(string)
				s.Save(r, w)
			}
		}

		t.Nonce = nonce

		query := r.URL.Query()
		code := query.Get("code")

		if c.Config.Auth.DisableRegistration && len(code) == 0 {
			c.Templates.ExecuteTemplate(w, "signup-disabled", t)
			return
		}

		if c.Config.Auth.VerifyEmail && code == "" {
			c.Templates.ExecuteTemplate(w, "verify-email", t)
			return
		}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		email, valid, err := c.GetEmailVerificationToken(ctx, code)
		if err != nil || !valid {
			log.Println(err)
			c.VerificationCodeInvalid(w, r)
			return
		}
		t.Email = email

		c.Templates.ExecuteTemplate(w, "signup", t)
	}
}

func (c *Client) SignupDisabled() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//us := LoggedInUser(r)

		type page struct {
			BasePage
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce

		c.Templates.ExecuteTemplate(w, "signup-disabled", t)
	}
}

func (c *Client) VerificationCodeInvalid(w http.ResponseWriter, r *http.Request) {
	us := LoggedInUser(r)

	type page struct {
		BasePage
	}

	nonce := secure.CSPNonce(r.Context())

	t := &page{}

	t.Nonce = nonce
	t.LoggedInUser = us

	c.Templates.ExecuteTemplate(w, "invalid-verification-code", t)
}

func (c *Client) PasswordResetCodeInvalid(w http.ResponseWriter, r *http.Request) {
	//us := LoggedInUser(r)

	type page struct {
		BasePage
	}

	nonce := secure.CSPNonce(r.Context())

	t := &page{}

	t.Nonce = nonce

	c.Templates.ExecuteTemplate(w, "invalid-password-reset-code", t)
}

// Copied from Dendrite clientapi/routing/register.go
const (
	minPasswordLength = 8   // http://matrix.org/docs/spec/client_server/r0.2.0.html#password-based
	maxPasswordLength = 512 // https://github.com/matrix-org/synapse/blob/v0.20.0/synapse/rest/client/v2_alpha/register.py#L161
	maxUsernameLength = 254 // http://matrix.org/speculator/spec/HEAD/intro.html#user-identifiers TODO account for domain
	sessionIDLength   = 24
)

//sign user up
func (c *Client) ValidateSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		repeat := r.FormValue("repeat")

		if j := RejectUsername(username); j {
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		if username == "" || password == "" ||
			len(username) < 3 ||
			len(username) > maxUsernameLength ||
			len(password) < minPasswordLength ||
			len(password) > maxPasswordLength {
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		if password != repeat {
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		type Auth struct {
			Type    string `json:"type"`
			Session string `json:"session"`
			Mac     string `json:"mac"`
		}

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		if strings.Contains(username, ":") {
			_, us := c.IsFederated(username)

			wk, err := WellKnown(c.URLScheme(us.ServerName))
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}
			serverName = c.URLScheme(wk.ServerName)
			//get rid of the @ prefix
			username = us.LocalPart[1:]
		}

		matrix, err := gomatrix.NewClient(serverName, "", "")
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		//check if username s available
		av, err := matrix.RegisterAvailable(&gomatrix.ReqRegisterAvailable{
			Username: username,
		})
		if err != nil {
			log.Println(err)

			s.AddFlash("Server Down", "server-down")
			s.Save(r, w)

			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		if av == nil || !av.Available {

			s.AddFlash("User Exists", "user-exists")
			s.Save(r, w)

			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		//actually register the user
		mac, err := ConstructMac(&NewUser{
			Username: username,
			Password: password,
		}, c.Config.Auth.SharedSecret)
		if err != nil {
			log.Println(err)
		}

		matrix.Prefix = "/_matrix/client/api/v1"

		req := &gomatrix.ReqLegacyRegister{
			Username: username,
			Password: password,
			Type:     "org.matrix.login.shared_secret",
			Mac:      mac,
		}

		resp, inter, err := matrix.LegacyRegister(req)

		if err != nil || (resp == nil && inter == nil) {
			log.Println(err)
			s.AddFlash("Server Down", "server-down")
			s.Save(r, w)

			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}

		matrix.Prefix = "/_matrix/client/r0"

		// create the user's timeline room
		matrix.SetCredentials(resp.UserID, resp.AccessToken)

		//let them join #public
		go func() {

			pub := fmt.Sprintf(`#public:%s`, c.Config.Client.Domain)
			_, err := matrix.JoinRoom(pub, "", nil)
			if err != nil {
				log.Println(err)
			}
		}()

		go func() {
			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, 7*time.Second)
			err := c.AddNewUser(ctx, resp.UserID, email, resp.AccessToken)
			if err != nil {
				log.Println(err)
			}

			err = c.UnsafeAddEmail(ctx, username, email)
			if err != nil {
				log.Println(err)
			}
		}()

		crr, err := matrix.CreateRoom(&gomatrix.ReqCreateRoom{
			Visibility:    "public",
			Preset:        "public_chat",
			RoomAliasName: fmt.Sprintf(`@%s`, username),
			Name:          fmt.Sprintf(`%s`, username),
			Topic:         fmt.Sprintf(`This is @%s's hummingbard profile.`, username),
			CreationContent: map[string]interface{}{
				"m.federate": true,
			},
			InitialState: []gomatrix.Event{gomatrix.Event{
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
					"type": "profile",
				},
			}, gomatrix.Event{
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
						resp.UserID:          100,
						c.DefaultUser.UserID: 100,
					},
					"users_default": 0,
				},
			}},
		})
		if err != nil || crr == nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		go func() {
			c.OperatorJoinRoom(string(crr.RoomID))

			/*
				text, html := InitialMessage()

				npe := gomatrix.CreatePostEvent{
					RoomID:        crr.RoomID,
					Text:          text,
					FormattedText: html,
				}

				_, err = matrix.CreatePost(&npe)
				if err != nil {
					log.Println(err)
				}
			*/
		}()

		rms, err := c.GetUserJoinedRooms(matrix)
		if err != nil {
			c.Error(w, r)
			return
		}
		//store session
		token := RandomString(64)
		u := User{
			AccessToken:       token,
			MatrixAccessToken: resp.AccessToken,
			DeviceID:          resp.DeviceID,
			HomeServer:        resp.HomeServer,
			UserID:            resp.UserID,
			RefreshToken:      resp.RefreshToken,
			JoinedRooms:       rms,
			RoomID:            crr.RoomID,
		}

		serialized, err := json.Marshal(u)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		err = c.Store.Set(token, resp.UserID, 0).Err()
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		err = c.Store.Set(resp.UserID, serialized, 0).Err()
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		go func() {
			err = c.RefreshRoomsCache()
			if err != nil {
				log.Println(err)
			}
		}()

		s.Values["access_token"] = token
		s.Values["device_id"] = resp.DeviceID

		s.AddFlash("Signed Up", "signed-up")
		s.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	}
}

//log user out, kill session in redis
func (c *Client) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		token, ok := s.Values["access_token"].(string)
		if ok {

			userid, err := c.Store.Get(token).Result()
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}

			_, err = c.Store.Del(userid).Result()
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}

			_, err = c.Store.Del(token).Result()
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}

			s.Values["access_token"] = ""
			s.Options.MaxAge = -1
			err = s.Save(r, w)
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func (c *Client) VerifySignupEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		email := r.FormValue("email")

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		exists, err := c.DoesUserExist(ctx, email)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		if !exists && len(email) > 0 {
			go c.SendSignupVerificationEmail(email)
		}

		type page struct {
			BasePage
			UserExists           bool
			ServerDown           bool
			SignupError          bool
			Interactive          bool
			HomeServer           string
			RegistrationDisabled bool
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
		}
		if s != nil {
			x := s.Flashes("user-exists")
			if len(x) > 0 {
				t.UserExists = true
				s.Save(r, w)
			}
			y := s.Flashes("server-down")
			if len(y) > 0 {
				t.ServerDown = true
				s.Save(r, w)
			}
			i := s.Flashes("interactive")
			if len(i) > 0 {
				t.Interactive = true
				t.HomeServer = i[0].(string)
				s.Save(r, w)
			}
		}

		t.Nonce = nonce

		c.Templates.ExecuteTemplate(w, "verify-success", t)
	}
}

func (c *Client) UsernameAvailable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type payload struct {
			Username string `json:"username"`
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
			Available bool `json:"available"`
		}

		ff := Response{}

		av, err := c.Matrix.RegisterAvailable(&gomatrix.ReqRegisterAvailable{
			Username: pay.Username,
		})

		if err != nil {
			log.Println(err)
		}

		if av != nil && av.Available {
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

func (c *Client) PasswordReset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type page struct {
			BasePage
			Code  string
			Email string
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce

		query := r.URL.Query()
		code := query.Get("code")

		if len(code) > 0 {

			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, 7*time.Second)
			email, valid, err := c.GetPasswordResetToken(ctx, code)
			if err != nil || !valid {
				log.Println(err)
				c.PasswordResetCodeInvalid(w, r)
				return
			}

			t.Code = code
			t.Email = email

			c.Templates.ExecuteTemplate(w, "password-reset-confirm", t)
			return
		}

		c.Templates.ExecuteTemplate(w, "password-reset", t)
	}
}

func (c *Client) ValidatePasswordReset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		email := r.FormValue("email")

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		exists, err := c.DoesUserExist(ctx, email)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		if exists && len(email) > 0 {
			go c.SendPasswordResetEmail(email)
		}

		type page struct {
			BasePage
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce

		c.Templates.ExecuteTemplate(w, "password-reset-code-sent", t)
	}
}

func (c *Client) ConfirmPasswordReset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		email := r.FormValue("email")
		password := r.FormValue("password")

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		userID, _, err := c.GetUser(ctx, email)
		if err != nil && len(userID) == 0 {
			log.Println(err)
			c.Error(w, r)
			return
		}

		username := GetLocalPart(userID)

		err = c.UnsafePasswordReset(ctx, username, password)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		if err != nil {
			log.Println(err)
		}

		err = c.InvalidatePasswordResetCode(ctx, email)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		if err != nil {
			log.Println(err)
		}

		type page struct {
			BasePage
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		s.AddFlash("Password Reset", "password-reset-success")
		s.Save(r, w)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}
