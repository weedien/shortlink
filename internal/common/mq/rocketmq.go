package mq

import (
	rmqclient "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"log/slog"
	"os"
	"time"
)

const (
	AppLinkTopic  = "app-short-link-topic"
	Endpoint      = "193.112.178.249:8081" // rocketmq-proxy
	ConsumerGroup = "app-short-link-group"
	AccessKey     = ""
	SecretKey     = ""
)

var (
	// maximum waiting time for receive func
	awaitDuration = time.Second * 5
	// maximum number of messages received at one time
	maxMessageNum int32 = 16
	// invisibleDuration should > 20s
	invisibleDuration = time.Second * 20
	// receive messages in a loop
)

func ConnectToRocketMQ() (rmqclient.Producer, rmqclient.SimpleConsumer, func()) {
	_ = os.Setenv("mq.consoleAppender.enabled", "true")
	rmqclient.ResetLogger()

	// In most case, you don't need to create many producers, singleton pattern is more recommended.
	producer, err := rmqclient.NewProducer(&rmqclient.Config{
		Endpoint: Endpoint,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    AccessKey,
			AccessSecret: SecretKey,
		},
	},
		rmqclient.WithTopics(AppLinkTopic),
	)

	if err != nil {
		panic(err)
	}
	err = producer.Start()
	if err != nil {
		slog.Error("Failed to start producer: ", "error", err)
	}

	// In most case, you don't need to create many consumers, singleton pattern is more recommended.
	simpleConsumer, err := rmqclient.NewSimpleConsumer(&rmqclient.Config{
		Endpoint:      Endpoint,
		ConsumerGroup: ConsumerGroup,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    AccessKey,
			AccessSecret: SecretKey,
		},
	},
		rmqclient.WithAwaitDuration(awaitDuration),
		rmqclient.WithSubscriptionExpressions(map[string]*rmqclient.FilterExpression{
			AppLinkTopic: rmqclient.SUB_ALL,
		}),
	)
	if err != nil {
		slog.Info("Failed to create simple consumer: %s\n", "error", err)
	}
	// start simpleConsumer
	err = simpleConsumer.Start()
	if err != nil {
		slog.Info("Failed to start simple consumer: %s\n", "error", err)
	}

	// graceful stop producer and simpleConsumer
	stopFn := func() {
		if err := producer.GracefulStop(); err != nil {
			slog.Error("Failed to graceful stop producer", "error", err)
		}
		if err := simpleConsumer.GracefulStop(); err != nil {
			slog.Error("Failed to graceful stop simple consumer", "error", err)
		}
	}

	return producer, simpleConsumer, stopFn
}
