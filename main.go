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

	listenSocket, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("error opening listen socket: %v", err)
	}

	m := lmqtt.New(
		lmqtt.WithTCPListener(listenSocket),
		lmqtt.WithConfig(mqttConfig),
		lmqtt.WithHook(lmqtt.Hooks{
			OnAccept: func(ctx context.Context, conn net.Conn) bool {
				log.Printf("\nOnAccept: conn = [%+v]", conn)
				return true
			},
			OnStop: func(ctx context.Context) {
				log.Printf("\nOnStop")
			},
			OnSubscribe: func(ctx context.Context, client lmqtt.Client, req *lmqtt.SubscribeRequest) error {
				log.Printf("\nOnSubscribe: client=[%+v] req=[%+v]", client, req)
				return nil
			},
			OnSubscribed: func(ctx context.Context, client lmqtt.Client, subscription *entities.Subscription) {
				log.Printf("\nOnSubscribed: client=[%+v] subscrioption=[%+v]", client, subscription)
			},
			OnUnsubscribe: func(ctx context.Context, client lmqtt.Client, req *lmqtt.UnsubscribeRequest) error {
				log.Printf("\nOnUnsubscribe: client=[%+v] req=[%+v]", client, req)
				return nil
			},
			OnUnsubscribed: func(ctx context.Context, client lmqtt.Client, topicName string) {
				log.Printf("\nOnUnsubscribed: client=[%+v] topicName=[%s]", client, topicName)
			},
			OnMsgArrived: func(ctx context.Context, client lmqtt.Client, req *lmqtt.MsgArrivedRequest) error {
				log.Printf("\nOnMsgArrived: client=[%+v] req=[%+v]", client, req)
				return nil
			},
			OnBasicAuth: func(ctx context.Context, client lmqtt.Client, req *lmqtt.ConnectRequest) error {
				log.Printf("\nOnBasicAuth: client=[%+v] req=[%+v]", client, req)
				return nil
			},

			OnEnhancedAuth: func(ctx context.Context, client lmqtt.Client, req *lmqtt.ConnectRequest) (*lmqtt.EnhancedAuthResponse, error) {
				log.Printf("\nOnEnhancedAuth: client=[%+v] req=[%+v]", client, req)
				return &lmqtt.EnhancedAuthResponse{
					Continue: false,
					OnAuth: func(ctx context.Context, client lmqtt.Client, req *lmqtt.AuthRequest) (*lmqtt.AuthResponse, error) {
						log.Printf("\nOnEnhancedAuth -> OnAuth client=[%+v] req=[%+v]", client, req)
						return &lmqtt.AuthResponse{
							Continue: false,
							AuthData: []byte{},
						}, err
					},
					AuthData: []byte{},
				}, nil
			},

			OnReAuth: func(ctx context.Context, client lmqtt.Client, auth *packets.Auth) (*lmqtt.AuthResponse, error) {
				log.Printf("\nOnReAuth: client=[%+v] auth=[%+v]", client, auth)

				// what would be the best response for debug?
				return &lmqtt.AuthResponse{
					Continue: false,
					AuthData: []byte{},
				}, err
			},
			OnConnected: func(ctx context.Context, client lmqtt.Client) {
				log.Printf("\nOnConnected: client=[%+v]", client)
			},
			OnSessionCreated: func(ctx context.Context, client lmqtt.Client) {
				log.Printf("\nOnSessionCreated: client=[%+v]", client)
			},
			OnSessionResumed: func(ctx context.Context, client lmqtt.Client) {
				log.Printf("\nOnSessionResumed: client=[%+v]", client)
			},
			OnSessionTerminated: func(ctx context.Context, clientID string, reason lmqtt.SessionTerminatedReason) {
				log.Printf("\nOnSessionTerminated: clientID=[%s] reason=[%+v]", clientID, reason)
			},
			OnDelivered: func(ctx context.Context, client lmqtt.Client, msg *entities.Message) {
				log.Printf("\nOnDelivered: client=[%+v] msg=[%+v]", client, msg)
			},
			OnClosed: func(ctx context.Context, client lmqtt.Client, err error) {
				log.Printf("\nOnClosed: client=[%+v] err=[%+v]", client, err)
			},
			OnMsgDropped: func(ctx context.Context, clientID string, msg *entities.Message, err error) {
				log.Printf("\nOnMsgDropped: clientID=[%s] msg=[%+v] err=[%+v]", clientID, msg, err)
			},
			OnWillPublish: func(ctx context.Context, clientID string, req *lmqtt.WillMsgRequest) {
				log.Printf("\nOnWillPublish: clientID=[%s] req=[%+v]", clientID, req)
			},
			OnWillPublished: func(ctx context.Context, clientID string, msg *entities.Message) {
				log.Printf("\nOnWillPublished: clientID=[%s] msg=[%+v]", clientID, msg)
			},
			OnPublish: func(ctx context.Context, client lmqtt.Client, msg *entities.Message) bool {
				log.Printf("\nOnPublish: client=[%+v] msg=[%+v]", client, msg)
				return true
			},
		}))

	err = m.Run()
	if err != nil {
		log.Printf("error running broker: %v", err)
	}
}
