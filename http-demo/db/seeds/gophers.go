package seeds

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai/config"
)

type GophersSeed struct {
	config *config.Config
}

func NewGophersSeed(config *config.Config) *GophersSeed {
	return &GophersSeed{
		config: config,
	}
}

func (s *GophersSeed) Name() string {
	return "gophers"
}

func (s *GophersSeed) Run(ctx context.Context, db *sql.DB) error {
	for name, job := range s.config.GetStringMapString("config.seed.gophers") {
		_, err := db.ExecContext(ctx, "INSERT INTO gophers (name, job) VALUES (?, ?)", name, job)
		if err != nil {
			return err
		}
	}

	return nil
}
