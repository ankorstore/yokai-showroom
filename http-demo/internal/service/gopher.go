package service

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/repository"
	"github.com/ankorstore/yokai/config"
	"github.com/prometheus/client_golang/prometheus"
)

// GopherServiceCounter is a counter for the operation on gophers.
var GopherServiceCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gophers_service_operations_total",
		Help: "Number of operations on the GopherService",
	},
	[]string{
		"operation",
	},
)

// GopherService is the service to manage the gophers.
type GopherService struct {
	config     *config.Config
	repository *repository.GopherRepository
}

// NewGopherService returns a new NewGopherService.
func NewGopherService(config *config.Config, repository *repository.GopherRepository) *GopherService {
	return &GopherService{
		config:     config,
		repository: repository,
	}
}

// List returns a list of all gophers.
func (s *GopherService) List(ctx context.Context, name string, job string) ([]model.Gopher, error) {
	GopherServiceCounter.WithLabelValues("list").Inc()

	return s.repository.FindAll(ctx, repository.GopherRepositoryFindAllParams{
		Name: sql.NullString{String: name, Valid: name != ""},
		Job:  sql.NullString{String: job, Valid: job != ""},
	})
}

// Create creates a new gopher.
func (s *GopherService) Create(ctx context.Context, name string, job string) (int, error) {
	GopherServiceCounter.WithLabelValues("create").Inc()

	return s.repository.Create(ctx, repository.GopherRepositoryCreateParams{
		Name: name,
		Job:  sql.NullString{String: job, Valid: job != ""},
	})
}

// Get returns a gopher by id.
func (s *GopherService) Get(ctx context.Context, id int) (model.Gopher, error) {
	GopherServiceCounter.WithLabelValues("get").Inc()

	return s.repository.Find(ctx, id)
}

// Delete deletes a gopher by id.
func (s *GopherService) Delete(ctx context.Context, id int) error {
	GopherServiceCounter.WithLabelValues("delete").Inc()

	return s.repository.Delete(ctx, id)
}
