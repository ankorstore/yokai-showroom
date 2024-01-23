package service

import (
	"context"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/repository"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// GopherService is the service to manage the [model.Gopher] model.
type GopherService struct {
	repository *repository.GopherRepository
	config     *config.Config
}

// NewGopherService returns a new NewGopherService.
func NewGopherService(repository *repository.GopherRepository, config *config.Config) *GopherService {
	return &GopherService{
		repository: repository,
		config:     config,
	}
}

// List returns a list of all [model.Gopher].
func (s *GopherService) List(ctx context.Context) ([]model.Gopher, error) {
	ctx, span := s.trace(ctx, "list gophers service")
	if span != nil {
		defer span.End()
	}

	s.log(ctx, "called list gophers")

	return s.repository.FindAll(ctx)
}

// Create creates a new [model.Gopher].
func (s *GopherService) Create(ctx context.Context, gopher *model.Gopher) error {
	ctx, span := s.trace(ctx, "create gopher service")
	if span != nil {
		defer span.End()
	}

	s.log(ctx, "called create gopher")

	return s.repository.Create(ctx, gopher)
}

// Get returns a [model.Gopher] by id.
func (s *GopherService) Get(ctx context.Context, id int) (*model.Gopher, error) {
	ctx, span := s.trace(ctx, "get gopher service")
	if span != nil {
		defer span.End()
	}

	s.log(ctx, "called get gopher")

	return s.repository.Find(ctx, id)
}

// Update updates a provided [model.Gopher].
func (s *GopherService) Update(ctx context.Context, id int, update *model.Gopher) (*model.Gopher, error) {
	ctx, span := s.trace(ctx, "update gopher service")
	if span != nil {
		defer span.End()
	}

	s.log(ctx, "called update gopher")

	gopher, err := s.repository.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	err = s.repository.Update(ctx, gopher, update)
	if err != nil {
		return nil, err
	}

	return gopher, nil
}

// Delete deletes a [model.Gopher] by id.
func (s *GopherService) Delete(ctx context.Context, id int) error {
	ctx, span := s.trace(ctx, "delete gopher service")
	if span != nil {
		defer span.End()
	}

	s.log(ctx, "called delete gopher")

	gopher, err := s.repository.Find(ctx, id)
	if err != nil {
		return err
	}

	err = s.repository.Delete(ctx, gopher)
	if err != nil {
		return err
	}

	return nil
}

func (s *GopherService) trace(ctx context.Context, spanName string) (context.Context, oteltrace.Span) {
	if s.config.GetBool("config.service.gopher.trace") {
		return trace.CtxTracerProvider(ctx).Tracer("gopher-service").Start(ctx, spanName)
	}

	return ctx, nil
}

func (s *GopherService) log(ctx context.Context, message string) {
	if s.config.GetBool("config.service.gopher.log") {
		log.CtxLogger(ctx).Info().Msg(message)
	}
}
