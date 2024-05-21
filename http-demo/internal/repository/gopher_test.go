package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestFindError(t *testing.T) {
	var db *sql.DB

	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	err := db.Close()
	assert.NoError(t, err)

	_, err = repo.Find(context.Background(), 1)
	assert.Error(t, err)
}

func TestFindAllError(t *testing.T) {
	var db *sql.DB

	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	err := db.Close()
	assert.NoError(t, err)

	_, err = repo.FindAll(context.Background(), repository.GopherRepositoryFindAllParams{})
	assert.Error(t, err)
}

func TestCreateError(t *testing.T) {
	var db *sql.DB

	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	ctx := context.Background()

	err := db.Close()
	assert.NoError(t, err)

	_, err = repo.Create(ctx, repository.GopherRepositoryCreateParams{})
	assert.Error(t, err)
}

func TestDeleteError(t *testing.T) {
	var db *sql.DB

	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	ctx := context.Background()

	err := db.Close()
	assert.NoError(t, err)

	err = repo.Delete(ctx, 1)
	assert.Error(t, err)
}
