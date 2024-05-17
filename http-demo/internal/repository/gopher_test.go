package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/db/sqlc"
	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestFindSuccess(t *testing.T) {
	var querier sqlc.Querier

	internal.RunTest(t, fx.Populate(&querier))

	repo := repository.NewGopherRepository(querier)

	ctx := context.Background()

	params := sqlc.CreateGopherParams{
		Name: "test",
		Job:  sql.NullString{String: "test", Valid: true},
	}

	id, err := repo.Create(ctx, params)
	assert.NoError(t, err)

	foundGopher, err := repo.Find(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, "test", foundGopher.Name)
	assert.Equal(t, "test", foundGopher.Job)
}

func TestFindError(t *testing.T) {
	var db *sql.DB
	var querier sqlc.Querier

	internal.RunTest(t, fx.Populate(&db, &querier))

	repo := repository.NewGopherRepository(querier)

	err := db.Close()
	assert.NoError(t, err)

	gopher, err := repo.Find(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, gopher)
}

func TestFindAllSuccess(t *testing.T) {
	var querier sqlc.Querier

	internal.RunTest(t, fx.Populate(&querier))

	repo := repository.NewGopherRepository(querier)

	ctx := context.Background()

	foundGophers, err := repo.FindAll(ctx)
	assert.NoError(t, err)

	assert.Len(t, foundGophers, 1)
	assert.Equal(t, "alice", foundGophers[0].Name)
	assert.Equal(t, "architect", foundGophers[0].Job)
	assert.Equal(t, "bob", foundGophers[1].Name)
	assert.Equal(t, "builder", foundGophers[1].Job)
	assert.Equal(t, "carl", foundGophers[2].Name)
	assert.Equal(t, "carpenter", foundGophers[2].Job)
}

func TestFindAllError(t *testing.T) {
	var db *sql.DB
	var querier sqlc.Querier

	internal.RunTest(t, fx.Populate(&db, &querier))

	repo := repository.NewGopherRepository(querier)

	err := db.Close()
	assert.NoError(t, err)

	gopher, err := repo.FindAll(context.Background())
	assert.Error(t, err)
	assert.Nil(t, gopher)
}

func TestCreateSuccess(t *testing.T) {
	var querier sqlc.Querier

	internal.RunTest(t, fx.Populate(&querier))

	repo := repository.NewGopherRepository(querier)

	ctx := context.Background()

	id, err := repo.Create(ctx, sqlc.CreateGopherParams{
		Name: "test",
		Job:  sql.NullString{String: "test", Valid: true},
	})
	assert.NoError(t, err)

	foundGopher, err := repo.Find(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, "test", foundGopher.Name)
	assert.Equal(t, "test", foundGopher.Job)
}

func TestCreateError(t *testing.T) {
	var db *sql.DB
	var querier sqlc.Querier

	internal.RunTest(t, fx.Populate(&db, &querier))

	repo := repository.NewGopherRepository(querier)

	ctx := context.Background()

	err := db.Close()
	assert.NoError(t, err)

	_, err = repo.Create(ctx, sqlc.CreateGopherParams{
		Name: "test",
		Job:  sql.NullString{String: "test", Valid: true},
	})
	assert.Error(t, err)
}

func TestDeleteSuccess(t *testing.T) {
	var querier sqlc.Querier

	internal.RunTest(t, fx.Populate(&querier))

	repo := repository.NewGopherRepository(querier)

	ctx := context.Background()

	err := repo.Delete(ctx, 1)
	assert.NoError(t, err)
}

func TestDeleteError(t *testing.T) {
	var db *sql.DB
	var querier sqlc.Querier

	internal.RunTest(t, fx.Populate(&db, &querier))

	repo := repository.NewGopherRepository(querier)

	ctx := context.Background()

	err := db.Close()
	assert.NoError(t, err)

	err = repo.Delete(ctx, 1)
	assert.NoError(t, err)
}
