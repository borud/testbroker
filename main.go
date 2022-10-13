package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/lab5e/lmqtt/pkg/config"
	"github.com/lab5e/lmqtt/pkg/entities"
	"github.com/lab5e/lmqtt/pkg/lmqtt"
	"github.com/lab5e/lmqtt/pkg/packets"
	_ "github.com/lab5e/lmqtt/pkg/persistence"     // required for side-effects
	_ "github.com/lab5e/lmqtt/pkg/topicalias/fifo" // required for side-effects
)

func main() {
	listenAddr := ":1883"
	if len(os.Args) > 1 {
		listenAddr = os.Args[1]
	}

	mqttConfig := config.DefaultConfig()
	mqttConfig.MQTT.AllowZeroLenClientID = false
	mqttConfig.Logging = true
	mqttConfig.DumpPacket = true

	listenSocket, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("error opening listen socket: %v", err)
	}

	m := lmqtt.New(
		lmqtt.WithTCPListener(listenSocket),
		lmqtt.WithConfig(mqttConfig),
		lmqtt.WithHook(lmqtt.Hooks{
			OnAccept: func(ctx context.Context, conn net.Conn) bool {
				return true
			},
			OnStop: func(ctx context.Context) {},
			OnSubscribe: func(ctx context.Context, client lmqtt.Client, req *lmqtt.SubscribeRequest) error {
				return nil
			},
			OnSubscribed: func(ctx context.Context, client lmqtt.Client, subscription *entities.Subscription) {},
			OnUnsubscribe: func(ctx context.Context, client lmqtt.Client, req *lmqtt.UnsubscribeRequest) error {
				return nil
			},
			OnUnsubscribed: func(ctx context.Context, client lmqtt.Client, topicName string) {},
			OnMsgArrived: func(ctx context.Context, client lmqtt.Client, req *lmqtt.MsgArrivedRequest) error {
				return nil
			},
			OnBasicAuth: func(ctx context.Context, client lmqtt.Client, req *lmqtt.ConnectRequest) error {
				return nil
			},

			OnEnhancedAuth: func(ctx context.Context, client lmqtt.Client, req *lmqtt.ConnectRequest) (*lmqtt.EnhancedAuthResponse, error) {
				return &lmqtt.EnhancedAuthResponse{
					Continue: false,
					OnAuth: func(ctx context.Context, client lmqtt.Client, req *lmqtt.AuthRequest) (*lmqtt.AuthResponse, error) {
						return &lmqtt.AuthResponse{
							Continue: false,
							AuthData: []byte{},
						}, err
					},
					AuthData: []byte{},
				}, nil
			},

			OnReAuth: func(ctx context.Context, client lmqtt.Client, auth *packets.Auth) (*lmqtt.AuthResponse, error) {
				// what would be the best response for debug?
				return &lmqtt.AuthResponse{
					Continue: false,
					AuthData: []byte{},
				}, err
			},
			OnConnected:         func(ctx context.Context, client lmqtt.Client) {},
			OnSessionCreated:    func(ctx context.Context, client lmqtt.Client) {},
			OnSessionResumed:    func(ctx context.Context, client lmqtt.Client) {},
			OnSessionTerminated: func(ctx context.Context, clientID string, reason lmqtt.SessionTerminatedReason) {},

			OnDelivered:     func(ctx context.Context, client lmqtt.Client, msg *entities.Message) {},
			OnClosed:        func(ctx context.Context, client lmqtt.Client, err error) {},
			OnMsgDropped:    func(ctx context.Context, clientID string, msg *entities.Message, err error) {},
			OnWillPublish:   func(ctx context.Context, clientID string, req *lmqtt.WillMsgRequest) {},
			OnWillPublished: func(ctx context.Context, clientID string, msg *entities.Message) {},
			OnPublish: func(ctx context.Context, client lmqtt.Client, msg *entities.Message) bool {
				return true
			},
		}))

	err = m.Run()
	if err != nil {
		log.Printf("error running broker: %v", err)
	}
}
