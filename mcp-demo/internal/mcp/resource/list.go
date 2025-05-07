package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// JobsListResource is the MCP resource template for jobs list.
type JobsListResource struct {
	config *config.Config
}

// NewJobsListResource returns a new JobsListResource instance.
func NewJobsListResource(config *config.Config) *JobsListResource {
	return &JobsListResource{
		config: config,
	}
}

// Name returns the JobsListResource name.
func (r *JobsListResource) Name() string {
	return "jobs-list"
}

// URI returns the JobsListResource URI.
func (r *JobsListResource) URI() string {
	return "jobs://jobs"
}

// Options returns the JobsListResource options.
func (r *JobsListResource) Options() []mcp.ResourceOption {
	return []mcp.ResourceOption{
		mcp.WithResourceDescription("list of all the possible gophers jobs"),
		mcp.WithMIMEType(`application/json`),
	}
}

// Handle returns the JobsListResource request handler.
func (r *JobsListResource) Handle() server.ResourceHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		ctx, span := trace.CtxTracer(ctx).Start(ctx, "JobsListResource.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("JobsListResource.Handle")

		var jobs []string
		for job := range r.config.GetStringMapString("config.gophers.jobs") {
			jobs = append(jobs, job)
		}

		sort.Strings(jobs)

		jsonJobs, err := json.Marshal(jobs)
		if err != nil {
			return nil, fmt.Errorf("cannot json marshall jobs: %w", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonJobs),
			},
		}, nil
	}
}
