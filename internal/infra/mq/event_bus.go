package mq

import (
	"context"
	rmqclient "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/bytedance/sonic"
	"log/slog"
	"reflect"
	"shortlink/common/base_event"
	"sync"
)

type RocketMqBasedEventBus struct {
	listeners    map[reflect.Type]*listenerGroup
	mu           sync.RWMutex
	typeRegistry map[string]reflect.Type
	producer     rmqclient.Producer
	consumer     rmqclient.SimpleConsumer
	stopFn       func()
}

type listenerGroup struct {
	mu        sync.RWMutex
	listeners []base_event.AppEventListener
}

func NewRocketMqBasedEventBus(ctx context.Context) *RocketMqBasedEventBus {
	producer, consumer, stopFn := ConnectToRocketMQ()

	bus := &RocketMqBasedEventBus{
		listeners:    make(map[reflect.Type]*listenerGroup),
		typeRegistry: make(map[string]reflect.Type),
		producer:     producer,
		consumer:     consumer,
		stopFn:       stopFn,
	}

	go bus.startReceivingMessages(ctx)

	return bus
}

func (bus *RocketMqBasedEventBus) Close() {
	bus.stopFn()
}

func (bus *RocketMqBasedEventBus) startReceivingMessages(ctx context.Context) {
	for {
		mvs, err := bus.consumer.Receive(ctx, 16, 20)
		if err != nil {
			slog.Error("Failed to receive message from RocketMQ", "error", err)
			continue
		}

		for _, mv := range mvs {
			eventType := bus.getTypeFromTag(mv.GetTag())
			if eventType == nil {
				slog.Warn("Unrecognized event type", "tag", mv.GetTag())
				continue
			}

			group := bus.getListenerGroup(eventType)
			if group == nil {
				slog.Warn("No listeners for event type", "eventType", eventType)
				continue
			}

			group.dispatchEvent(ctx, mv.GetBody(), eventType)
		}
	}
}

func (bus *RocketMqBasedEventBus) getTypeFromTag(tag *string) reflect.Type {
	bus.mu.RLock()
	defer bus.mu.RUnlock()
	return bus.typeRegistry[*tag]
}

func (bus *RocketMqBasedEventBus) getListenerGroup(eventType reflect.Type) *listenerGroup {
	bus.mu.RLock()
	defer bus.mu.RUnlock()
	return bus.listeners[eventType]
}

func (group *listenerGroup) dispatchEvent(ctx context.Context, data []byte, eventType reflect.Type) {
	group.mu.RLock()
	defer group.mu.RUnlock()

	for _, listener := range group.listeners {
		go func(listener base_event.AppEventListener) {
			eventPtr := reflect.New(eventType).Interface()
			if err := sonic.Unmarshal(data, eventPtr); err != nil {
				slog.Error("Failed to unmarshal message", "error", err)
				return
			}
			if err := listener.Process(ctx, eventPtr.(base_event.AppEvent)); err != nil {
				slog.Error("Failed to process message", "error", err)
			}
		}(listener)
	}
}

func (bus *RocketMqBasedEventBus) Publish(ctx context.Context, event base_event.AppEvent) {
	marshal, err := sonic.Marshal(event)
	if err != nil {
		slog.Error("Failed to marshal event", "error", err)
		return
	}

	msg := &rmqclient.Message{
		Topic: "app-short-link",
		Body:  marshal,
	}
	eventType := reflect.TypeOf(event)
	msg.SetTag(eventType.Name())

	resp, err := bus.producer.Send(ctx, msg)
	if err != nil {
		slog.Error("Failed to send message to RocketMQ", "error", err, "msg", marshal)
		return
	}
	slog.Info("Sent message to RocketMQ", "resp", resp)
}

func (bus *RocketMqBasedEventBus) Subscribe(eventType reflect.Type, listener base_event.AppEventListener) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	group, exists := bus.listeners[eventType]
	if !exists {
		group = &listenerGroup{}
		bus.listeners[eventType] = group
		bus.typeRegistry[eventType.Name()] = eventType
	}

	group.addListener(listener)
}

func (group *listenerGroup) addListener(listener base_event.AppEventListener) {
	group.mu.Lock()
	defer group.mu.Unlock()
	group.listeners = append(group.listeners, listener)
}