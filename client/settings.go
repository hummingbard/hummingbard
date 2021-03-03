package client

import (
	"net/http"

	"github.com/unrolled/secure"
)

func (c *Client) Settings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		us := LoggedInUser(r)

		t := TimelinePage{}

		nonce := secure.CSPNonce(r.Context())

		t.Nonce = nonce
		t.LoggedInUser = us

		c.Templates.ExecuteTemplate(w, "settings", t)
	}
}
