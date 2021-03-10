package config

const (
	defaultDebug            = false
	defaultMqttClientID     = "cloud.0"
	defaultMqttURL          = "mqtt://admin:public@localhost:1883"
	defaultBaseURL          = "http://localhost/gosd"
	defaultRootURL          = "http://localhost"
	defaultRedisURL         = "redis://localhost:6379/0"
	defaultBasePath         = ""
	defaultDatabaseURL      = "postgres://postgres:password@localhost/gosd?sslmode=disable"
	defaultDatabaseMaxConns = 20
	defaultDatabaseMinConns = 1
	defaultListenAddr       = "127.0.0.1:8000"
)

// Options contains configuration options.
type Options struct {
	debug            bool
	mqttClientID     string
	mqttURL          string
	baseURL          string
	rootURL          string
	redisURL         string
	basePath         string
	databaseURL      string
	databaseMaxConns int
	databaseMinConns int
	listenAddr       string
}

// NewOptions returns Options with default values.
func NewOptions() *Options {
	return &Options{
		debug:            defaultDebug,
		mqttClientID:     defaultMqttClientID,
		mqttURL:          defaultMqttURL,
		baseURL:          defaultBaseURL,
		rootURL:          defaultRootURL,
		redisURL:         defaultRedisURL,
		basePath:         defaultBasePath,
		databaseURL:      defaultDatabaseURL,
		databaseMaxConns: defaultDatabaseMaxConns,
		databaseMinConns: defaultDatabaseMinConns,
		listenAddr:       defaultListenAddr,
	}
}

func (o *Options) HasDebugMode() bool {
	return o.debug
}

func (o *Options) MqttClientID() string {
	return o.mqttClientID
}

func (o *Options) MqttURL() string {
	return o.mqttURL
}

func (o *Options) BaseURL() string {
	return o.baseURL
}

func (o *Options) RootURL() string {
	return o.rootURL
}

func (o *Options) RedisURL() string {
	return o.redisURL
}

func (o *Options) DatabaseURL() string {
	return o.databaseURL
}

func (o *Options) DatabaseMaxConns() int {
	return o.databaseMaxConns
}

func (o *Options) DatabaseMinConns() int {
	return o.databaseMinConns
}

func (o *Options) ListenAddr() string {
	return o.listenAddr
}
