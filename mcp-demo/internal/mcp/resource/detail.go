package resource

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// JobDetailResource is the MCP resource for job details.
type JobDetailResource struct {
	config *config.Config
}

// NewJobDetailResource returns a new JobDetailResource instance.
func NewJobDetailResource(config *config.Config) *JobDetailResource {
	return &JobDetailResource{
		config: config,
	}
}

// Name returns the JobDetailResource name.
func (r *JobDetailResource) Name() string {
	return "job-detail"
}

// URI returns the JobDetailResource URI.
func (r *JobDetailResource) URI() string {
	return "jobs://jobs/{name}"
}

// Options returns the JobDetailResource options.
func (r *JobDetailResource) Options() []mcp.ResourceTemplateOption {
	return []mcp.ResourceTemplateOption{
		mcp.WithTemplateDescription("details about a specific gopher job"),
		mcp.WithTemplateMIMEType("text/plain"),
	}
}

// Handle returns the JobDetailResource request handler.
func (r *JobDetailResource) Handle() server.ResourceTemplateHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		ctx, span := trace.CtxTracer(ctx).Start(ctx, "ListJobsResource.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("ListJobsResource.Handle")

		uri, err := url.Parse(request.Params.URI)
		if err != nil {
			return nil, fmt.Errorf("cannot parse request URI: %w", err)
		}

		job := path.Base(uri.Path)

		jobDetail, exists := r.config.GetStringMapString("config.gophers.jobs")[job]
		if !exists {
			return nil, fmt.Errorf("job %s does not exist", job)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/plain",
				Text:     jobDetail,
			},
		}, nil
	}
}
