package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/unrolled/secure"
)

func (c *Client) SendSignupVerificationEmail(email string) error {

	token := RandomString(73)

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 7*time.Second)
	err := c.AddEmailVerification(ctx, email, token)
	if err != nil {
		log.Println(err)
	}

	link := fmt.Sprintf(`%s/signup?code=%s`, c.Config.Client.Domain, token)

	from := fmt.Sprintf(`support@%s`, c.Config.Client.Domain)
	password := c.Config.Email.Password

	to := []string{email}

	body := fmt.Sprintf(`
		We're sending you this email to make sure that this address actually belongs to you. Follow the link below to create your new Hummingbard account.<br><br>
		%s
	`, link)

	log.Println("waht is body", body)

	message := []byte("From:" + "support@hummingbard.com" + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: Create Your Hummingbard Account\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", password, password, c.Config.Email.Server)

	ad := fmt.Sprintf(`%s:%s`, c.Config.Email.Server, c.Config.Email.Port)

	err = smtp.SendMail(ad, auth, from, to, message)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *Client) SendEmailUpdateVerificationEmail(email string) error {

	token := RandomString(73)

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 7*time.Second)
	err := c.AddEmailVerification(ctx, email, token)
	if err != nil {
		log.Println(err)
	}

	link := fmt.Sprintf(`%s/verify?code=%s`, c.Config.Client.Domain, token)

	from := fmt.Sprintf(`support@%s`, c.Config.Client.Domain)
	password := c.Config.Email.Password

	to := []string{email}

	body := fmt.Sprintf(`
		You've added this email address to your Hummingbard account. Follow the link below to verify this email.<br><br>
		%s
	`, link)

	log.Println("waht is body", body)

	message := []byte("From:" + "support@hummingbard.com" + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: Verify Email Address Update\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", password, password, c.Config.Email.Server)

	ad := fmt.Sprintf(`%s:%s`, c.Config.Email.Server, c.Config.Email.Port)

	err = smtp.SendMail(ad, auth, from, to, message)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *Client) SendPasswordResetEmail(email string) error {

	token := RandomString(73)

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 7*time.Second)
	err := c.AddPasswordReset(ctx, email, token)
	if err != nil {
		log.Println(err)
	}

	link := fmt.Sprintf(`%s/password/reset?code=%s`, c.Config.Client.Domain, token)

	from := fmt.Sprintf(`support@%s`, c.Config.Client.Domain)
	password := c.Config.Email.Password

	to := []string{email}

	body := fmt.Sprintf(`
		You requested a password reset. Follow the link below to reset your password.<br><br>
		%s<br><br>
		If you did not request a password reset, you can safely ignore this email.
	`, link)

	log.Println("waht is body", body)

	message := []byte("From:" + "support@hummingbard.com" + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: Reset Your Hummingbard Password\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", password, password, c.Config.Email.Server)

	ad := fmt.Sprintf(`%s:%s`, c.Config.Email.Server, c.Config.Email.Port)

	err = smtp.SendMail(ad, auth, from, to, message)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *Client) SendInvitationEmail(email string, user *User) error {

	token := RandomString(73)

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 7*time.Second)
	err := c.AddEmailVerification(ctx, email, token)
	if err != nil {
		log.Println(err)
	}
	err = c.AddInvite(ctx, user.UserID, email, token)
	if err != nil {
		log.Println(err)
	}

	link := fmt.Sprintf(`%s/signup?code=%s`, c.Config.Client.Domain, token)

	from := fmt.Sprintf(`support@%s`, c.Config.Client.Domain)
	password := c.Config.Email.Password

	to := []string{email}

	username := user.UserID

	if !user.Federated {
		username = GetLocalPart(user.UserID)
		username = "@" + username
	}

	body := fmt.Sprintf(`
		Your friend <a href="https://%s/%s">%s</a> sent you an invite to join Hummingbard. Follow the link below to create your Hummingbard account.<br><br>
		%s
	`, c.Config.Client.Domain, username, username, link)

	log.Println("waht is body", body)

	message := []byte("From:" + "support@hummingbard.com" + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: You've been invited to join Hummingbard!\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", password, password, c.Config.Email.Server)

	ad := fmt.Sprintf(`%s:%s`, c.Config.Email.Server, c.Config.Email.Port)

	err = smtp.SendMail(ad, auth, from, to, message)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *Client) VerifyEmailUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		us := LoggedInUser(r)

		query := r.URL.Query()
		code := query.Get("code")

		if len(code) == 0 {
			c.NotFound(w, r)
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

		err = c.VerifyEmail(ctx, email)
		if err != nil {
			log.Println(err)
		} else {
			username := GetLocalPart(us.UserID)
			err = c.UnsafeUpdateEmail(ctx, username, email)
			if err != nil {
				log.Println(err)
			}
		}

		type page struct {
			BasePage
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce
		t.LoggedInUser = us

		c.Templates.ExecuteTemplate(w, "email-update-success", t)
	}
}
