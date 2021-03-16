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

func (c *Client) UpdatePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			OldPassword    string `json:"old_password"`
			NewPassword    string `json:"new_password"`
			RepeatPassword string `json:"repeat_password"`
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
			Error   bool `json:"error"`
		}

		res := Response{}

		if (len(pay.OldPassword) < 8 ||
			len(pay.NewPassword) < 8 ||
			len(pay.RepeatPassword) < 8) ||
			(pay.NewPassword != pay.RepeatPassword) {

			res.Error = true
			js, err := json.Marshal(res)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		username := GetLocalPart(user.UserID)

		rl := &gomatrix.ReqLogin{
			Type:     "m.login.password",
			User:     username,
			Password: pay.OldPassword,
		}

		_, err = matrix.Login(rl)
		if err != nil {
			log.Println(err)
			res.Error = true
			js, err := json.Marshal(res)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}

		_, err = matrix.UpdatePassword(map[string]interface{}{
			"new_password": pay.NewPassword,
			"auth": map[string]interface{}{
				"user":     username,
				"password": pay.OldPassword,
				"type":     "m.login.password",
			},
		})

		if err != nil {
			log.Println(err)
			res.Error = true
		}

		res.Updated = true

		js, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

/*
func (c *Client) NewMFA() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type Response struct {
			Secret string `json:"secret"`
			Image  []byte `json:"image"`
		}

		ff := Response{}

		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      c.Config.Name,
			AccountName: user.UserID,
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var buf bytes.Buffer
		img, err := key.Image(400, 400)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		png.Encode(&buf, img)

		ff.Secret = key.Secret()
		ff.Image = buf.Bytes()

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
*/
