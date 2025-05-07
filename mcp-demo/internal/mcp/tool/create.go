package tool

import (
	"context"
	"errors"
	"fmt"

	"github.com/ankorstore/yokai-showroom/mcp-demo/internal/domain"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type CreateGopherTool struct {
	service *domain.GopherService
}

func NewCreateGopherTool(service *domain.GopherService) *CreateGopherTool {
	return &CreateGopherTool{
		service: service,
	}
}

func (t *CreateGopherTool) Name() string {
	return "create-gopher"
}

func (t *CreateGopherTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("create one new gopher"),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("name of the new gopher to create"),
		),
		mcp.WithString(
			"job",
			mcp.Required(),
			mcp.Description("job of the new gopher to create"),
			mcp.Enum("backend", "frontend"),
		),
	}
}

func (t *CreateGopherTool) Handle() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, span := trace.CtxTracer(ctx).Start(ctx, "CreateGopherTool.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("CreateGopherTool.Handle")

		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("gopher name must be a string")
		}

		job, ok := request.Params.Arguments["job"].(string)
		if !ok {
			return nil, errors.New("gopher job must be a string")
		}

		id, err := t.service.Create(ctx, name, job)
		if err != nil {
			return nil, fmt.Errorf("cannot create gopher: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("gopher was created with id %d", id)), nil
	}
}
