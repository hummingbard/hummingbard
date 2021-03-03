package client

import (
	"context"
	"encoding/json"
	"fmt"
	"hummingbard/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hummingbard/gomatrix"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	"github.com/robfig/cron/v3"

	cache "hummingbard/cache"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Client struct {
	Config      *config.Config
	Router      *chi.Mux
	HTTP        *http.Server
	Templates   *Template
	Sessions    *sessions.CookieStore
	Store       *redis.Client
	Matrix      *gomatrix.Client
	DefaultUser User
	Cache       *cache.Cache
	DB          *DB
	Cron        *cron.Cron
}

func (c *Client) Activate() {

	log.Println("started server")

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)

		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		if err := c.HTTP.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
			log.Printf("Shutdown by user")
		}
		close(idleConnsClosed)
	}()

	if err := c.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

func Start() {
	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	conf, err := config.Read()
	if err != nil {
		panic(err)
	}

	tmpl, err := NewTemplate()
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	redis := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	server := fmt.Sprintf(`http://%s:%d`, conf.Matrix.Server, conf.Matrix.Port)
	matrix, err := gomatrix.NewClient(server, "", "")
	if err != nil {
		panic(err)
	}

	//log into the default matrix account
	defUser := User{}

	username := conf.Client.Domain
	password := conf.Matrix.Password

	resp, err := matrix.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     username,
		Password: password,
	})

	if resp != nil {
		defUser.UserID = resp.UserID
		defUser.AccessToken = resp.AccessToken
		matrix.SetCredentials(resp.UserID, resp.AccessToken)
	}

	//default account doesn't exist yet, let's create it
	if err != nil {
		log.Println(err)

		av, err := matrix.RegisterAvailable(&gomatrix.ReqRegisterAvailable{
			Username: username,
		})
		if err != nil {
			log.Println(err)
		}

		log.Println(av)

		if av == nil || !av.Available {
			panic(err)
		}

		type Auth struct {
			Type    string
			Session string
		}
		resp, _, err := matrix.Register(&gomatrix.ReqRegister{
			Username: username,
			Password: password,
			Auth: Auth{
				Type: "m.login.dummy",
			},
		})
		log.Println(resp)

		if err != nil || resp == nil {
			panic(err)
		}
		defUser.UserID = resp.UserID
		defUser.AccessToken = resp.AccessToken
		matrix.SetCredentials(resp.UserID, resp.AccessToken)

	}

	//does default #public room exist?
	//create #public room
	un := fmt.Sprintf(`#public:%s`, conf.Matrix.Server)
	res, err := matrix.ResolveAlias(un)
	if err != nil || res == nil {
		log.Println(err)
		//no, let's create it
		crr, err := matrix.CreateRoom(&gomatrix.ReqCreateRoom{
			Preset:        "public_chat",
			Visibility:    "public",
			RoomAliasName: "public",
			Name:          "Public",
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
		pub := fmt.Sprintf(`#public:%s`, conf.Client.Domain)
		_, err = matrix.JoinRoom(pub, "", nil)
		if err != nil {
			log.Println(err)
		}

	}

	sess := NewSession(conf.Client.SecureCookie)
	sess.Options.Domain = fmt.Sprintf(`.%s`, conf.Client.Domain)

	cache, err := cache.NewCache()
	if err != nil {
		panic(err)
	}

	cron := cron.New()

	c := &Client{
		DB:     db,
		Config: conf,
		Matrix: matrix,
		HTTP: &http.Server{
			ReadTimeout:  21 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  120 * time.Second,
			Addr:         conf.Client.Port,
			Handler:      router,
		},
		Router:      router,
		Templates:   tmpl,
		Sessions:    sess,
		Store:       redis,
		Cache:       cache,
		DefaultUser: defUser,
		Cron:        cron,
	}

	c.Middleware()
	c.Routes()

	c.Build()

	//let's cache bigroom items from json dump
	bigrooms, err := os.Open("bigrooms.json")
	if err != nil {
		panic(err)
	}
	defer bigrooms.Close()

	b, err := ioutil.ReadAll(bigrooms)
	if err != nil {
		panic(err)
	}

	var rpub *gomatrix.RespPublicRooms

	json.Unmarshal(b, &rpub)

	for i, _ := range rpub.Chunk {
		if rpub.Chunk[i].CanonicalAlias != "" {
			c.Cache.LargeRooms.Set(rpub.Chunk[i].CanonicalAlias[1:], true, 1)
		}
	}

	go c.RefreshRoomsCache()

	go c.Cron.AddFunc("*/15 * * * *", c.RefreshCache)
	go c.Cron.Start()

	c.Activate()
}
