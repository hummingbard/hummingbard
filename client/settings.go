package client

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/unrolled/secure"
)

func (c *Client) Settings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		us := LoggedInUser(r)

		type Page struct {
			BasePage
		}

		t := Page{}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		email, verified, err := c.GetEmailFromUserID(ctx, us.UserID)
		if err != nil {
			log.Println(err)
		}
		us.Email = email
		us.EmailVerified = verified

		nonce := secure.CSPNonce(r.Context())

		t.Nonce = nonce
		t.LoggedInUser = us

		c.Templates.ExecuteTemplate(w, "settings", t)
	}
}

func (c *Client) ChangeEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Email string `json:"email"`
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
			Updated bool `json:"updated"`
			Exists  bool `json:"exists"`
		}

		ff := Response{}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)

		exists, err := c.DoesEmailExist(ctx, pay.Email, user.UserID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}
		if exists {
			log.Println(err)
			ff.Exists = true
			js, err := json.Marshal(ff)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}

		err = c.UpdateEmail(ctx, user.UserID, pay.Email)
		if err != nil {
			log.Println(err)
		} else {
			go c.SendEmailUpdateVerificationEmail(pay.Email)
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

func (c *Client) ResendEmailVerification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Email string `json:"email"`
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
			Resent bool `json:"resent"`
			Exists bool `json:"exists"`
		}

		ff := Response{}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)

		exists, err := c.DoesEmailExist(ctx, pay.Email, user.UserID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}
		if exists {
			log.Println(err)
			ff.Exists = true
			js, err := json.Marshal(ff)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}

		err = c.UpdateEmail(ctx, user.UserID, pay.Email)
		if err != nil {
			log.Println(err)
		} else {
			go c.SendEmailUpdateVerificationEmail(pay.Email)
		}

		ff.Resent = true

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}
