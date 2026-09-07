package getmessages

import (
	"context"
	"reflect"
	"testing"

	"github.com/AmadeusAI-dev/telegram/internal/client"
	"github.com/AmadeusAI-dev/telegram/internal/mcp/tools/testutil"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SpyRepo struct {
	calls int
}

func (s *SpyRepo) Get(ctx context.Context, username string, limit int) ([]client.Message, error) {
	s.calls++
	return []client.Message{
		{
			ID:       100,
			Text:     "Hi",
			Outgoing: false,
		},
		{
			ID:       200,
			Text:     "Hello, World!",
			Outgoing: true,
		},
	}, nil
}

func TestReturnsMessages(t *testing.T) {
	server, cs := testutil.NewTestMcp(t)
	repo := &SpyRepo{}
	mcp.AddTool(server, Info(), Handle(repo))
	want := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: `{"Messages":[{"ID":100,"Outgoing":false,"Text":"Hi"},{"ID":200,"Outgoing":true,"Text":"Hello, World!"}]}`},
		},
		StructuredContent: map[string]any{
			"Messages": []any{
				map[string]any{
					"ID":       100.0,
					"Text":     "Hi",
					"Outgoing": false,
				},
				map[string]any{
					"ID":       200.0,
					"Text":     "Hello, World!",
					"Outgoing": true,
				},
			},
		},
	}

	got := testutil.CallTool(t, cs, "get_messages", map[string]any{
		"username": "mr_TheKiryuKha",
		"limit":    2,
	})

	testutil.AssertCallToolResultsMatch(t, got, want)
	if repo.calls != 1 {
		t.Fatalf("want %d repo calls, got %d", 1, repo.calls)
	}
}

func TestGetMessagesInfo(t *testing.T) {
	want := &mcp.Tool{
		Name:        "get_messages",
		Description: "retrieves messages from telegram chat, based on the provided username and limit of messages",
	}

	got := Info()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("get_messages tool info mismatch. got: %v, want: %v", got, want)
	}
}
