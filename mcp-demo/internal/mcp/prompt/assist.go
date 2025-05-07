package prompt

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

// JobsAssistPrompt is the MCP prompt for jobs assistance.
type JobsAssistPrompt struct {
	config *config.Config
}

// NewJobsAssistPrompt returns a new JobsAssistPrompt instance.
func NewJobsAssistPrompt(config *config.Config) *JobsAssistPrompt {
	return &JobsAssistPrompt{
		config: config,
	}
}

// Name returns the JobsAssistPrompt name.
func (p *JobsAssistPrompt) Name() string {
	return "jobs-assist"
}

// Options returns the JobsAssistPrompt options.
func (p *JobsAssistPrompt) Options() []mcp.PromptOption {
	return []mcp.PromptOption{
		mcp.WithPromptDescription("gophers jobs assistance"),
	}
}

// Handle returns the JobsAssistPrompt request handler.
func (p *JobsAssistPrompt) Handle() server.PromptHandlerFunc {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		ctx, span := trace.CtxTracer(ctx).Start(ctx, "JobsAssistPrompt.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("JobsAssistPrompt.Handle")

		var jobs []string
		for job := range p.config.GetStringMapString("config.gophers.jobs") {
			jobs = append(jobs, job)
		}

		sort.Strings(jobs)

		jsonJobs, err := json.Marshal(jobs)
		if err != nil {
			return nil, fmt.Errorf("cannot json marshall jobs: %w", err)
		}

		return mcp.NewGetPromptResult(
			"gophers jobs assistance",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent("help by providing information about the gopher jobs"),
				),
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewEmbeddedResource(mcp.TextResourceContents{
						URI:      "jobs://jobs",
						MIMEType: "application/json",
						Text:     string(jsonJobs),
					}),
				),
			},
		), nil
	}
}
