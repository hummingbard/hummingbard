package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type Client struct {
	Domain       string `toml:"domain"`
	Port         string `toml:"port"`
	CookieName   string `toml:"cookie_name"`
	SecureCookie string `toml:"secure_cookie"`
}

type Matrix struct {
	Server            string `toml:"server"`
	Port              int    `toml:"port"`
	Password          string `toml:"password"`
	AnonymousPassword string `toml:"anonymous_password"`
}

type DB struct {
	Client                  string `toml:"client"`
	DendriteUserAPIAccounts string `toml:"dendrite_user_api_accounts"`
}

type Redis struct {
	Address  string `toml:"address"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type Spaces struct {
	Prefix string `toml:"prefix"`
}

type Auth struct {
	VerifyEmail            bool   `toml:"verify_email"`
	DisableRegistration    bool   `toml:"disable_registration"`
	DisableFederatedLogin  bool   `toml:"disable_federated_login"`
	DisableProfileCreation bool   `toml:"disable_profile_creation"`
	SharedSecret           string `toml:"shared_secret"`
}

type Privacy struct {
	DisablePublic bool `toml:"disable_public"`
}

type Email struct {
	Server   string `toml:"server"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type Config struct {
	Name       string  `toml:"name"`
	Mode       string  `toml:"mode"`
	Client     Client  `toml:"client"`
	Matrix     Matrix  `toml:"matrix"`
	DB         DB      `toml:"db"`
	Redis      Redis   `toml:"redis"`
	YoutubeKey string  `toml:"youtube_key"`
	Spaces     Spaces  `toml:"spaces"`
	Auth       Auth    `toml:"auth"`
	Privacy    Privacy `toml:"privacy"`
	Email      Email   `toml:"email"`
}

var conf Config

// Read reads the config file and returns the Values struct
func Read() (*Config, error) {
	file, err := os.Open("config.toml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if _, err := toml.Decode(string(b), &conf); err != nil {
		panic(err)
	}

	return &conf, err
}
