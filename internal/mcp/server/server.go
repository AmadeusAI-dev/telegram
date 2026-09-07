package server

import (
	"net/http"

	"github.com/AmadeusAI-dev/telegram/internal/mcp/tools/getmessages"
	"github.com/AmadeusAI-dev/telegram/internal/mcp/tools/sendmessage"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func New(sender sendmessage.Sender, repo getmessages.MessageRepo) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "telegram-mcp", Version: "1.0.0"}, nil)

	registerTools(server, sender, repo)

	return server
}

func registerTools(server *mcp.Server, sender sendmessage.Sender, repo getmessages.MessageRepo) {
	mcp.AddTool(server, sendmessage.Info(), sendmessage.Handle(sender))
	mcp.AddTool(server, getmessages.Info(), getmessages.Handle(repo))
}

func Run(ch chan<- error, url string, server *mcp.Server) {
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	go func() {
		err := http.ListenAndServe(url, handler)
		if err != nil {
			ch <- err
		}
	}()
}
