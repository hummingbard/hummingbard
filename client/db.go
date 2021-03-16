package client

import (
	"context"
	"database/sql"
	"hummingbard/config"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	*sqlx.DB
}

// NewDB returns a new database instace
func NewDB(dbType string) (*DB, error) {

	c, err := config.Read()
	if err != nil {
		panic(err)
	}

	conn := c.DB.Client

	if dbType == "userapi_accounts" {
		conn = c.DB.DendriteUserAPIAccounts
	}

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	store := &DB{db}
	return store, nil
}

func Slugify(title string) string {
	sp := strings.Split(title, " ")
	jp := strings.Join(sp, "-")
	lp := strings.ToLower(jp)

	reg := regexp.MustCompile("[^a-zA-Z0-9-]+")
	slug := reg.ReplaceAllString(lp, "")

	return slug
}

func (c *Client) DoesSlugExist(ctx context.Context, roomPath, slug string) (bool, error) {
	var exists bool
	err := c.DB.QueryRow("select exists(select 1 from slug_to_event where room_path=$1 and slug=$2)", roomPath, slug).Scan(&exists)
	if err != nil || err == sql.ErrNoRows {
		log.Println(err)
		return true, err
	}
	return exists, nil
}

func (c *Client) GetSlugEventID(ctx context.Context, roomPath, slug string) (string, error) {
	var event string
	err := c.DB.QueryRow("select event_id from slug_to_event where room_path=$1 and slug=$2", roomPath, slug).Scan(&event)
	if err != nil || err == sql.ErrNoRows {
		log.Println(err)
		return "", err
	}
	return event, nil
}

func (c *Client) UpdateEventSlug(ctx context.Context, roomPath, slug, event string) (bool, error) {

	_, err := c.DB.Exec(`INSERT INTO slug_to_event(room_path, slug, event_id) VALUES($1, $2, $3)`, roomPath, slug, event)
	log.Println(err)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) AddRoom(ctx context.Context, userID, roomID, roomAlias, path string) error {

	_, err := c.DB.Exec(`INSERT INTO rooms(user_id, room_id, room_alias, room_path) VALUES($1, $2, $3, $4)`, userID, roomID, roomAlias, path)
	log.Println(err)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetUserRooms(ctx context.Context, userID string) ([]*OwnedRoom, error) {

	query := `
		select 
            rooms.room_id,
            rooms.room_alias
		FROM rooms 
		WHERE user_id=$1
        LIMIt 53
    `
	sargs := []interface{}{userID}

	rows, err := c.DB.Queryx(query, sargs...)
	if err != nil || err == sql.ErrNoRows {
		log.Println(err)
		return nil, err
	}

	rooms := []*OwnedRoom{}

	for rows.Next() {
		results := make(map[string]interface{})
		err = rows.MapScan(results)

		room := OwnedRoom{}

		id, ok := results["room_id"].(string)
		if ok {
			room.RoomID = id
		}

		alias, ok := results["room_alias"].(string)
		if ok {
			room.RoomAlias = alias
		}

		rooms = append(rooms, &room)

	}

	return rooms, nil
}

func (c *Client) GetAllRooms(ctx context.Context) (map[string]string, error) {

	query := `
		select 
            rooms.room_id,
            rooms.room_alias
		FROM rooms 
    `

	rows, err := c.DB.Queryx(query)
	if err != nil || err == sql.ErrNoRows {
		log.Println(err)
		return nil, err
	}

	rooms := map[string]string{}

	for rows.Next() {
		results := make(map[string]interface{})
		err = rows.MapScan(results)

		id, ok := results["room_id"].(string)
		alias, ok := results["room_alias"].(string)
		if ok {
			rooms[id] = alias
		}

	}

	return rooms, nil
}

func (c *Client) AddEmailVerification(ctx context.Context, email, token string) error {

	_, err := c.DB.Exec(`INSERT INTO email_verification(email, token) VALUES($1, $2)`, email, token)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AddInvite(ctx context.Context, user_id, email, token string) error {

	_, err := c.DB.Exec(`INSERT INTO invites(invited_by, invitee_email, token) VALUES($1, $2, $3)`, user_id, email, token)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetEmailVerificationToken(ctx context.Context, code string) (string, bool, error) {
	var email string
	var valid bool
	err := c.DB.QueryRow("select email, valid from email_verification where token=$1 and valid=true", code).Scan(&email, &valid)
	if err != nil || err == sql.ErrNoRows {
		log.Println(err)
		return "", false, err
	}
	return email, valid, nil
}

func (c *Client) AddNewUser(ctx context.Context, user_id, email, token string) error {

	_, err := c.DB.Exec(`INSERT INTO users(user_id, email, email_verified, access_token) VALUES($1, $2, true, $3)`, user_id, email, token)
	if err != nil {
		return err
	}

	_, err = c.DB.Exec(`UPDATE email_verification set valid=false where email=$1`, email)
	if err != nil {
		return err
	}
	_, err = c.DB.Exec(`UPDATE invites set valid=false where invitee_email=$1`, email)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) VerifyEmail(ctx context.Context, email string) error {

	_, err := c.DB.Exec(`UPDATE users set email_verified=true where email=$1`, email)
	if err != nil {
		return err
	}

	_, err = c.DB.Exec(`UPDATE email_verification set valid=false where email=$1`, email)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DoesUserExist(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := c.DB.QueryRow("select exists(select 1 from users where email=$1)", email).Scan(&exists)
	if err != nil || err == sql.ErrNoRows {
		log.Println(err)
		return true, err
	}
	return exists, nil
}

func (c *Client) DoesEmailExist(ctx context.Context, email, userID string) (bool, error) {
	var exists bool
	err := c.DB.QueryRow("select exists(select 1 from users where email=$1 and user_id!=$2)", email, userID).Scan(&exists)
	if err != nil || err == sql.ErrNoRows {
		log.Println(err)
		return true, err
	}
	return exists, nil
}

func (c *Client) AddPasswordReset(ctx context.Context, email, token string) error {

	_, err := c.DB.Exec(`INSERT INTO password_resets(email, token) VALUES($1, $2)`, email, token)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetPasswordResetToken(ctx context.Context, code string) (string, bool, error) {
	var email string
	var valid bool
	err := c.DB.QueryRow("select email, valid from password_resets where token=$1 and valid=true", code).Scan(&email, &valid)
	if err != nil || err == sql.ErrNoRows {
		log.Println(err)
		return "", false, err
	}
	return email, valid, nil
}

func (c *Client) InvalidatePasswordResetCode(ctx context.Context, email string) error {

	_, err := c.DB.Exec(`UPDATE password_resets set valid=false where email=$1`, email)
	log.Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetUser(ctx context.Context, email string) (string, string, error) {
	var user_id, token string
	err := c.DB.QueryRow("select user_id, access_token from users where email=$1", email).Scan(&user_id, &token)
	if err != nil || err == sql.ErrNoRows {
		log.Println(err)
		return "", "", err
	}
	return user_id, token, nil
}

func (c *Client) GetEmailFromUserID(ctx context.Context, id string) (string, bool, error) {
	var user_id string
	var verified bool
	err := c.DB.QueryRow("select email, email_verified from users where user_id=$1", id).Scan(&user_id, &verified)
	if err != nil || err == sql.ErrNoRows {
		log.Println(err)
		return "", false, err
	}
	return user_id, verified, nil
}

func (c *Client) UnsafePasswordReset(ctx context.Context, username, password string) error {

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return err
	}

	hash := string(hashBytes)

	_, err = c.UserAPIAccountsDB.Exec(`UPDATE account_accounts set password_hash=$1 where localpart=$2`, hash, username)
	log.Println(err)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UnsafeAddEmail(ctx context.Context, username, email string) error {

	_, err := c.UserAPIAccountsDB.Exec(`INSERT into account_threepid(threepid, medium, localpart) VALUES($1, $2, $3)`, email, "email", username)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateEmail(ctx context.Context, userId, email string) error {

	_, err := c.DB.Exec(`UPDATE users SET email=$1, email_verified=false where user_id=$2`, email, userId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UnsafeUpdateEmail(ctx context.Context, username, email string) error {

	_, err := c.UserAPIAccountsDB.Exec(`UPDATE account_threepid SET threepid=$1 where localpart=$2`, email, username)
	if err != nil {
		return err
	}

	return nil
}
