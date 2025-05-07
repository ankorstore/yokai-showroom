package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/ankorstore/yokai-showroom/mcp-demo/internal/domain"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type GetGopherTool struct {
	service *domain.GopherService
}

func NewGetGopherTool(service *domain.GopherService) *GetGopherTool {
	return &GetGopherTool{
		service: service,
	}
}

func (t *GetGopherTool) Name() string {
	return "get-gopher"
}

func (t *GetGopherTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("retrieve one specific gopher by id"),
		mcp.WithString(
			"id",
			mcp.Required(),
			mcp.Description("id of the gopher to retrieve"),
		),
	}
}

func (t *GetGopherTool) Handle() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, span := trace.CtxTracer(ctx).Start(ctx, "GetGopherTool.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("GetGopherTool.Handle")

		id := 0
		idParam, ok := request.Params.Arguments["id"]
		if ok {
			var err error

			id, err = strconv.Atoi(idParam.(string))
			if err != nil {
				return nil, errors.New("gopher id must be a integer")
			}
		}

		gopher, err := t.service.Get(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("cannot get gopher: %w", err)
		}

		jsonGopher, err := json.Marshal(gopher)
		if err != nil {
			return nil, fmt.Errorf("cannot json marshall gopher: %w", err)
		}

		return mcp.NewToolResultText(string(jsonGopher)), nil
	}
}
