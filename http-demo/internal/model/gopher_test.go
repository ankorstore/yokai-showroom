package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGopherModelProperties(t *testing.T) {
	t.Parallel()

	now := time.Now()

	gopher := &model.Gopher{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: gorm.DeletedAt{},
		Name:      "gopher 1",
		Job:       "job 1",
	}

	assert.Equal(t, uint(1), gopher.ID)
	assert.Equal(t, now, gopher.CreatedAt)
	assert.Equal(t, now, gopher.UpdatedAt)
	assert.Equal(t, gorm.DeletedAt{}, gopher.DeletedAt)
	assert.Equal(t, "gopher 1", gopher.Name)
	assert.Equal(t, "job 1", gopher.Job)
}

func TestGopherModelAsJson(t *testing.T) {
	t.Parallel()

	fixedDate := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local)

	gopher := &model.Gopher{
		ID:        1,
		CreatedAt: fixedDate,
		UpdatedAt: fixedDate,
		DeletedAt: gorm.DeletedAt{},
		Name:      "gopher 1",
		Job:       "job 1",
	}

	jsonData, err := json.Marshal(gopher)
	assert.NoError(t, err)

	var parsedGopher *model.Gopher
	err = json.Unmarshal(jsonData, &parsedGopher)
	assert.NoError(t, err)

	assert.Equal(t, gopher, parsedGopher)
}
