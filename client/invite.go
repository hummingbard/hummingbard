package client

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/unrolled/secure"
)

func (c *Client) Invite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		us := LoggedInUser(r)

		type page struct {
			BasePage
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce
		t.LoggedInUser = us

		c.Templates.ExecuteTemplate(w, "invite", t)
	}
}

func (c *Client) ValidateInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		us := LoggedInUser(r)

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
			go c.SendInvitationEmail(email, us)
		}

		type page struct {
			BasePage
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce
		t.LoggedInUser = us

		c.Templates.ExecuteTemplate(w, "invite-sent", t)
	}
}
