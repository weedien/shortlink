package config

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"os"
	"shortlink/pkg/toolkit"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigWithDefault struct {
	key          string
	defaultValue string
}

var (
	DSN                  = ConfigWithDefault{"DSN", "host=remote user=weedien password=031209 dbname=wespace search_path=link port=5432 sslmode=disable TimeZone=Asia/Shanghai"}
	RedisAddr            = ConfigWithDefault{"REDIS_ADDR", "localhost:6379"}
	RedisPassword        = ConfigWithDefault{"REDIS_PASSWORD", ""}
	RedisDB              = ConfigWithDefault{"REDIS_DB", "0"}
	EnableSharding       = ConfigWithDefault{"ENABLE_SHARDING", "false"}
	BaseRoutePrefix      = ConfigWithDefault{"BASE_ROUTE_PREFIX", "/api/short-link/v1"}
	Port                 = ConfigWithDefault{"PORT", "8080"}
	ShortLinkDomain      = ConfigWithDefault{"SHORT_LINK_DOMAIN", "http://localhost:8080"}
	EnableWhiteList      = ConfigWithDefault{"ENABLE_WHITE_LIST", "false"}
	DomainWhiteList      = ConfigWithDefault{"DOMAIN_WHITE_LIST", ""}
	DomainWhiteListNames = ConfigWithDefault{"DOMAIN_WHITE_LIST_NAMES", "掘金,知乎,简书,博客园,CSDN,开源中国,SegmentFault,思否,博客,博客园,博客园首页,博客首页,博客园博客"}
	DefaultFavicon       = ConfigWithDefault{"DEFAULT_FAVICON", "https://cdn.jsdelivr.net/gh/weedien/shortlink@main/static/favicon.ico"}
)

func (c ConfigWithDefault) String() string {
	return Default(c.key, c.defaultValue)
}

func (c ConfigWithDefault) Int() int {
	value := Default(c.key, c.defaultValue)
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Warn(fmt.Sprintf("Config key %s value %s is not a int", c.key, value))
		return 0
	}
	return i
}

func (c ConfigWithDefault) Bool() bool {
	value := Default(c.key, c.defaultValue)
	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Warn(fmt.Sprintf("Config key %s value %s is not a bool", c.key, value))
		return false
	}
	return b
}

func (c ConfigWithDefault) Array() []string {
	value := Default(c.key, c.defaultValue)
	// split by comma
	return toolkit.Split(value, ',')
}

// Config func to get env value
func Config(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}

// Default func to get env value with default value
func Default(key string, defaultValue string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	}
	return defaultValue
}

// DefaultInt func to get env value with default value
func DefaultInt(key string, defaultValue int) int {
	c := Config(key)
	if c == "" {
		return defaultValue
	}
	ci, err := strconv.Atoi(c)
	if err != nil {
		return ci
	}
	return defaultValue
}

func DefaultBool(key string, defaultValue bool) bool {
	c := Config(key)
	if c == "" {
		return defaultValue
	}
	cb, err := strconv.ParseBool(c)
	if err != nil {
		return cb
	}
	return defaultValue
}
