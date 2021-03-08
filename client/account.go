package client

import (
	"encoding/json"
	"hummingbard/gomatrix"
	"log"
	"net/http"
)

func (c *Client) UpdatePreferences() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			*gomatrix.UserPreferences
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
		}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}
		res := Response{}

		err = matrix.UpdateAccountPreferences(user.UserID, pay)
		if err != nil {
			log.Println(err)
		} else {
			res.Updated = true
			c.RefreshPreferences(r)
		}

		js, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
func (c *Client) DeleteAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Password string `json:"password"`
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
			Deactivated bool `json:"deactivated"`
			Error       bool `json:"error"`
		}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		inter, _ := matrix.DeactivateInitiate()

		_, err = matrix.Deactivate(map[string]interface{}{
			"auth": map[string]interface{}{
				"user":     user.UserID,
				"password": pay.Password,
				"type":     "m.login.password",
				"session":  inter.Session,
			},
		})

		res := Response{}
		if err != nil {
			log.Println(err)
			res.Error = true
		} else {
			res.Deactivated = true
		}

		js, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
