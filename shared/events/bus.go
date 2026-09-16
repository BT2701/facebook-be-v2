package events

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	ChannelCreate = "social.notifications"
	ChannelDelete = "social.notifications.delete"
)

type NotificationEvent struct {
	User     string `json:"user"`
	Receiver string `json:"receiver"`
	Post     string `json:"post"`
	Content  string `json:"content"`
	Action   int    `json:"action_n"`
}

type Bus struct {
	rdb *redis.Client
}

func New(addr string) *Bus {
	if addr == "" {
		return &Bus{}
	}
	return &Bus{
		rdb: redis.NewClient(&redis.Options{
			Addr:         addr,
			DialTimeout:  2 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		}),
	}
}

func (b *Bus) enabled() bool {
	return b != nil && b.rdb != nil
}

func (b *Bus) PublishCreate(ctx context.Context, event NotificationEvent) {
	b.publish(ctx, ChannelCreate, event)
}

func (b *Bus) PublishDelete(ctx context.Context, event NotificationEvent) {
	b.publish(ctx, ChannelDelete, event)
}

func (b *Bus) publish(ctx context.Context, channel string, event NotificationEvent) {
	if !b.enabled() {
		return
	}
	if event.Post == "" {
		event.Post = "0"
	}
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("events: encode %s: %v", channel, err)
		return
	}
	if err := b.rdb.Publish(ctx, channel, payload).Err(); err != nil {
		log.Printf("events: publish %s: %v", channel, err)
	}
}

func (b *Bus) Subscribe(ctx context.Context, onCreate, onDelete func(NotificationEvent)) {
	if !b.enabled() {
		return
	}

	pubsub := b.rdb.Subscribe(ctx, ChannelCreate, ChannelDelete)
	defer func() {
		_ = pubsub.Close()
	}()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			var event NotificationEvent
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Printf("events: decode %s: %v", msg.Channel, err)
				continue
			}
			switch msg.Channel {
			case ChannelCreate:
				if onCreate != nil {
					onCreate(event)
				}
			case ChannelDelete:
				if onDelete != nil {
					onDelete(event)
				}
			}
		}
	}
}
