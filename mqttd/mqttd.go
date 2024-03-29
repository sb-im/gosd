package mqttd

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sb-im/jsonrpc-lite"
	log "github.com/sirupsen/logrus"

	packets "github.com/eclipse/paho.golang/packets"
	paho "github.com/eclipse/paho.golang/paho"
)

type Mqtt struct {
	Connect *paho.Connect
	Config  *MqttdConfig
	Client  *paho.Client
	rdb     *redis.Client
	cache   chan MqttRpc
	I       <-chan MqttRpc
	O       chan<- MqttRpc
}

type MqttRpc struct {
	ID      string
	Payload []byte
}

func NewMqttd(broker string, rdb *redis.Client, i <-chan MqttRpc, o chan<- MqttRpc) *Mqtt {
	config, err := loadMqttConfigDefault()
	if err != nil {
		log.Error(err)
	}

	config.Broker = broker
	opt, err := url.Parse(broker)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Infof("%+v\n", config)

	password, _ := opt.User.Password()
	cache := make(chan MqttRpc, 128)

	// 2h
	//sessionExpiryInterval := uint32(7200)

	return &Mqtt{
		I:      i,
		O:      o,
		rdb:    rdb,
		cache:  cache,
		Config: config,
		Client: paho.NewClient(paho.ClientConfig{
			ClientID: fmt.Sprint(config.Client, config.ID),
			Router: paho.NewSingleHandlerRouter(func(p *paho.Publish) {
				id := strings.Split(p.Topic, "/")[1]

				// nodes/%s/status
				// nodes/%s/network
				// nodes/%s/rpc/recv
				// nodes/%s/msg/%s

				// tasks/%s/term

				switch strings.Split(p.Topic, "/")[2] {
				case "rpc":
					o <- MqttRpc{
						ID:      id,
						Payload: p.Payload,
					}
					rdb.Set(context.Background(), p.Topic, p.Payload, 0)
				default:
					rdb.Set(context.Background(), p.Topic, p.Payload, 0)
				}

			}),
		}),
		Connect: paho.ConnectFromPacketConnect(&packets.Connect{
			//WillProperties: &packets.Properties{},

			//WillFlag:    true,
			//WillMessage: raw,
			//WillRetain:  true,
			//WillTopic:   fmt.Sprintf(config.Status, config.ID),
			//WillQOS:     1,
			Password: []byte(password),
			Username: opt.User.Username(),
			// https://stackoverflow.com/questions/65314401/cannot-connect-to-mosquitto-2-0-with-paho-library
			PasswordFlag: true,
			UsernameFlag: true,
			ClientID:     fmt.Sprintf(config.Client, config.ID),
			//CleanStart: false,
			CleanStart: true,
			// interval 10s
			KeepAlive: 10,
			// TODO:
			Properties: &packets.Properties{
				// PayloadFormat indicates the format of the payload of the message
				// 0 is unspecified bytes
				// 1 is UTF8 encoded character data
				//PayloadFormat: 1,
				// MessageExpiry is the lifetime of the message in seconds
				//MessageExpiry *uint32
				//// ContentType is a UTF8 string describing the content of the message
				//// for example it could be a MIME type
				//ContentType string
				//// ResponseTopic is a UTF8 string indicating the topic name to which any
				//// response to this message should be sent
				//ResponseTopic string
				//// CorrelationData is binary data used to associate future response
				//// messages with the original request message
				//CorrelationData []byte
				//// SubscriptionIdentifier is an identifier of the subscription to which
				//// the Publish matched
				//SubscriptionIdentifier *uint32
				//// SessionExpiryInterval is the time in seconds after a client disconnects
				//// that the server should retain the session information (subscriptions etc)
				//SessionExpiryInterval: &sessionExpiryInterval,
				//// AssignedClientID is the server assigned client identifier in the case
				//// that a client connected without specifying a clientID the server
				//// generates one and returns it in the Connack
				//AssignedClientID string
				//// ServerKeepAlive allows the server to specify in the Connack packet
				//// the time in seconds to be used as the keep alive value
				//ServerKeepAlive *uint16
				//// AuthMethod is a UTF8 string containing the name of the authentication
				//// method to be used for extended authentication
				//AuthMethod string
				//// AuthData is binary data containing authentication data
				//AuthData []byte
				//// RequestProblemInfo is used by the Client to indicate to the server to
				//// include the Reason String and/or User Properties in case of failures
				//RequestProblemInfo *byte
				//// WillDelayInterval is the number of seconds the server waits after the
				//// point at which it would otherwise send the will message before sending
				//// it. The client reconnecting before that time expires causes the server
				//// to cancel sending the will
				//WillDelayInterval *uint32
				//// RequestResponseInfo is used by the Client to request the Server provide
				//// Response Information in the Connack
				//RequestResponseInfo *byte
				//// ResponseInfo is a UTF8 encoded string that can be used as the basis for
				//// createing a Response Topic. The way in which the Client creates a
				//// Response Topic from the Response Information is not defined. A common
				//// use of this is to pass a globally unique portion of the topic tree which
				//// is reserved for this Client for at least the lifetime of its Session. This
				//// often cannot just be a random name as both the requesting Client and the
				//// responding Client need to be authorized to use it. It is normal to use this
				//// as the root of a topic tree for a particular Client. For the Server to
				//// return this information, it normally needs to be correctly configured.
				//// Using this mechanism allows this configuration to be done once in the
				//// Server rather than in each Client
				//ResponseInfo string
				//// ServerReference is a UTF8 string indicating another server the client
				//// can use
				//ServerReference string
				//// ReasonString is a UTF8 string representing the reason associated with
				//// this response, intended to be human readable for diagnostic purposes
				//ReasonString string
				//// ReceiveMaximum is the maximum number of QOS1 & 2 messages allowed to be
				//// 'inflight' (not having received a PUBACK/PUBCOMP response for)
				//ReceiveMaximum *uint16
				//// TopicAliasMaximum is the highest value permitted as a Topic Alias
				//TopicAliasMaximum *uint16
				//// TopicAlias is used in place of the topic string to reduce the size of
				//// packets for repeated messages on a topic
				//TopicAlias *uint16
				//// MaximumQOS is the highest QOS level permitted for a Publish
				//MaximumQOS *byte
				//// RetainAvailable indicates whether the server supports messages with the
				//// retain flag set
				//RetainAvailable *byte
				//// User is a map of user provided properties
				//User map[string]string
				//// MaximumPacketSize allows the client or server to specify the maximum packet
				//// size in bytes that they support
				//MaximumPacketSize *uint32
				//// WildcardSubAvailable indicates whether wildcard subscriptions are permitted
				//WildcardSubAvailable *byte
				//// SubIDAvailable indicates whether subscription identifiers are supported
				//SubIDAvailable *byte
				//// SharedSubAvailable indicates whether shared subscriptions are supported
				//SharedSubAvailable *byte
			},
		}),
	}
}

func (t *Mqtt) Run(ctx context.Context) {
	//t.Client.SetDebugLogger(logger.New(os.Stdout, "[DEBUG]: ", logger.LstdFlags | logger.Lshortfile))
	//t.Client.SetErrorLogger(logger.New(os.Stdout, "[ERROR]: ", logger.LstdFlags | logger.Lshortfile))

	opt, err := url.Parse(t.Config.Broker)
	if err != nil {
		log.Error(err)
	}

	go func() {
		keyspace := "__keyspace@" + strconv.Itoa(t.rdb.Options().DB) + "__:%s"
		ch := t.rdb.PSubscribe(
			context.Background(),
			fmt.Sprintf(keyspace, "nodes/*"),
			fmt.Sprintf(keyspace, "tasks/*"),
		).Channel()

		for {
			select {
			case m := <-ch:
				log.Debugf("%s: message: %s\n", m.Channel, m.Payload)
				if m.Payload != "set" {
					continue
				}

				// topic:
				// plans/1/running
				// nodes/4/msg/weather
				topic := strings.Split(m.Channel, ":")[1]
				data := t.rdb.Get(context.Background(), topic)
				log.Trace(data)
				raw, err := data.Bytes()
				if err != nil {
					log.Error(err)
					continue
				}

				// prefix:
				// tasks
				// nodes
				prefix := strings.Split(topic, "/")[0]

				switch prefix {
				case "nodes":
					// nodes/+/rpc/send
					if ks := strings.Split(topic, "/"); len(ks) == 4 && ks[3] == "send" {
						res, err := t.Client.Publish(ctx, &paho.Publish{
							Payload: raw,
							Topic:   topic,
							QoS:     0,
							Retain:  false,
						})

						if err != nil {
							if res != nil {
								log.Errorf("%+v\n", res)
							}
							log.Error(err)
						}
					}
				case "tasks":
					if strings.Split(topic, "/")[2] == "term" {
						continue
					}
					res, err := t.Client.Publish(ctx, &paho.Publish{
						Payload: raw,
						Topic:   topic,
						QoS:     1,
						Retain:  true,
					})

					if err != nil {
						if res != nil {
							log.Errorf("%+v\n", res)
						}
						log.Error(err)
					}
				default:
					log.Warnf("ignore -- %s: %s\n", topic, raw)
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Info("MQTT Try Connect")
			if conn, err := net.Dial("tcp", opt.Hostname()+":"+opt.Port()); err != nil {
				log.Error(err)
			} else {
				log.Info("MQTT TCP Connected")
				t.Client.Conn = conn
				t.doRun(ctx)
				conn.Close()
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (t *Mqtt) doRun(parent context.Context) {
	//pinger := NewPingHandler(t.Client, fmt.Sprintf(t.Config.Network, t.Config.ID))
	//pinger.SetDebug(logger.New(os.Stdout, "[Pinger]: ", logger.LstdFlags | logger.Lshortfile))
	//t.Client.PingHandler = pinger

	ctx, cancel := context.WithCancel(parent)
	t.Client.OnServerDisconnect = func(p *paho.Disconnect) {
		log.Warn("OnServerDisconnect")
		log.Error(p)
		cancel()
	}

	t.Client.OnClientError = func(err error) {
		log.Warn("OnClientError")
		log.Error(err)
		cancel()
	}

	defer log.Info("MQTT Close")
	t.Client.Connect(ctx, t.Connect)
	log.Info("MQTT Connected")

	res, err := t.Client.Subscribe(ctx, &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			fmt.Sprintf(t.Config.Rpc.I, "+"): {
				QoS: 2,
				//RetainHandling    byte
				//NoLocal           bool
				//RetainAsPublished bool
			},
			fmt.Sprintf(t.Config.Status, "+"): {
				QoS: 1,
			},
			fmt.Sprintf(t.Config.Network, "+"): {
				QoS: 1,
			},
			fmt.Sprintf(t.Config.Gtran.Prefix, "+", "+"): {
				QoS: 1,
			},
			fmt.Sprintf("tasks/%s/term", "+"): {
				QoS:     1,
				NoLocal: true,
			},
		},
	})

	if err != nil {
		if res != nil {
			log.Errorf("%+v\n", res)
		}
		log.Error(err)
		return
	}

	for {
		select {
		case raw := <-t.cache:
			if err := t.send(ctx, raw); err != nil {
				t.cache <- raw
				return
			}
		case raw := <-t.I:
			if err := t.send(ctx, raw); err != nil {
				t.cache <- raw
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (t *Mqtt) send(ctx context.Context, packet MqttRpc) error {
	id := packet.ID
	raw := packet.Payload
	if rpc, err := jsonrpc.Parse(raw); err == nil && (rpc.Type == jsonrpc.TypeRequest || rpc.Type == jsonrpc.TypeNotify) {
		//fmt.Println("[REQ]: ", string(raw))
		// JSON-RPC Request Ignore

		// {"jsonrpc":"2.0","method":"test","params":[]}
		// {"jsonrpc":"2.0","id":"test.0-1553321035000","method":"test","params":[]}

		// {"jsonrpc":"2.0","method":"ncp_offline"}

		res, err := t.Client.Publish(ctx, &paho.Publish{
			Payload: raw,
			Topic:   fmt.Sprintf(t.Config.Rpc.O, id),
			QoS:     2,
		})

		if err != nil {
			if res != nil {
				log.Errorf("%+v\n", res)
			}
			log.Error(err)
			return err
		}
	} else {
		log.Warnf("[Drop]: %s", raw)
	}
	return nil
}
