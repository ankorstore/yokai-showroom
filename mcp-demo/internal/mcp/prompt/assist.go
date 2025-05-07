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

type JobsAssistPrompt struct {
	config *config.Config
}

func NewJobsAssistPrompt(config *config.Config) *JobsAssistPrompt {
	return &JobsAssistPrompt{
		config: config,
	}
}

func (p *JobsAssistPrompt) Name() string {
	return "jobs-assist"
}

func (p *JobsAssistPrompt) Options() []mcp.PromptOption {
	return []mcp.PromptOption{
		mcp.WithPromptDescription("gophers jobs assistance"),
	}
}

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
