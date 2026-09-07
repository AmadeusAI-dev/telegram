package client

import (
	"context"
	"fmt"
	"slices"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

type Message struct {
	ID       int
	Text     string
	Outgoing bool
}

type MessageRepo struct {
	client *tg.Client
}

func NewMessageRepo(client *tg.Client) *MessageRepo {
	return &MessageRepo{client: client}
}

func (m *MessageRepo) Get(ctx context.Context, username string, limit int) ([]Message, error) {
	var msgs []Message

	p, err := peers.Options{}.Build(m.client).Resolve(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve user peer by username")
	}

	history, err := m.client.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		Peer:  p.InputPeer(),
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	messages, err := extractMessages(history)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {

		message, ok := message.(*tg.Message)
		if !ok {
			continue
		}

		msgs = append(msgs, Message{
			ID:       message.ID,
			Text:     message.Message,
			Outgoing: message.Out,
		})
	}

	slices.Reverse(msgs)

	return msgs, nil
}

func extractMessages(history tg.MessagesMessagesClass) ([]tg.MessageClass, error) {
	messages, ok := history.(*tg.MessagesMessagesSlice)
	if !ok {
		return nil, fmt.Errorf("failed to get MessagesMessages from MessagesMessagesClass")
	}

	return messages.Messages, nil
}
