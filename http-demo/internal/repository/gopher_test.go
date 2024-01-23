package repository_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func TestFindSuccess(t *testing.T) {
	var db *gorm.DB
	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	ctx := context.Background()

	gopher := &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	}

	err := repo.Create(ctx, gopher)
	assert.NoError(t, err)

	foundGopher, err := repo.Find(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "gopher 1", foundGopher.Name)
	assert.Equal(t, "job 1", foundGopher.Job)
}

func TestFindError(t *testing.T) {
	var db *gorm.DB
	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	d, err := db.DB()
	assert.NoError(t, err)

	err = d.Close()
	assert.NoError(t, err)

	gopher, err := repo.Find(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, gopher)
}

func TestFindAllSuccess(t *testing.T) {
	var db *gorm.DB
	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	ctx := context.Background()

	gopher := &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	}

	err := repo.Create(ctx, gopher)
	assert.NoError(t, err)

	foundGophers, err := repo.FindAll(ctx)
	assert.NoError(t, err)

	assert.Len(t, foundGophers, 1)
	assert.Equal(t, "gopher 1", foundGophers[0].Name)
	assert.Equal(t, "job 1", foundGophers[0].Job)
}

func TestFindAllError(t *testing.T) {
	var db *gorm.DB
	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	d, err := db.DB()
	assert.NoError(t, err)

	err = d.Close()
	assert.NoError(t, err)

	gopher, err := repo.FindAll(context.Background())
	assert.Error(t, err)
	assert.Nil(t, gopher)
}

func TestCreateSuccess(t *testing.T) {
	var db *gorm.DB
	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	ctx := context.Background()

	gopher := &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	}

	err := repo.Create(ctx, gopher)
	assert.NoError(t, err)

	foundGopher, err := repo.Find(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "gopher 1", foundGopher.Name)
	assert.Equal(t, "job 1", foundGopher.Job)
}

func TestCreateError(t *testing.T) {
	var db *gorm.DB
	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	d, err := db.DB()
	assert.NoError(t, err)

	err = d.Close()
	assert.NoError(t, err)

	gopher := &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	}

	err = repo.Create(context.Background(), gopher)
	assert.Error(t, err)
}

func TestUpdateSuccess(t *testing.T) {
	var db *gorm.DB
	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	ctx := context.Background()

	gopher := &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	}

	err := repo.Create(ctx, gopher)
	assert.NoError(t, err)

	err = repo.Update(ctx, gopher, &model.Gopher{
		Name: "new gopher 1",
		Job:  "new job 1",
	})
	assert.NoError(t, err)
	assert.Equal(t, "new gopher 1", gopher.Name)
	assert.Equal(t, "new job 1", gopher.Job)
}

func TestUpdateError(t *testing.T) {
	var db *gorm.DB
	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	d, err := db.DB()
	assert.NoError(t, err)

	err = d.Close()
	assert.NoError(t, err)

	gopher := &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	}

	err = repo.Update(context.Background(), gopher, &model.Gopher{
		Name: "new gopher 1",
		Job:  "new job 1",
	})
	assert.Error(t, err)
}

func TestDeleteSuccess(t *testing.T) {
	var db *gorm.DB
	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	ctx := context.Background()

	gopher := &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	}

	err := repo.Create(ctx, gopher)
	assert.NoError(t, err)

	err = repo.Delete(ctx, gopher)
	assert.NoError(t, err)
}

func TestDeleteError(t *testing.T) {
	var db *gorm.DB
	internal.RunTest(t, fx.Populate(&db))

	repo := repository.NewGopherRepository(db)

	d, err := db.DB()
	assert.NoError(t, err)

	err = d.Close()
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), &model.Gopher{})
	assert.Error(t, err)
}
