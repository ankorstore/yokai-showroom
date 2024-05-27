package seeds

import (
	"context"
	"database/sql"
	"sort"

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
	var txErr error

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, txErr = tx.ExecContext(ctx, "DELETE FROM gophers")

	seedData := s.config.GetStringMapString("config.seed.gophers")

	names := make([]string, 0, len(seedData))
	for name := range seedData {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		_, txErr = tx.ExecContext(ctx, "INSERT INTO gophers (name, job) VALUES (?, ?)", name, seedData[name])
	}

	if txErr != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}
