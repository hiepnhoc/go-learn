package config

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config App config struct
type Config struct {
	ServiceName    string
	Server         Server
	Postgres       Postgres
	DaprComponents DaprComponents
	Cookie         Cookie
	Session        Session
	Metrics        Metrics
	Logger         Logger
	Jaeger         Jaeger
}

// Server config struct
type Server struct {
	AppVersion        string
	Port              string
	PprofPort         string
	Mode              string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
	MaxConnectionIdle time.Duration
	Timeout           time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// Postgres config
type Postgres struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

// DaprComponents config struct
type DaprComponents struct {
	StateStore StateStore
	PubSub     PubSub
}

// StateStore config
type StateStore struct {
	ComponentName string
}

// PubSub config
type PubSub struct {
	ComponentName string
}

// Cookie config
type Cookie struct {
	Name     string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

// Session config
type Session struct {
	Prefix string
	Name   string
	Expire int
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

// Jaeger config
type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

// LoadConfig file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvPrefix(viper.GetString("ENV"))
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// ParseConfig file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

func GetConfig(configPath string) (*Config, error) {

	cfgFile, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
