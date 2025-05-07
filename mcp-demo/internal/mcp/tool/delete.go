package tool

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ankorstore/yokai-showroom/mcp-demo/internal/domain"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type DeleteGopherTool struct {
	service *domain.GopherService
}

func NewDeleteGopherTool(service *domain.GopherService) *DeleteGopherTool {
	return &DeleteGopherTool{
		service: service,
	}
}

func (t *DeleteGopherTool) Name() string {
	return "delete-gopher"
}

func (t *DeleteGopherTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("delete one specific gopher by id"),
		mcp.WithString(
			"id",
			mcp.Required(),
			mcp.Description("id of the gopher to delete"),
		),
	}
}

func (t *DeleteGopherTool) Handle() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, span := trace.CtxTracer(ctx).Start(ctx, "DeleteGopherTool.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("DeleteGopherTool.Handle")

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

		err = t.service.Delete(ctx, int(gopher.ID))
		if err != nil {
			return nil, fmt.Errorf("cannot delete gopher: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("gopher with id %d was deleted", id)), nil
	}
}
