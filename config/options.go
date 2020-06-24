package config

const (
	defaultDebug            = false
	defaultMqttURL          = "mqtt://admin:public@localhost:1883"
	defaultBaseURL          = "http://localhost/gosd"
	defaultRootURL          = "http://localhost"
	defaultBasePath         = ""
	defaultDatabaseURL      = "postgres://postgres:password@localhost/gosd?sslmode=disable"
	defaultDatabaseMaxConns = 20
	defaultDatabaseMinConns = 1
	defaultListenAddr       = "127.0.0.1:8000"
)

// Options contains configuration options.
type Options struct {
	debug            bool
	mqttURL          string
	baseURL          string
	rootURL          string
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
		mqttURL:          defaultMqttURL,
		baseURL:          defaultBaseURL,
		rootURL:          defaultRootURL,
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

func (o *Options) MqttURL() string {
	return o.mqttURL
}

func (o *Options) BaseURL() string {
	return o.baseURL
}

func (o *Options) RootURL() string {
	return o.rootURL
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
