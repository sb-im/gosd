package config

import (
	"errors"
	"fmt"
	url_parser "net/url"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	opts *Options
}

func NewParser() *Parser {
	return &Parser{
		opts: NewOptions(),
	}
}

func (p *Parser) ParseEnvironmentVariables() (*Options, error) {
	err := p.parseLines(os.Environ())
	if err != nil {
		return nil, err
	}
	return p.opts, nil
}

func (p *Parser) parseLines(lines []string) (err error) {
	var port string
	for _, line := range lines {
		fields := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		switch key {
		case "DEBUG":
			p.opts.debug = parseBool(value, defaultDebug)
		case "MQTT_CLIENT_ID":
			p.opts.mqttClientID = parseString(value, defaultMqttClientID)
		case "MQTT_URL":
			p.opts.mqttURL = parseString(value, defaultMqttURL)
		case "REDIS_URL":
			url, err := url_parser.Parse(value)
			if err == nil && url.Path != "" && url.Path != "/0"{
				old := value
				url.Path = "0"
				value = url.String()
				fmt.Println("ERROR: Only use redis 0, Automatic conversion ===")
				fmt.Printf("%s => %s\n", old , url.String())
			}
			p.opts.redisURL = parseString(value, defaultRedisURL)
		case "BASE_URL":
			p.opts.baseURL, p.opts.rootURL, p.opts.basePath, err = parseBaseURL(value)
			if err != nil {
				return err
			}
		case "PORT":
			port = value
		case "LISTEN_ADDR":
			p.opts.listenAddr = parseString(value, defaultListenAddr)
		case "DATABASE_URL":
			p.opts.databaseURL = parseString(value, defaultDatabaseURL)
		case "DATABASE_MAX_CONNS":
			p.opts.databaseMaxConns = parseInt(value, defaultDatabaseMaxConns)
		case "DATABASE_MIN_CONNS":
			p.opts.databaseMinConns = parseInt(value, defaultDatabaseMinConns)
		case "LOG_FILE":
			p.opts.logFile = parseString(value, defaultLogFile)
		case "LOG_LEVEL":
			p.opts.logLevel = parseString(value, defaultLogLevel)
		case "OAUTH_CLIENT_ID":
			p.opts.oauthID = parseString(value, defaultOauthID)
		case "OAUTH_CLIENT_SECRET":
			p.opts.oauthSecret = parseString(value, defaultOauthSecret)
		case "LUA_FILE":
			p.opts.luaFile = parseString(value, defaultLuaFile)
		}
	}

	if port != "" {
		p.opts.listenAddr = ":" + port
	}
	return nil
}

func parseBaseURL(value string) (string, string, string, error) {
	if value == "" {
		return defaultBaseURL, defaultRootURL, "", nil
	}

	if value[len(value)-1:] == "/" {
		value = value[:len(value)-1]
	}

	url, err := url_parser.Parse(value)
	if err != nil {
		return "", "", "", fmt.Errorf("Invalid BASE_URL: %v", err)
	}

	scheme := strings.ToLower(url.Scheme)
	if scheme != "https" && scheme != "http" {
		return "", "", "", errors.New("Invalid BASE_URL: scheme must be http or https")
	}

	basePath := url.Path
	url.Path = ""
	return value, url.String(), basePath, nil
}

func parseBool(value string, fallback bool) bool {
	if value == "" {
		return fallback
	}

	value = strings.ToLower(value)
	if value == "1" || value == "yes" || value == "true" || value == "on" {
		return true
	}

	return false
}

func parseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return v
}

func parseString(value string, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
