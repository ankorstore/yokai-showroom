package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ankorstore/yokai-showroom/mcp-demo/internal/domain"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ListGophersTool struct {
	service *domain.GopherService
}

func NewListGophersTool(service *domain.GopherService) *ListGophersTool {
	return &ListGophersTool{
		service: service,
	}
}

func (t *ListGophersTool) Name() string {
	return "list-gophers"
}

func (t *ListGophersTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("list several gophers, optionally by job"),
		mcp.WithString(
			"job",
			mcp.DefaultString(""),
			mcp.Description("optional job of the gophers to list, empty value means all jobs"),
			mcp.Enum("", "frontend", "backend"),
		),
	}
}

func (t *ListGophersTool) Handle() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx = context.WithoutCancel(ctx)

		ctx, span := trace.CtxTracer(ctx).Start(ctx, "ListGophersTool.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("ListGophersTool.Handle")

		job := ""
		jobParam, ok := request.Params.Arguments["job"]
		if ok {
			job = jobParam.(string)
		}

		gophers, err := t.service.List(ctx, "", job)
		if err != nil {
			return nil, fmt.Errorf("cannot list gophers: %w", err)
		}

		jsonGophers, err := json.Marshal(gophers)
		if err != nil {
			return nil, fmt.Errorf("cannot json marshall gophers list: %w", err)
		}

		return mcp.NewToolResultText(string(jsonGophers)), nil
	}
}
