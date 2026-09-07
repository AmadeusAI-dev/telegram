package getmessages

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/AmadeusAI-dev/telegram/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Info() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_messages",
		Description: "retrieves messages from telegram chat, based on the provided username and limit of messages",
	}
}

type GetMessagesInput struct {
	Username string `json:"username" jsonschema:"username of the user from the chat with whom you want to return messages"`
	Limit    int    `json:"limit" jsonschema:"amount of messages to retrieve"`
}

type GetMessagesOutput struct {
	Messages []client.Message
}

type MessageRepo interface {
	Get(context.Context, string, int) ([]client.Message, error)
}

var (
	UserNotFound        = errors.New("user not found")
	InternalServerError = errors.New("internal error(LOL)")
)

func Handle(repo MessageRepo) mcp.ToolHandlerFor[GetMessagesInput, GetMessagesOutput] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input GetMessagesInput) (
		*mcp.CallToolResult,
		GetMessagesOutput,
		error,
	) {
		messages, err := repo.Get(ctx, input.Username, input.Limit)

		switch {
		case errors.Is(err, client.ErrResolveUserPeer):
			err = UserNotFound

		case err != nil:
			slog.Error("failed to get messages", "err", err)
			err = fmt.Errorf("%w:%v", InternalServerError, err)
		}

		return nil, GetMessagesOutput{
			Messages: messages,
		}, err
	}
}
